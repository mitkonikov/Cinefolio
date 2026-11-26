package media

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cinefolio/internal/store"
)

type Asset struct {
	ID          uint64 `json:"id"`
	StorageKey  string `json:"-"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

func (a Asset) URL() string { return fmt.Sprintf("/media/%d", a.ID) }

type Kind string

const (
	Film      Kind = "film"
	Trailer   Kind = "trailer"
	Thumbnail Kind = "thumbnail"
)

func (k Kind) Limit() int64 {
	switch k {
	case Film:
		return 500 << 20
	case Trailer:
		return 15 << 20
	case Thumbnail:
		return 512 << 10
	default:
		panic("unknown media kind")
	}
}

func (k Kind) accepts(contentType string) bool {
	switch k {
	case Thumbnail:
		return contentType == "image/jpeg"
	case Trailer:
		return contentType == "video/mp4"
	case Film:
		switch contentType {
		case "video/mp4", "video/avi", "video/webm", "video/mpeg", "application/ogg":
			return true
		}
	}
	return false
}

type Storage struct{ dir string }

func New(dir string) (*Storage, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &Storage{dir: dir}, nil
}

func (s *Storage) Stage(kind Kind, filename string, source io.Reader) (Asset, error) {
	limit := kind.Limit()
	var header [512]byte
	n, err := io.ReadFull(source, header[:])
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return Asset{}, err
	}
	contentType := http.DetectContentType(header[:n])
	// MPEG program streams are not identified by DetectContentType.
	if n >= 4 && string(header[:4]) == "\x00\x00\x01\xba" {
		contentType = "video/mpeg"
	}
	if !kind.accepts(contentType) {
		return Asset{}, store.ValidationError("unsupported media type")
	}
	if contentType == "application/ogg" {
		contentType = "video/ogg"
	}
	filename = filepath.Base(strings.ReplaceAll(filename, "\\", "/"))
	if len(filename) > 255 || filename == "." || filename == "" || strings.ContainsAny(filename, "\r\n\x00") {
		return Asset{}, store.ValidationError("invalid media filename")
	}

	file, err := os.CreateTemp(s.dir, ".upload-*")
	if err != nil {
		return Asset{}, err
	}
	defer os.Remove(file.Name())
	written, err := file.Write(header[:n])
	if err == nil {
		var extra int64
		extra, err = io.Copy(file, io.LimitReader(source, limit-int64(written)+1))
		written += int(extra)
	}
	if err == nil && int64(written) > limit {
		err = store.ValidationError("media exceeds its upload limit")
	}
	if err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err != nil {
		return Asset{}, err
	}
	if closeErr != nil {
		return Asset{}, closeErr
	}
	var random [24]byte
	if _, err := rand.Read(random[:]); err != nil {
		return Asset{}, err
	}
	key := hex.EncodeToString(random[:])
	if err := os.Rename(file.Name(), filepath.Join(s.dir, key)); err != nil {
		return Asset{}, err
	}
	return Asset{StorageKey: key, Filename: filename, ContentType: contentType, Size: int64(written)}, nil
}

func (s *Storage) Open(key string) (*os.File, error) {
	path, err := s.path(key)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func (s *Storage) Remove(key string) error {
	path, err := s.path(key)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *Storage) path(key string) (string, error) {
	if _, err := hex.DecodeString(key); err != nil || len(key) != 48 {
		return "", fmt.Errorf("invalid storage key")
	}
	return filepath.Join(s.dir, key), nil
}
