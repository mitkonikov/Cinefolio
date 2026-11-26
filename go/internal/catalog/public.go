package catalog

import (
	"context"
	"errors"
	"fmt"

	"cinefolio/internal/media"
	"gorm.io/gorm"
)

type Portfolio struct {
	Trailer   *PublicTrailer   `json:"trailer"`
	Playlists []PublicPlaylist `json:"playlists"`
}

type PublicTrailer struct {
	URL  string     `json:"url"`
	Film PublicFilm `json:"film"`
}

type PublicPlaylist struct {
	ID    uint64       `json:"id"`
	Name  string       `json:"name"`
	Films []PublicFilm `json:"films"`
}

type PublicFilm struct {
	ID          uint64           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Published   string           `json:"published"`
	Link        string           `json:"link"`
	VideoURL    string           `json:"videoUrl"`
	Thumbnail   *PublicThumbnail `json:"thumbnail"`
}

type PublicThumbnail struct {
	URL  string `json:"url"`
	Crop bool   `json:"crop"`
}

func publicFilm(film Film) PublicFilm {
	result := PublicFilm{ID: film.ID, Title: film.Title, Description: film.Description, Published: film.Published, Link: film.Link}
	if film.AssetID != nil {
		result.VideoURL = fmt.Sprintf("/media/%d", *film.AssetID)
	}
	if film.Thumbnail != nil {
		result.Thumbnail = &PublicThumbnail{URL: fmt.Sprintf("/media/%d", film.Thumbnail.AssetID), Crop: film.Thumbnail.Crop}
	}
	return result
}

func (r *Repository) Public(ctx context.Context) (Portfolio, error) {
	p := Portfolio{Playlists: []PublicPlaylist{}}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		siteID, err := r.activeSite(tx)
		if err != nil {
			return err
		}
		var lists []Playlist
		if err := tx.Where("site_id = ? AND visible", siteID).
			Order("position, id").Preload("Films", orderedFilms).
			Preload("Films.Thumbnail").Find(&lists).Error; err != nil {
			return err
		}
		for _, list := range lists {
			entry := PublicPlaylist{ID: list.ID, Name: list.Name, Films: make([]PublicFilm, 0, len(list.Films))}
			for _, film := range list.Films {
				entry.Films = append(entry.Films, publicFilm(film))
			}
			p.Playlists = append(p.Playlists, entry)
		}
		var trailer Trailer
		err = tx.Where("site_id = ? AND featured", siteID).Preload("Film.Thumbnail").Take(&trailer).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		p.Trailer = &PublicTrailer{URL: fmt.Sprintf("/media/%d", trailer.AssetID), Film: publicFilm(*trailer.Film)}
		return nil
	})
	return p, err
}

const visibleFilmIDs = `
	SELECT pf.film_id FROM playlist_films pf
	JOIN playlists p ON p.id = pf.playlist_id WHERE p.site_id = (?) AND p.visible
	UNION
	SELECT film_id FROM trailers WHERE site_id = (?) AND featured
`

func (r *Repository) Asset(ctx context.Context, id uint64, admin bool) (media.Asset, error) {
	var asset media.Asset
	db := r.db.WithContext(ctx)
	if !admin {
		q := `EXISTS (SELECT 1 FROM trailers WHERE site_id = (?) AND featured AND asset_id = assets.id)
			OR EXISTS (SELECT 1 FROM films f LEFT JOIN thumbnails t ON t.id = f.thumbnail_id
				WHERE (f.asset_id = assets.id OR t.asset_id = assets.id) AND f.id IN (` + visibleFilmIDs + `))`
		siteID := r.siteID(db)
		db = db.Where(q, siteID, siteID, siteID)
	}
	err := db.Take(&asset, id).Error
	return asset, err
}
