package store

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// BackupFile snapshots an existing database without running application migrations.
func BackupFile(source, output string) error {
	absolute, err := filepath.Abs(source)
	if err != nil {
		return err
	}
	db, err := connect(absolute, "mode=ro&_pragma=busy_timeout(5000)")
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	return Backup(db, output)
}

func Backup(db *gorm.DB, output string) error {
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
