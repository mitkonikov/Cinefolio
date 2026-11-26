package legacyimport

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"

	"cinefolio/internal/catalog"
	"cinefolio/internal/media"
	"gorm.io/gorm"
)

type Options struct {
	Export    string
	MediaRoot string
	DryRun    bool
}

type Report struct {
	DryRun     bool           `json:"dryRun"`
	Created    map[string]int `json:"created"`
	Existing   map[string]int `json:"existing"`
	MediaBytes int64          `json:"mediaBytes"`
}

type importer struct {
	tx            *gorm.DB
	catalog       *catalog.Repository
	report        Report
	legacy        *os.Root
	target        *media.Storage
	staged        []string
	ids           map[string]map[string]uint64
	sites         map[string]string
	playlistSites map[string]string
}

type mapping struct {
	Kind     string `gorm:"primaryKey"`
	LegacyID string `gorm:"primaryKey"`
	RecordID uint64
}

func (mapping) TableName() string { return "legacy_ids" }

const defaultSite = "cinefolio"

var dryRollback = errors.New("dry-run rollback")
var nonSlug = regexp.MustCompile("[^a-z0-9]+")

func Run(ctx context.Context, db *gorm.DB, target *media.Storage, options Options) (Report, error) {
	report := Report{DryRun: options.DryRun, Created: map[string]int{}, Existing: map[string]int{}}
	input, err := os.ReadFile(options.Export)
	if err != nil {
		return report, err
	}
	var data source
	if err := json.Unmarshal(input, &data); err != nil {
		return report, fmt.Errorf("decode export: %w", err)
	}
	if data.Sites == nil || data.Films == nil || data.Playlists == nil || data.Memberships == nil || data.Thumbnails == nil || data.Trailers == nil {
		return report, errors.New("export must contain all six collection arrays")
	}
	legacy, err := os.OpenRoot(options.MediaRoot)
	if err != nil {
		return report, err
	}
	defer legacy.Close()
	i := importer{report: report, legacy: legacy, target: target, ids: make(map[string]map[string]uint64),
		sites: make(map[string]string), playlistSites: make(map[string]string)}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		i.tx = tx
		i.catalog = catalog.New(tx, defaultSite)
		if err := i.apply(ctx, data); err != nil {
			return err
		}
		if options.DryRun {
			return dryRollback
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, dryRollback) {
			err = nil
		}
		for _, key := range i.staged {
			if cleanErr := target.Remove(key); cleanErr != nil {
				err = errors.Join(err, cleanErr)
			}
		}
	}
	return i.report, err
}

func (i *importer) apply(ctx context.Context, data source) error {
	legacyZone, err := time.LoadLocation("Europe/Skopje")
	if err != nil {
		return err
	}
	for _, site := range data.Sites {
		if err := i.site(site); err != nil {
			return err
		}
	}
	for _, thumbnail := range data.Thumbnails {
		if err := i.entry("thumbnails", thumbnail.ID, func() (uint64, error) {
			asset, err := i.asset(media.Thumbnail, "thumbs", thumbnail.Filepath)
			if err != nil {
				return 0, err
			}
			t, _, err := i.catalog.SaveThumbnail(ctx, 0, catalog.ThumbnailInput{Title: thumbnail.Title, Crop: thumbnail.Crop}, asset)
			return t.ID, err
		}); err != nil {
			return err
		}
	}
	for _, film := range data.Films {
		if err := i.entry("films", film.ID, func() (uint64, error) {
			var thumbnail *uint64
			if film.Thumbnail != "" {
				id, err := i.lookup("thumbnails", film.Thumbnail)
				if err != nil {
					return 0, err
				}
				thumbnail = &id
			}
			date, err := publicationDate(film.Published, legacyZone)
			if err != nil {
				return 0, err
			}
			var asset *media.Asset
			if film.Filepath != "" {
				asset, err = i.asset(media.Film, "files", film.Filepath)
				if err != nil {
					return 0, err
				}
			}
			f, _, err := i.catalog.SaveFilm(ctx, 0, catalog.FilmInput{Title: film.Title, Description: film.Description,
				Published: date, Link: film.Link, ThumbnailID: thumbnail}, asset)
			return f.ID, err
		}); err != nil {
			return err
		}
	}
	for _, playlist := range data.Playlists {
		slug, ok := i.sites[playlist.Site]
		if !ok {
			return fmt.Errorf("playlist %s references unknown site %s", playlist.ID, playlist.Site)
		}
		i.playlistSites[playlist.ID] = slug
		if err := i.entry("playlists", playlist.ID, func() (uint64, error) {
			visible := playlist.Visible == nil || *playlist.Visible
			p, err := catalog.New(i.tx, slug).SavePlaylist(ctx, 0, catalog.PlaylistInput{Name: playlist.Name, Position: playlist.Order, Visible: visible})
			return p.ID, err
		}); err != nil {
			return err
		}
	}
	pairs := make(map[[2]string]string, len(data.Memberships))
	for _, member := range data.Memberships {
		pair := [2]string{member.PlaylistID, member.FileID}
		if previous, ok := pairs[pair]; ok {
			return fmt.Errorf("memberships %s and %s link the same playlist and film", previous, member.ID)
		}
		pairs[pair] = member.ID
		if err := i.entry("memberships", member.ID, func() (uint64, error) {
			playlistID, err := i.lookup("playlists", member.PlaylistID)
			if err != nil {
				return 0, err
			}
			filmID, err := i.lookup("films", member.FileID)
			if err != nil {
				return 0, err
			}
			err = catalog.New(i.tx, i.playlistSites[member.PlaylistID]).SetMembership(ctx, playlistID, filmID, true)
			return playlistID, err
		}); err != nil {
			return err
		}
	}
	var featured string
	for _, trailer := range data.Trailers {
		if trailer.Visible {
			if featured != "" {
				return fmt.Errorf("trailers %s and %s are both visible; choose one in the export", featured, trailer.ID)
			}
			featured = trailer.ID
		}
		if err := i.entry("trailers", trailer.ID, func() (uint64, error) {
			filmID, err := i.lookup("films", trailer.Film)
			if err != nil {
				return 0, err
			}
			asset, err := i.asset(media.Trailer, "trailers", trailer.Filepath)
			if err != nil {
				return 0, err
			}
			t, _, err := i.catalog.SaveTrailer(ctx, 0, catalog.TrailerInput{FilmID: filmID, Featured: trailer.Visible}, asset)
			return t.ID, err
		}); err != nil {
			return err
		}
	}
	return nil
}

func publicationDate(raw string, zone *time.Location) (string, error) {
	if _, err := time.Parse(time.DateOnly, raw); err == nil {
		return raw, nil
	}
	date, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return "", errors.New("missing or invalid publication date")
	}
	// AdminJS stored the date picker’s local midnight as a UTC instant.
	return date.In(zone).Format(time.DateOnly), nil
}

func (i *importer) site(item legacySite) error {
	name := strings.TrimSpace(item.Name)
	slug := strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(name), "-"), "-")
	if slug == "" {
		return fmt.Errorf("site %s has no usable slug", item.ID)
	}
	err := i.entry("sites", item.ID, func() (uint64, error) {
		site := catalog.Site{Name: name, Slug: slug}
		if slug == defaultSite && strings.EqualFold(name, "Cinefolio") {
			err := i.tx.Where("slug = ?", defaultSite).Take(&site).Error
			return site.ID, err
		}
		err := i.tx.Create(&site).Error
		return site.ID, err
	})
	if err != nil {
		return err
	}
	var actual catalog.Site
	if err := i.tx.Take(&actual, i.ids["sites"][item.ID]).Error; err != nil {
		return err
	}
	i.sites[item.ID] = actual.Slug
	return nil
}

func (i *importer) entry(kind, legacyID string, create func() (uint64, error)) error {
	value, err := hex.DecodeString(legacyID)
	if err != nil || len(value) != 12 {
		return fmt.Errorf("%s: invalid legacy ObjectID %q", kind, legacyID)
	}
	if i.ids[kind] == nil {
		i.ids[kind] = make(map[string]uint64)
	}
	if i.ids[kind][legacyID] != 0 {
		return fmt.Errorf("%s: duplicate legacy ObjectID %s", kind, legacyID)
	}
	var record mapping
	err = i.tx.Where("kind = ? AND legacy_id = ?", kind, legacyID).Take(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		id, err := create()
		if err != nil {
			return fmt.Errorf("%s %s: %w", kind, legacyID, err)
		}
		record = mapping{Kind: kind, LegacyID: legacyID, RecordID: id}
		if err := i.tx.Create(&record).Error; err != nil {
			return err
		}
		i.report.Created[kind]++
	} else if err != nil {
		return err
	} else {
		i.report.Existing[kind]++
	}
	i.ids[kind][legacyID] = record.RecordID
	return nil
}

func (i *importer) lookup(kind, id string) (uint64, error) {
	if recordID := i.ids[kind][id]; recordID != 0 {
		return recordID, nil
	}
	return 0, fmt.Errorf("references unknown %s record %q", kind, id)
}

func (i *importer) asset(kind media.Kind, folder, name string) (*media.Asset, error) {
	if !fs.ValidPath(name) || name == "." || strings.Contains(name, "\\") {
		return nil, errors.New("invalid legacy media path")
	}
	file, err := i.legacy.Open(path.Join(folder, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("legacy media must be a regular file")
	}
	sourceHash := sha256.New()
	asset, err := i.target.Stage(kind, path.Base(name), io.TeeReader(file, sourceHash))
	if err != nil {
		return nil, err
	}
	i.staged = append(i.staged, asset.StorageKey)
	staged, err := i.target.Open(asset.StorageKey)
	if err != nil {
		return nil, err
	}
	defer staged.Close()
	targetHash := sha256.New()
	size, err := io.Copy(targetHash, staged)
	if err != nil {
		return nil, err
	}
	if size != info.Size() || !bytes.Equal(sourceHash.Sum(nil), targetHash.Sum(nil)) {
		return nil, errors.New("media checksum or size mismatch")
	}
	i.report.MediaBytes += size
	return &asset, nil
}
