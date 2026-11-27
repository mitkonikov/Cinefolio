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
	handler := New(db, storage, "cinefolio", Config{WebDir: dir, Origin: origin, TrustLoopbackProxy: true})
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
		request.RemoteAddr = "127.0.0.1:12345"
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
	attacker := map[string]string{"X-Forwarded-For": "198.51.100.10"}
	for range loginBurst {
		do("POST", "/api/admin/login", "{", "application/json", false, attacker, 400)
	}
	do("POST", "/api/admin/login", "{", "application/json", false, attacker, 429)
	do("POST", "/api/admin/login", `{"email":"owner@example.com","password":"a testing password"}`,
		"application/json", false, map[string]string{"X-Forwarded-For": "198.51.100.11"}, 200)
	do("POST", "/api/admin/films", `{"title":"Invalid","published":"2024-10-01","link":"https://example.com:99999/video"}`,
		"application/json", true, nil, 400)
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
	revision := func(resource, id string) string {
		t.Helper()
		rec := do("GET", "/api/admin/"+resource, "", "", true, nil, 200)
		var page struct {
			Items []struct{ ID, Revision uint64 }
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
			t.Fatal(err)
		}
		for _, row := range page.Items {
			if strconv.FormatUint(row.ID, 10) == id {
				return strconv.FormatUint(row.Revision, 10)
			}
		}
		t.Fatalf("missing %s %s", resource, id)
		return ""
	}
	staleFilm := revision("films", film)
	var created catalog.Trailer
	trailer := upload("/api/admin/trailers", `{"filmId":`+film+`,"featured":true}`, "trailer.mp4", video, 201)
	if err := json.Unmarshal(trailer.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	other := index(upload("/api/admin/trailers", `{"filmId":`+film+`,"featured":true}`, "second.mp4", video, 201))
	do("PUT", "/api/admin/trailers/"+strconv.FormatUint(created.ID, 10),
		`{"revision":`+strconv.FormatUint(created.Revision, 10)+`,"filmId":`+film+`,"featured":true}`,
		"application/json", true, nil, 409)
	do("GET", "/media/"+strconv.FormatUint(created.AssetID, 10), "", "", false, nil, 404)
	list := index(do("POST", "/api/admin/playlists", `{"name":"Shorts","visible":true,"position":2,"filmIds":[`+film+`]}`, "application/json", true, nil, 201))
	stalePlaylist := revision("playlists", list)
	do("PUT", "/api/admin/films/"+film, `{"revision":`+staleFilm+`,"title":"Stale title","published":"2024-10-01"}`,
		"application/json", true, nil, 409)
	updateFilm := func(ids string) {
		do("PUT", "/api/admin/films/"+film, `{"revision":`+revision("films", film)+`,"title":"First film","published":"2024-10-01","playlistIds":`+ids+`}`, "application/json", true, nil, 200)
	}
	updateFilm("[]")
	do("PUT", "/api/admin/playlists/"+list, `{"revision":`+stalePlaylist+`,"name":"Stale name","visible":true,"position":2,"filmIds":[`+film+`]}`,
		"application/json", true, nil, 409)
	updateFilm("[" + list + "]")
	do("DELETE", "/api/admin/playlists/"+list+"?revision="+stalePlaylist, "", "", true, nil, 409)
	do("PUT", "/api/admin/films/"+film, `{"revision":`+revision("films", film)+`,"title":"First film updated","published":"2024-10-01"}`,
		"application/json", true, nil, 200)
	rec := do("GET", "/api/admin/films?q=First", "", "", true, nil, 200)
	var films catalog.Page[catalog.Film]
	if err := json.Unmarshal(rec.Body.Bytes(), &films); err != nil || len(films.Items) != 1 || len(films.Items[0].Playlists) != 1 {
		t.Fatalf("admin list: %s (%v)", rec.Body.String(), err)
	}
	var portfolio catalog.Portfolio
	rec = do("GET", "/api/public/portfolio", "", "", false, nil, 200)
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
	do("DELETE", "/api/admin/films/"+film+"?revision="+revision("films", film), "", "", true, nil, 409)
	do("DELETE", "/api/admin/trailers/"+other+"?revision="+revision("trailers", other), "", "", true, nil, 204)
	do("PUT", "/api/admin/playlists/"+list, `{"revision":`+revision("playlists", list)+`,"name":"Shorts","visible":false,"position":2}`, "application/json", true, nil, 200)
	do("GET", url, "", "", false, nil, 404)
	do("GET", strings.Replace(url, "/media/", "/api/admin/media/", 1), "", "", true, nil, 200)
	do("GET", "/missing.js", "", "", false, nil, 404)
	before := revision("films", film)
	do("DELETE", "/api/admin/playlists/"+list+"?revision="+revision("playlists", list), "", "", true, nil, 204)
	if before == revision("films", film) {
		t.Fatal("deleting a playlist did not revise its films")
	}

	jpeg := []byte{0xff, 0xd8, 0xff, 0xd9}
	thumb := index(upload("/api/admin/thumbnails", `{"title":"Frame"}`, "frame.jpg", jpeg, 201))
	otherThumb := index(upload("/api/admin/thumbnails", `{"title":"Another frame"}`, "other.jpg", jpeg, 201))
	firstRevision := revision("thumbnails", thumb)
	staleThumb := revision("thumbnails", otherThumb)
	do("PUT", "/api/admin/thumbnails/"+otherThumb, `{"revision":`+staleThumb+`,"title":"Renamed frame"}`, "application/json", true, nil, 200)
	crop := `{"records":[{"id":` + thumb + `,"revision":` + firstRevision + `},{"id":` + otherThumb + `,"revision":`
	do("POST", "/api/admin/thumbnails/enable-crop", crop+staleThumb+`}]}`, "application/json", true, nil, 409)
	if revision("thumbnails", thumb) != firstRevision {
		t.Fatal("stale bulk edit applied partially")
	}
	do("POST", "/api/admin/thumbnails/enable-crop", crop+revision("thumbnails", otherThumb)+`}]}`,
		"application/json", true, nil, 204)
	do("PUT", "/api/admin/thumbnails/"+thumb, `{"revision":`+firstRevision+`,"title":"Uncrop"}`, "application/json", true, nil, 409)
}
