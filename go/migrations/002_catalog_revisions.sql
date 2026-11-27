ALTER TABLE films ADD COLUMN revision INTEGER NOT NULL DEFAULT 1 CHECK (revision > 0);

CREATE TRIGGER films_revise AFTER UPDATE OF title, description, published, link, asset_id, thumbnail_id ON films
BEGIN
    UPDATE films SET revision = revision + 1 WHERE id = NEW.id;
END;

ALTER TABLE playlists ADD COLUMN revision INTEGER NOT NULL DEFAULT 1 CHECK (revision > 0);

CREATE TRIGGER playlists_revise AFTER UPDATE OF name, visible, position ON playlists
BEGIN
    UPDATE playlists SET revision = revision + 1 WHERE id = NEW.id;
END;

ALTER TABLE thumbnails ADD COLUMN revision INTEGER NOT NULL DEFAULT 1 CHECK (revision > 0);

CREATE TRIGGER thumbnails_revise AFTER UPDATE OF title, crop, asset_id ON thumbnails
BEGIN
    UPDATE thumbnails SET revision = revision + 1 WHERE id = NEW.id;
END;

ALTER TABLE trailers ADD COLUMN revision INTEGER NOT NULL DEFAULT 1 CHECK (revision > 0);

CREATE TRIGGER trailers_revise AFTER UPDATE OF film_id, asset_id, featured ON trailers
BEGIN
    UPDATE trailers SET revision = revision + 1 WHERE id = NEW.id;
END;

CREATE TRIGGER playlist_films_revise_insert AFTER INSERT ON playlist_films
BEGIN
    UPDATE films SET revision = revision + 1 WHERE id = NEW.film_id;
    UPDATE playlists SET revision = revision + 1 WHERE id = NEW.playlist_id;
END;

CREATE TRIGGER playlist_films_revise_delete AFTER DELETE ON playlist_films
BEGIN
    UPDATE films SET revision = revision + 1 WHERE id = OLD.film_id;
    UPDATE playlists SET revision = revision + 1 WHERE id = OLD.playlist_id;
END;
