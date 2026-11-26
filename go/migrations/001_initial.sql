CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);
CREATE TABLE sessions (
    token_hash TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at INTEGER NOT NULL
);
CREATE INDEX sessions_user ON sessions(user_id);
CREATE INDEX sessions_expiry ON sessions(expires_at);

CREATE TABLE sites (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE
);
INSERT INTO sites (name, slug) VALUES ('Cinefolio', 'cinefolio');

CREATE TABLE assets (
    id INTEGER PRIMARY KEY,
    storage_key TEXT NOT NULL UNIQUE,
    filename TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size INTEGER NOT NULL CHECK (size > 0)
);
CREATE TABLE thumbnails (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL DEFAULT '',
    crop BOOLEAN NOT NULL DEFAULT 0,
    asset_id INTEGER NOT NULL UNIQUE REFERENCES assets(id)
);
CREATE TABLE films (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL CHECK (length(trim(title)) > 0),
    description TEXT NOT NULL DEFAULT '',
    published TEXT NOT NULL,
    link TEXT NOT NULL DEFAULT '',
    asset_id INTEGER UNIQUE REFERENCES assets(id),
    thumbnail_id INTEGER REFERENCES thumbnails(id),
    CHECK (length(link) > 0 OR asset_id IS NOT NULL)
);
CREATE INDEX films_thumbnail ON films(thumbnail_id);
CREATE TABLE playlists (
    id INTEGER PRIMARY KEY,
    site_id INTEGER NOT NULL REFERENCES sites(id),
    name TEXT NOT NULL,
    visible BOOLEAN NOT NULL DEFAULT 1,
    position INTEGER NOT NULL,
    UNIQUE (site_id, name)
);
CREATE INDEX playlists_public ON playlists(site_id, visible, position, id);
CREATE TABLE playlist_films (
    playlist_id INTEGER NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    film_id INTEGER NOT NULL REFERENCES films(id) ON DELETE RESTRICT,
    PRIMARY KEY (playlist_id, film_id)
);
CREATE INDEX playlist_films_film ON playlist_films(film_id);
CREATE TABLE trailers (
    id INTEGER PRIMARY KEY,
    site_id INTEGER NOT NULL REFERENCES sites(id),
    film_id INTEGER NOT NULL REFERENCES films(id) ON DELETE RESTRICT,
    asset_id INTEGER NOT NULL UNIQUE REFERENCES assets(id),
    featured BOOLEAN NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX one_featured_trailer ON trailers(site_id) WHERE featured;
CREATE INDEX trailers_film ON trailers(film_id);

CREATE TABLE legacy_ids (
    kind TEXT NOT NULL,
    legacy_id TEXT NOT NULL,
    record_id INTEGER NOT NULL,
    PRIMARY KEY (kind, legacy_id)
);
