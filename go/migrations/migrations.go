package migrations

import (
	"embed"
	"fmt"
	"io/fs"

	"gorm.io/gorm"
)

//go:embed *.sql
var scripts embed.FS

func Apply(db *gorm.DB) error {
	entries, err := fs.Glob(scripts, "*.sql")
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY)").Error; err != nil {
			return err
		}
		for _, name := range entries {
			var count int64
			if err := tx.Table("schema_migrations").Where("version = ?", name).Count(&count).Error; err != nil {
				return err
			}
			if count != 0 {
				continue
			}
			content, err := fs.ReadFile(scripts, name)
			if err != nil {
				return err
			}
			if err := tx.Exec(string(content)).Error; err != nil {
				return fmt.Errorf("apply migration %s: %w", name, err)
			}
			if err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
