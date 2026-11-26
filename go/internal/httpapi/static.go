package httpapi

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
)

func static(directory string) http.Handler {
	files := http.Dir(directory)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		name := path.Clean(r.URL.Path)
		file, err := files.Open(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		info, err := file.Stat()
		if err != nil {
			file.Close()
			http.NotFound(w, r)
			return
		}
		if info.IsDir() {
			file.Close()
			if !strings.HasSuffix(r.URL.Path, "/") {
				http.Redirect(w, r, r.URL.Path+"/", http.StatusPermanentRedirect)
				return
			}
			file, err = files.Open(path.Join(name, "index.html"))
			if errors.Is(err, os.ErrNotExist) || err != nil {
				http.NotFound(w, r)
				return
			}
			name = path.Join(name, "index.html")
			info, err = file.Stat()
			if err != nil || info.IsDir() {
				file.Close()
				http.NotFound(w, r)
				return
			}
		}
		defer file.Close()
		w.Header().Set("Cache-Control", "no-cache")
		if strings.HasPrefix(name, "/_app/immutable/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		http.ServeContent(w, r, path.Base(name), info.ModTime(), file)
	})
}
