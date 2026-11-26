package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cinefolio/internal/auth"
	"cinefolio/internal/catalog"
	"cinefolio/internal/media"
	"cinefolio/internal/store"
)

type adminHandler = func(http.ResponseWriter, *http.Request, auth.User)

func (s *server) catalogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/public/portfolio", func(w http.ResponseWriter, r *http.Request) {
		p, err := s.catalog.Public(r.Context())
		if err != nil {
			fail(w, err)
			return
		}
		respond(w, http.StatusOK, p)
	})
	mux.HandleFunc("GET /media/{id}", s.serveMedia(false))
	mux.HandleFunc("GET /api/admin/media/{id}", s.protected(func(w http.ResponseWriter, r *http.Request, _ auth.User) { s.serveMedia(true)(w, r) }))

	mux.HandleFunc("GET /api/admin/films", s.protected(list(s.catalog.Films)))
	mux.HandleFunc("GET /api/admin/playlists", s.protected(list(s.catalog.Playlists)))
	mux.HandleFunc("GET /api/admin/thumbnails", s.protected(list(s.catalog.Thumbnails)))
	mux.HandleFunc("GET /api/admin/trailers", s.protected(list(s.catalog.Trailers)))
	s.registerMediaEditor(mux, "films", editMedia(s, s.catalog.SaveFilm, media.Film), s.catalog.DeleteFilm)
	s.registerMediaEditor(mux, "thumbnails", editMedia(s, s.catalog.SaveThumbnail, media.Thumbnail), s.catalog.DeleteThumbnail)
	s.registerMediaEditor(mux, "trailers", editMedia(s, s.catalog.SaveTrailer, media.Trailer), s.catalog.DeleteTrailer)
	mux.HandleFunc("POST /api/admin/playlists", s.protected(s.savePlaylist))
	mux.HandleFunc("PUT /api/admin/playlists/{id}", s.protected(s.savePlaylist))
	mux.HandleFunc("DELETE /api/admin/playlists/{id}", s.protected(remove(s, func(ctx context.Context, id uint64) (string, error) { return "", s.catalog.DeletePlaylist(ctx, id) })))
	mux.HandleFunc("PUT /api/admin/playlists/{id}/films/{film}", s.protected(s.membership))
	mux.HandleFunc("DELETE /api/admin/playlists/{id}/films/{film}", s.protected(s.membership))
	mux.HandleFunc("POST /api/admin/thumbnails/enable-crop", s.protected(s.enableCrop))
}

func idParam(r *http.Request, name string) (uint64, error) {
	id, err := strconv.ParseUint(r.PathValue(name), 10, 63)
	if err != nil || id == 0 {
		return 0, store.ErrNotFound
	}
	return id, nil
}

func recordID(r *http.Request) (uint64, error) {
	if r.Method == http.MethodPost {
		return 0, nil
	}
	return idParam(r, "id")
}

func list[T any](find func(context.Context, catalog.Filter) (catalog.Page[T], error)) adminHandler {
	return func(w http.ResponseWriter, r *http.Request, _ auth.User) {
		q := r.URL.Query()
		filter := catalog.Filter{Query: strings.TrimSpace(q.Get("q")), Limit: 50}
		fields := []struct {
			field    string
			min, max int
			target   *int
		}{
			{"offset", 0, 1000000, &filter.Offset},
			{"limit", 1, 100, &filter.Limit},
		}
		for _, field := range fields {
			if raw := q.Get(field.field); raw != "" {
				n, err := strconv.Atoi(raw)
				if err != nil || n < field.min || n > field.max {
					fail(w, store.ValidationError("invalid "+field.field))
					return
				}
				*field.target = n
			}
		}
		items, err := find(r.Context(), filter)
		if err != nil {
			fail(w, err)
			return
		}
		respond(w, http.StatusOK, items)
	}
}

func (s *server) registerMediaEditor(mux *http.ServeMux, name string, save adminHandler, deleteRecord func(context.Context, uint64) (string, error)) {
	base := "/api/admin/" + name
	mux.HandleFunc("POST "+base, s.protected(save))
	mux.HandleFunc("PUT "+base+"/{id}", s.protected(save))
	mux.HandleFunc("DELETE "+base+"/{id}", s.protected(remove(s, deleteRecord)))
}

func editMedia[I, R any](s *server, save func(context.Context, uint64, I, *media.Asset) (R, string, error), kind media.Kind) adminHandler {
	return func(w http.ResponseWriter, r *http.Request, _ auth.User) {
		id, err := recordID(r)
		if err != nil {
			fail(w, err)
			return
		}
		var input I
		var upload *media.Asset
		committed := false
		defer func() {
			if !committed && upload != nil {
				s.cleanup(upload.StorageKey)
			}
		}()
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			r.Body = http.MaxBytesReader(w, r.Body, kind.Limit()+(256<<10))
			reader, err := r.MultipartReader()
			if err != nil {
				fail(w, store.ValidationError("invalid multipart request"))
				return
			}
			hasInput := false
			for {
				part, err := reader.NextPart()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					fail(w, store.ValidationError("invalid or oversized upload"))
					return
				}
				switch part.FormName() {
				case "input":
					if hasInput {
						fail(w, store.ValidationError("duplicate input"))
						return
					}
					hasInput = true
					decoder := json.NewDecoder(io.LimitReader(part, 128<<10))
					decoder.DisallowUnknownFields()
					if decoder.Decode(&input) != nil || !errors.Is(decoder.Decode(new(any)), io.EOF) {
						fail(w, store.ValidationError("invalid input"))
						return
					}
				case "file":
					if upload != nil || part.FileName() == "" {
						fail(w, store.ValidationError("expected one media file"))
						return
					}
					staged, err := s.media.Stage(kind, part.FileName(), part)
					if err != nil {
						fail(w, err)
						return
					}
					upload = &staged
				default:
					fail(w, store.ValidationError("unexpected upload field"))
					return
				}
			}
			if !hasInput {
				fail(w, store.ValidationError("missing input"))
				return
			}
		} else if err := decode(w, r, &input); err != nil {
			fail(w, err)
			return
		}
		record, obsolete, err := save(r.Context(), id, input, upload)
		if err != nil {
			fail(w, err)
			return
		}
		committed = true
		s.cleanup(obsolete)
		status := http.StatusOK
		if id == 0 {
			status = http.StatusCreated
		}
		respond(w, status, record)
	}
}

func remove(s *server, deleteRecord func(context.Context, uint64) (string, error)) adminHandler {
	return func(w http.ResponseWriter, r *http.Request, _ auth.User) {
		id, err := idParam(r, "id")
		if err != nil {
			fail(w, err)
			return
		}
		obsolete, err := deleteRecord(r.Context(), id)
		if err != nil {
			fail(w, err)
			return
		}
		s.cleanup(obsolete)
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) cleanup(key string) {
	if key != "" {
		if err := s.media.Remove(key); err != nil {
			slog.Error("remove unreferenced media", "key", key, "error", err)
		}
	}
}

func (s *server) savePlaylist(w http.ResponseWriter, r *http.Request, _ auth.User) {
	id, err := recordID(r)
	if err != nil {
		fail(w, err)
		return
	}
	var input catalog.PlaylistInput
	if err := decode(w, r, &input); err != nil {
		fail(w, err)
		return
	}
	record, err := s.catalog.SavePlaylist(r.Context(), id, input)
	if err != nil {
		fail(w, err)
		return
	}
	status := http.StatusOK
	if id == 0 {
		status = http.StatusCreated
	}
	respond(w, status, record)
}

func (s *server) membership(w http.ResponseWriter, r *http.Request, _ auth.User) {
	id, err := idParam(r, "id")
	if err != nil {
		fail(w, err)
		return
	}
	film, err := idParam(r, "film")
	if err != nil {
		fail(w, err)
		return
	}
	if err := s.catalog.SetMembership(r.Context(), id, film, r.Method == http.MethodPut); err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) enableCrop(w http.ResponseWriter, r *http.Request, _ auth.User) {
	var input struct {
		IDs []uint64 `json:"ids"`
	}
	if err := decode(w, r, &input); err != nil {
		fail(w, err)
		return
	}
	if err := s.catalog.EnableCrop(r.Context(), input.IDs); err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) serveMedia(admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := idParam(r, "id")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		asset, err := s.catalog.Asset(r.Context(), id, admin)
		if catalog.Missing(err) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			fail(w, err)
			return
		}
		file, err := s.media.Open(asset.StorageKey)
		if err != nil {
			slog.Error("missing media", "id", id, "error", err)
			http.NotFound(w, r)
			return
		}
		defer file.Close()
		w.Header().Set("Content-Type", asset.ContentType)
		if !admin {
			w.Header().Set("Cache-Control", "no-cache")
		}
		w.Header().Set("ETag", `"`+asset.StorageKey+`"`)
		http.ServeContent(w, r, asset.Filename, time.Time{}, file)
	}
}
