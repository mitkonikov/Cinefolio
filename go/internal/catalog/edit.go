package catalog

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cinefolio/internal/media"
	"cinefolio/internal/store"
	"gorm.io/gorm"
)

func (r *Repository) SaveFilm(ctx context.Context, id uint64, in FilmInput, upload *media.Asset) (Film, string, error) {
	film := Film{ID: id, Title: strings.TrimSpace(in.Title), Description: in.Description, Published: in.Published, Link: strings.TrimSpace(in.Link), ThumbnailID: in.ThumbnailID}
	if film.Title == "" || len(film.Title) > 200 || len(film.Description) > 50000 {
		return film, "", store.ValidationError("check the film title and description length")
	}
	if _, err := time.Parse(time.DateOnly, film.Published); err != nil {
		return film, "", store.ValidationError("enter a publication date")
	}
	if film.Link != "" {
		link, err := url.Parse(film.Link)
		if err != nil || len(film.Link) > 4096 || (link.Scheme != "https" && link.Scheme != "http") || link.Hostname() == "" || link.User != nil || !validPort(link) {
			return film, "", store.ValidationError("enter an HTTP or HTTPS film link")
		}
	}
	var obsolete string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var old Film
		if id != 0 {
			if err := tx.Take(&old, id).Error; err != nil {
				return err
			}
			if err := requireRevision(old.Revision, in.Revision); err != nil {
				return err
			}
		}
		var err error
		film.AssetID, err = swapAsset(tx, old.AssetID, upload, in.RemoveUpload)
		if err != nil {
			return err
		}
		if film.AssetID == nil && film.Link == "" {
			return store.ValidationError("add a video upload or an external link")
		}
		if id == 0 {
			if err := tx.Create(&film).Error; err != nil {
				return err
			}
		} else if err := tx.Model(&film).Select("title", "description", "published", "link", "asset_id", "thumbnail_id").Updates(&film).Error; err != nil {
			return err
		}
		if in.PlaylistIDs != nil {
			if err := r.replaceFilmMemberships(tx, film.ID, in.PlaylistIDs); err != nil {
				return err
			}
		}
		obsolete, err = discardAsset(tx, old.AssetID, film.AssetID)
		if err != nil {
			return err
		}
		return readRevision(tx, &film, film.ID, &film.Revision)
	})
	return film, obsolete, err
}

func validPort(link *url.URL) bool {
	port := link.Port()
	if port == "" {
		return true
	}
	_, err := strconv.ParseUint(port, 10, 16)
	return err == nil
}

func (r *Repository) replaceFilmMemberships(tx *gorm.DB, filmID uint64, ids []uint64) error {
	playlists := tx.Session(&gorm.Session{NewDB: true}).Model(&Playlist{}).Where("site_id = (?)", r.siteID(tx))
	if len(ids) > 0 {
		var count int64
		if err := playlists.Session(&gorm.Session{}).Where("id IN ?", ids).Count(&count).Error; err != nil {
			return err
		}
		if int(count) != len(ids) {
			return store.ValidationError("select distinct playlists belonging to this site")
		}
	}
	stale := tx.Where("film_id = ? AND playlist_id IN (?)", filmID, playlists.Select("id"))
	if len(ids) > 0 {
		stale = stale.Where("playlist_id NOT IN ?", ids)
	}
	rows := make([]playlistFilm, 0, len(ids))
	for _, id := range ids {
		rows = append(rows, playlistFilm{PlaylistID: id, FilmID: filmID})
	}
	return syncMemberships(tx, stale, rows)
}

func (r *Repository) SaveThumbnail(ctx context.Context, id uint64, in ThumbnailInput, upload *media.Asset) (Thumbnail, string, error) {
	t := Thumbnail{ID: id, Title: strings.TrimSpace(in.Title), Crop: in.Crop}
	if len(t.Title) > 200 {
		return t, "", store.ValidationError("thumbnail title must be at most 200 bytes")
	}
	var obsolete string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var old *uint64
		if id != 0 {
			var current Thumbnail
			if err := tx.Take(&current, id).Error; err != nil {
				return err
			}
			if err := requireRevision(current.Revision, in.Revision); err != nil {
				return err
			}
			old = &current.AssetID
		}
		next, err := swapAsset(tx, old, upload, false)
		if err != nil {
			return err
		}
		if err := requiredAsset(next); err != nil {
			return err
		}
		t.AssetID = *next
		if id == 0 {
			return tx.Create(&t).Error
		}
		if err := tx.Model(&t).Select("title", "crop", "asset_id").Updates(&t).Error; err != nil {
			return err
		}
		obsolete, err = discardAsset(tx, old, next)
		if err != nil {
			return err
		}
		return readRevision(tx, &t, t.ID, &t.Revision)
	})
	return t, obsolete, err
}

func (r *Repository) EnableCrop(ctx context.Context, records []RecordRevision) error {
	if len(records) == 0 || len(records) > 200 {
		return store.ValidationError("select between 1 and 200 thumbnails")
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, expected := range records {
			var current Thumbnail
			if err := tx.Take(&current, expected.ID).Error; err != nil {
				return err
			}
			if err := requireRevision(current.Revision, expected.Revision); err != nil {
				return err
			}
			if !current.Crop {
				if err := tx.Model(&current).Update("crop", true).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *Repository) SaveTrailer(ctx context.Context, id uint64, in TrailerInput, upload *media.Asset) (Trailer, string, error) {
	t := Trailer{ID: id, FilmID: in.FilmID}
	if in.FilmID == 0 {
		return t, "", store.ValidationError("select a film")
	}
	var obsolete string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		siteID, err := r.activeSite(tx)
		if err != nil {
			return err
		}
		t.SiteID = siteID
		var old *uint64
		if id != 0 {
			var current Trailer
			if err := tx.Where("site_id = ?", siteID).Take(&current, id).Error; err != nil {
				return err
			}
			if err := requireRevision(current.Revision, in.Revision); err != nil {
				return err
			}
			old = &current.AssetID
		}
		next, err := swapAsset(tx, old, upload, false)
		if err != nil {
			return err
		}
		if err := requiredAsset(next); err != nil {
			return err
		}
		t.AssetID = *next
		if id == 0 {
			if err := tx.Create(&t).Error; err != nil {
				return err
			}
		} else if err := tx.Model(&t).Select("film_id", "asset_id", "featured").Updates(&t).Error; err != nil {
			return err
		}
		if in.Featured {
			if err := r.Feature(ctx, tx, siteID, t.ID); err != nil {
				return err
			}
			t.Featured = true
		}
		obsolete, err = discardAsset(tx, old, next)
		if err != nil {
			return err
		}
		return readRevision(tx, &t, t.ID, &t.Revision)
	})
	return t, obsolete, err
}

func (r *Repository) DeleteFilm(ctx context.Context, id, revision uint64) (string, error) {
	return removeRecord(r.db.WithContext(ctx), id, revision, func(f Film) (uint64, *uint64) { return f.Revision, f.AssetID })
}
func (r *Repository) DeleteThumbnail(ctx context.Context, id, revision uint64) (string, error) {
	return removeRecord(r.db.WithContext(ctx), id, revision, func(t Thumbnail) (uint64, *uint64) { return t.Revision, &t.AssetID })
}
func (r *Repository) DeleteTrailer(ctx context.Context, id, revision uint64) (string, error) {
	return removeRecord(r.db.WithContext(ctx).Where("site_id = (?)", r.siteID(r.db)), id, revision, func(t Trailer) (uint64, *uint64) { return t.Revision, &t.AssetID })
}

func removeRecord[T any](db *gorm.DB, id, revision uint64, fields func(T) (uint64, *uint64)) (string, error) {
	var obsolete string
	err := db.Transaction(func(tx *gorm.DB) error {
		var record T
		if err := tx.Take(&record, id).Error; err != nil {
			return err
		}
		current, asset := fields(record)
		if err := requireRevision(current, revision); err != nil {
			return err
		}
		tx = tx.Session(&gorm.Session{NewDB: true})
		if err := tx.Delete(&record).Error; err != nil {
			return err
		}
		var err error
		obsolete, err = discardAsset(tx, asset, nil)
		return err
	})
	return obsolete, err
}
