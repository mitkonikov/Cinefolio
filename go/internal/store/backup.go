package store

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// BackupFile snapshots an existing database without running application migrations.
func BackupFile(source, output string) error {
	scratch, err := os.MkdirTemp(filepath.Dir(output), ".cinefolio-snapshot-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(scratch)
	db, err := ReadOnly(source)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	staged := filepath.Join(scratch, "snapshot.db")
	if err := snapshot(db, staged); err != nil {
		return err
	}
	if err := PurgeSessions(staged); err != nil {
		return fmt.Errorf("sanitize backup sessions: %w", err)
	}
	// Publish only after sanitizing and refuse to overwrite an existing backup.
	return os.Link(staged, output)
}

// PurgeSessions opens an existing database without applying migrations.
func PurgeSessions(source string) error {
	absolute, err := filepath.Abs(source)
	if err != nil {
		return err
	}
	db, err := connect(absolute, "mode=rw&_pragma=busy_timeout(5000)&_pragma=secure_delete(ON)")
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	return db.Exec("DELETE FROM sessions").Error
}

func ReadOnly(source string) (*gorm.DB, error) {
	absolute, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	return connect(absolute, "mode=ro&_pragma=busy_timeout(5000)")
}

func snapshot(db *gorm.DB, output string) error {
	path, err := filepath.Abs(output)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("create backup destination: %w", err)
	}
	if err := file.Close(); err != nil {
		os.Remove(path)
		return err
	}
	// VACUUM INTO takes a consistent snapshot, including committed WAL data.
	if err := db.Exec("VACUUM INTO ?", path).Error; err != nil {
		os.Remove(path)
		return fmt.Errorf("write backup: %w", err)
	}
	return nil
}
