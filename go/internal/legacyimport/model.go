package legacyimport

type source struct {
	Sites       []legacySite       `json:"sites"`
	Films       []legacyFilm       `json:"files"`
	Playlists   []legacyPlaylist   `json:"playlists"`
	Memberships []legacyMembership `json:"filetoplaylists"`
	Thumbnails  []legacyThumbnail  `json:"thumbnails"`
	Trailers    []legacyTrailer    `json:"trailers"`
}

type legacySite struct {
	ID   string `json:"_id"`
	Name string
}

type legacyFilm struct {
	ID          string `json:"_id"`
	Filepath    string
	Link        string
	Title       string
	Description string
	Published   string
	Thumbnail   string
}

type legacyPlaylist struct {
	ID      string `json:"_id"`
	Name    string
	Site    string
	Visible *bool
	Order   int
}

type legacyMembership struct {
	ID         string `json:"_id"`
	FileID     string `json:"fileId"`
	PlaylistID string `json:"playlistId"`
}

type legacyThumbnail struct {
	ID       string `json:"_id"`
	Title    string
	Filepath string
	Crop     bool
}

type legacyTrailer struct {
	ID       string `json:"_id"`
	Film     string
	Filepath string
	Visible  bool
}
