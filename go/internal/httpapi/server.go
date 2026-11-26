package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"cinefolio/internal/auth"
	"cinefolio/internal/catalog"
	"cinefolio/internal/media"
	"cinefolio/internal/store"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

const cookieName = "cinefolio_session"

type Config struct {
	WebDir       string
	Origin       string
	SecureCookie bool
}

type server struct {
	auth    *auth.Repository
	catalog *catalog.Repository
	media   *media.Storage
	config  Config
	logins  *rate.Limiter
}

func New(db *gorm.DB, storage *media.Storage, site string, config Config) http.Handler {
	s := &server{auth: auth.New(db), catalog: catalog.New(db, site), media: storage, config: config, logins: rate.NewLimiter(rate.Every(2*time.Second), 8)}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err == nil {
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()
			err = sqlDB.PingContext(ctx)
		}
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /api/admin/login", s.login)
	mux.HandleFunc("GET /api/admin/session", s.protected(func(w http.ResponseWriter, r *http.Request, u auth.User) { respond(w, http.StatusOK, u) }))
	mux.HandleFunc("POST /api/admin/logout", s.protected(s.logout))
	s.catalogRoutes(mux)
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) { fail(w, store.ErrNotFound) })
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) { fail(w, store.ErrNotFound) })
	mux.HandleFunc("/media/", http.NotFound)
	mux.Handle("/", static(config.WebDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Cache-Control", "no-store")
			if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Header.Get("Origin") != s.config.Origin {
				fail(w, store.ErrForbidden)
				return
			}
		}
		mux.ServeHTTP(w, r)
	})
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	if !s.logins.Allow() {
		respond(w, http.StatusTooManyRequests, map[string]string{"error": "try again shortly"})
		return
	}
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decode(w, r, &input); err != nil {
		fail(w, err)
		return
	}
	token, user, err := s.auth.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		fail(w, err)
		return
	}
	s.cookie(w, token, int(auth.Lifetime/time.Second))
	respond(w, http.StatusOK, user)
}

func (s *server) logout(w http.ResponseWriter, r *http.Request, _ auth.User) {
	if err := s.auth.Logout(r.Context(), requestToken(r)); err != nil {
		fail(w, err)
		return
	}
	s.cookie(w, "", -1)
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) cookie(w http.ResponseWriter, token string, age int) {
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: token, Path: "/", HttpOnly: true, Secure: s.config.SecureCookie, SameSite: http.SameSiteStrictMode, MaxAge: age})
}

func requestToken(r *http.Request) string {
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		return cookie.Value
	}
	return ""
}

func (s *server) protected(next func(http.ResponseWriter, *http.Request, auth.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := s.auth.Current(r.Context(), requestToken(r))
		if err != nil {
			fail(w, err)
			return
		}
		next(w, r, user)
	}
}

func decode(w http.ResponseWriter, r *http.Request, target any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 128<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return store.ValidationError("invalid JSON request")
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return store.ValidationError("expected one JSON value")
	}
	return nil
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Debug("write response", "error", err)
	}
}

func fail(w http.ResponseWriter, err error) {
	var invalid store.ValidationError
	status, message := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.As(err, &invalid):
		status, message = http.StatusBadRequest, invalid.Error()
	case errors.Is(err, auth.ErrUnauthorized):
		status, message = http.StatusUnauthorized, err.Error()
	case errors.Is(err, auth.ErrHashBusy):
		status, message = http.StatusServiceUnavailable, err.Error()
	case errors.Is(err, store.ErrForbidden):
		status, message = http.StatusForbidden, err.Error()
	case catalog.Missing(err):
		status, message = http.StatusNotFound, "record not found"
	case errors.Is(err, store.ErrConflict), sqliteConflict(err):
		status, message = http.StatusConflict, "the record is in use or conflicts with existing data"
	default:
		slog.Error("request failed", "error", err)
	}
	respond(w, status, map[string]string{"error": message})
}

func sqliteConflict(err error) bool {
	var sqliteErr interface{ Code() int }
	// SQLite extended constraint codes for unique, primary key, and referenced rows.
	if !errors.As(err, &sqliteErr) {
		return false
	}
	switch sqliteErr.Code() {
	case 1555, 1811, 2067, 787:
		return true
	default:
		return false
	}
}
