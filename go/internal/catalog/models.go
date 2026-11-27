package catalog

import "cinefolio/internal/media"

type RecordRevision struct {
	ID       uint64 `json:"id"`
	Revision uint64 `json:"revision"`
}

type Site struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Film struct {
	Revision    uint64       `json:"revision" gorm:"default:1"`
	ID          uint64       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Published   string       `json:"published"`
	Link        string       `json:"link"`
	AssetID     *uint64      `json:"assetId"`
	Asset       *media.Asset `json:"asset,omitempty"`
	ThumbnailID *uint64      `json:"thumbnailId"`
	Thumbnail   *Thumbnail   `json:"thumbnail,omitempty"`
	Playlists   []Playlist   `json:"playlists,omitempty" gorm:"many2many:playlist_films"`
}

type FilmInput struct {
	Revision     uint64   `json:"revision"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Published    string   `json:"published"`
	Link         string   `json:"link"`
	ThumbnailID  *uint64  `json:"thumbnailId"`
	RemoveUpload bool     `json:"removeUpload"`
	PlaylistIDs  []uint64 `json:"playlistIds"`
}

type Thumbnail struct {
	Revision uint64       `json:"revision" gorm:"default:1"`
	ID       uint64       `json:"id"`
	Title    string       `json:"title"`
	Crop     bool         `json:"crop"`
	AssetID  uint64       `json:"assetId"`
	Asset    *media.Asset `json:"asset,omitempty"`
}

type ThumbnailInput struct {
	Revision uint64 `json:"revision"`
	Title    string `json:"title"`
	Crop     bool   `json:"crop"`
}

type Playlist struct {
	Revision uint64 `json:"revision" gorm:"default:1"`
	ID       uint64 `json:"id"`
	SiteID   uint64 `json:"siteId"`
	Name     string `json:"name"`
	Visible  bool   `json:"visible"`
	Position int    `json:"position"`
	Films    []Film `json:"films" gorm:"many2many:playlist_films"`
}

type PlaylistInput struct {
	Revision uint64   `json:"revision"`
	Name     string   `json:"name"`
	Visible  bool     `json:"visible"`
	Position int      `json:"position"`
	FilmIDs  []uint64 `json:"filmIds"`
}

type Trailer struct {
	Revision uint64       `json:"revision" gorm:"default:1"`
	ID       uint64       `json:"id"`
	SiteID   uint64       `json:"siteId"`
	FilmID   uint64       `json:"filmId"`
	Film     *Film        `json:"film,omitempty"`
	AssetID  uint64       `json:"assetId"`
	Asset    *media.Asset `json:"asset,omitempty"`
	Featured bool         `json:"featured"`
}

type TrailerInput struct {
	Revision uint64 `json:"revision"`
	FilmID   uint64 `json:"filmId"`
	Featured bool   `json:"featured"`
}
