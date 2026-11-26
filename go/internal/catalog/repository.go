package catalog

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cinefolio/internal/media"
	"cinefolio/internal/store"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db   *gorm.DB
	site string
}

func New(db *gorm.DB, site string) *Repository { return &Repository{db: db, site: site} }

type Page[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}

type Filter struct {
	Query  string
	Offset int
	Limit  int
}

func page[T any](q *gorm.DB, filter Filter, fields ...string) (Page[T], error) {
	p := Page[T]{Items: []T{}}
	if filter.Query != "" && len(fields) > 0 {
		needle := "%" + strings.NewReplacer("!", "!!", "%", "!%", "_", "!_").Replace(filter.Query) + "%"
		conditions := make([]string, 0, len(fields))
		args := make([]any, 0, len(fields))
		for _, field := range fields {
			conditions = append(conditions, field+" LIKE ? ESCAPE '!'")
			args = append(args, needle)
		}
		q = q.Where(strings.Join(conditions, " OR "), args...)
	}
	if err := q.Session(&gorm.Session{}).Count(&p.Total).Error; err != nil {
		return p, err
	}
	err := q.Offset(filter.Offset).Limit(filter.Limit).Find(&p.Items).Error
	return p, err
}

func (r *Repository) Films(ctx context.Context, f Filter) (Page[Film], error) {
	return page[Film](r.db.WithContext(ctx).Model(&Film{}).Preload("Asset").Preload("Thumbnail.Asset").
		Preload("Playlists", "site_id = (?)", r.siteID(r.db)).Order("published DESC, id"), f, "title", "description")
}
func (r *Repository) Thumbnails(ctx context.Context, f Filter) (Page[Thumbnail], error) {
	return page[Thumbnail](r.db.WithContext(ctx).Model(&Thumbnail{}).Preload("Asset").Order("title, id"), f, "title")
}
func (r *Repository) Playlists(ctx context.Context, f Filter) (Page[Playlist], error) {
	return page[Playlist](r.db.WithContext(ctx).Model(&Playlist{}).Where("site_id = (?)", r.siteID(r.db)).Preload("Films", orderedFilms).Order("position, id"), f, "name")
}
func (r *Repository) Trailers(ctx context.Context, f Filter) (Page[Trailer], error) {
	return page[Trailer](r.db.WithContext(ctx).Model(&Trailer{}).Where("site_id = (?)", r.siteID(r.db)).Preload("Asset").Preload("Film").Order("featured DESC, id"), f)
}

func (r *Repository) siteID(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{NewDB: true}).Model(&Site{}).Select("id").Where("slug = ?", r.site)
}

func (r *Repository) activeSite(tx *gorm.DB) (uint64, error) {
	var site Site
	err := tx.Where("slug = ?", r.site).Take(&site).Error
	return site.ID, err
}

func orderedFilms(db *gorm.DB) *gorm.DB { return db.Order("published DESC, id") }

func (r *Repository) SetMembership(ctx context.Context, playlistID, filmID uint64, present bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p Playlist
		if err := tx.Where("site_id = (?)", r.siteID(tx)).Take(&p, playlistID).Error; err != nil {
			return err
		}
		membership := playlistFilm{PlaylistID: playlistID, FilmID: filmID}
		if present {
			return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&membership).Error
		}
		return tx.Delete(&membership).Error
	})
}

type playlistFilm struct {
	PlaylistID uint64 `gorm:"primaryKey"`
	FilmID     uint64 `gorm:"primaryKey"`
}

func (r *Repository) SavePlaylist(ctx context.Context, id uint64, in PlaylistInput) (Playlist, error) {
	p := Playlist{ID: id, Name: strings.TrimSpace(in.Name), Visible: in.Visible, Position: in.Position}
	if p.Name == "" || len(p.Name) > 200 {
		return p, store.ValidationError("playlist name must contain 1–200 bytes")
	}
	if p.Position < 0 {
		return p, store.ValidationError("position must be nonnegative")
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		siteID, err := r.activeSite(tx)
		if err != nil {
			return err
		}
		p.SiteID = siteID
		if id == 0 {
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
		} else {
			result := tx.Model(&p).Where("site_id = ?", siteID).Select("name", "visible", "position").Updates(&p)
			if err := changed(result); err != nil {
				return err
			}
		}
		if in.FilmIDs == nil {
			return nil
		} // Optional on edit, so metadata changes can leave memberships intact.
		if err := tx.Where("playlist_id = ?", p.ID).Delete(&playlistFilm{}).Error; err != nil {
			return err
		}
		seen := make(map[uint64]bool, len(in.FilmIDs))
		links := make([]playlistFilm, 0, len(in.FilmIDs))
		for _, filmID := range in.FilmIDs {
			if seen[filmID] {
				return store.ValidationError("duplicate film membership")
			}
			seen[filmID] = true
			links = append(links, playlistFilm{PlaylistID: p.ID, FilmID: filmID})
		}
		if len(links) == 0 {
			return nil
		}
		return tx.CreateInBatches(&links, 100).Error
	})
	return p, err
}

func (r *Repository) DeletePlaylist(ctx context.Context, id uint64) error {
	return changed(r.db.WithContext(ctx).Where("site_id = (?)", r.siteID(r.db)).Delete(&Playlist{}, id))
}

func changed(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return store.ErrNotFound
	}
	return nil
}

func Missing(err error) bool {
	return errors.Is(err, store.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound)
}

func swapAsset(tx *gorm.DB, old *uint64, upload *media.Asset, remove bool) (*uint64, error) {
	if upload != nil {
		if err := tx.Create(upload).Error; err != nil {
			return nil, err
		}
		return &upload.ID, nil
	}
	if remove {
		return nil, nil
	}
	return old, nil
}

func discardAsset(tx *gorm.DB, old, next *uint64) (string, error) {
	if old == nil || (next != nil && *old == *next) {
		return "", nil
	}
	var asset media.Asset
	if err := tx.Take(&asset, *old).Error; err != nil {
		return "", err
	}
	if err := tx.Delete(&asset).Error; err != nil {
		return "", err
	}
	return asset.StorageKey, nil
}

func requiredAsset(id *uint64) error {
	if id == nil {
		return store.ValidationError("a media upload is required")
	}
	return nil
}

func (r *Repository) Feature(ctx context.Context, tx *gorm.DB, siteID, id uint64) error {
	if err := tx.WithContext(ctx).Model(&Trailer{}).Where("site_id = ? AND id <> ? AND featured", siteID, id).Update("featured", false).Error; err != nil {
		return err
	}
	return changed(tx.Model(&Trailer{}).Where("site_id = ? AND id = ?", siteID, id).Update("featured", true))
}

func (r *Repository) CheckSite(ctx context.Context) error {
	_, err := r.activeSite(r.db.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("public site %q: %w", r.site, err)
	}
	return nil
}
