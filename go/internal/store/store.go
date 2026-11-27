package store

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"

	"cinefolio/migrations"
)

var (
	ErrNotFound  = errors.New("resource not found")
	ErrForbidden = errors.New("access forbidden")
	ErrConflict  = errors.New("the record changed; reopen it and try again")
)

type ValidationError string

func (e ValidationError) Error() string { return string(e) }

func Open(path string) (*gorm.DB, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(absolute), 0700); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(absolute, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}

	db, err := connect(absolute, "_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}
	if err := migrations.Apply(db); err != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		return nil, err
	}
	return db, nil
}

func connect(absolute, query string) (*gorm.DB, error) {
	uriPath := filepath.ToSlash(absolute)
	if runtime.GOOS == "windows" && len(uriPath) >= 2 && uriPath[1] == ':' {
		uriPath = "/" + uriPath
	}
	dsn := url.URL{Scheme: "file", Path: uriPath, RawQuery: query}
	db, err := gorm.Open(sqlite.Dialector{DriverName: "sqlite", DSN: dsn.String()}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open SQLite: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// Serialize transactions within each process.
	sqlDB.SetMaxOpenConns(1)
	return db, nil
}
