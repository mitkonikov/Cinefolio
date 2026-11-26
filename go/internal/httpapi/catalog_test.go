package httpapi

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"cinefolio/internal/auth"
	"cinefolio/internal/catalog"
	"cinefolio/internal/media"
	"cinefolio/internal/store"
)

func TestCatalogAndPublicMedia(t *testing.T) {
	dir := t.TempDir()
	db, err := store.Open(filepath.Join(dir, "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	storage, err := media.New(filepath.Join(dir, "media"))
	if err != nil {
		t.Fatal(err)
	}
	const origin = "http://127.0.0.1:5002"
	handler := New(db, storage, "cinefolio", Config{WebDir: dir, Origin: origin})
	ctx := t.Context()
	if err := auth.New(db).SetPassword(ctx, "owner@example.com", "a testing password"); err != nil {
		t.Fatal(err)
	}
	token, _, err := auth.New(db).Login(ctx, "owner@example.com", "a testing password")
	if err != nil {
		t.Fatal(err)
	}
	do := func(method, path, body, contentType string, admin bool, headers map[string]string, status int) *httptest.ResponseRecorder {
		t.Helper()
		request := httptest.NewRequest(method, path, strings.NewReader(body))
		request.Header.Set("Content-Type", contentType)
		request.Header.Set("Origin", origin)
		for name, value := range headers {
			request.Header.Set(name, value)
		}
		if admin {
			request.AddCookie(&http.Cookie{Name: cookieName, Value: token})
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, request)
		if rec.Code != status {
			t.Fatalf("%s %s: status %d want %d: %s", method, path, rec.Code, status, rec.Body.String())
		}
		return rec
	}
	do("GET", "/api/admin/films", "", "", false, nil, 401)
	do("POST", "/api/admin/films", "", "", true, map[string]string{"Origin": "https://other.example"}, 403)
	upload := func(path, input, filename string, content []byte, status int) *httptest.ResponseRecorder {
		t.Helper()
		var body bytes.Buffer
		form := multipart.NewWriter(&body)
		if err := form.WriteField("input", input); err != nil {
			t.Fatal(err)
		}
		part, err := form.CreateFormFile("file", filename)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := form.Close(); err != nil {
			t.Fatal(err)
		}
		return do("POST", path, body.String(), form.FormDataContentType(), true, nil, status)
	}
	index := func(rec *httptest.ResponseRecorder) string {
		t.Helper()
		var result struct {
			ID      uint64
			AssetID uint64
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return strconv.FormatUint(result.ID, 10)
	}
	video := []byte{0, 0, 0, 24, 'f', 't', 'y', 'p', 'm', 'p', '4', '2', 0, 0, 0, 0, 'm', 'p', '4', '2', 'i', 's', 'o', 'm'}
	film := index(upload("/api/admin/films", `{"title":"First film","published":"2024-10-01"}`, "test.mp4", video, 201))
	var created catalog.Trailer
	trailer := upload("/api/admin/trailers", `{"filmId":`+film+`,"featured":true}`, "trailer.mp4", video, 201)
	if err := json.Unmarshal(trailer.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	other := index(upload("/api/admin/trailers", `{"filmId":`+film+`,"featured":true}`, "second.mp4", video, 201))
	do("GET", "/media/"+strconv.FormatUint(created.AssetID, 10), "", "", false, nil, 404)
	list := index(do("POST", "/api/admin/playlists", `{"name":"Shorts","visible":true,"position":2,"filmIds":[`+film+`]}`, "application/json", true, nil, 201))
	var portfolio catalog.Portfolio
	rec := do("GET", "/api/public/portfolio", "", "", false, nil, 200)
	if err := json.Unmarshal(rec.Body.Bytes(), &portfolio); err != nil {
		t.Fatal(err)
	}
	if len(portfolio.Playlists) != 1 || len(portfolio.Playlists[0].Films) != 1 || portfolio.Trailer == nil {
		t.Fatalf("portfolio = %s", rec.Body.String())
	}
	url := portfolio.Playlists[0].Films[0].VideoURL
	rec = do("GET", url, "", "", false, map[string]string{"Range": "bytes=0-3"}, 206)
	if rec.Body.Len() != 4 {
		t.Fatalf("range length %d", rec.Body.Len())
	}
	do("DELETE", "/api/admin/films/"+film, "", "", true, nil, 409)
	do("DELETE", "/api/admin/trailers/"+other, "", "", true, nil, 204)
	do("PUT", "/api/admin/playlists/"+list, `{"name":"Shorts","visible":false,"position":2}`, "application/json", true, nil, 200)
	do("GET", url, "", "", false, nil, 404)
	do("GET", strings.Replace(url, "/media/", "/api/admin/media/", 1), "", "", true, nil, 200)
	do("GET", "/missing.js", "", "", false, nil, 404)
}
