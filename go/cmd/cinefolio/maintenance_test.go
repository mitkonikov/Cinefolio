package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"cinefolio/internal/auth"
	"cinefolio/internal/store"
)

func TestBackupsAndRestoresDoNotReactivateSessions(t *testing.T) {
	root := t.TempDir()
	data := filepath.Join(root, "live")
	db, err := store.Open(filepath.Join(data, "cinefolio.db"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	if err := os.Mkdir(filepath.Join(data, "media"), 0700); err != nil {
		t.Fatal(err)
	}
	ctx := t.Context()
	repo := auth.New(db)
	if err := repo.SetPassword(ctx, "owner@example.com", "testing password"); err != nil {
		t.Fatal(err)
	}
	token, _, err := repo.Login(ctx, "owner@example.com", "testing password")
	if err != nil {
		t.Fatal(err)
	}

	paired := filepath.Join(root, "paired")
	if err := backupData(data, paired); err != nil {
		t.Fatal(err)
	}
	copyDB := filepath.Join(root, "database-only.db")
	if err := store.BackupFile(filepath.Join(data, "cinefolio.db"), copyDB); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(root, "legacy.db")
	if err := db.Exec("VACUUM INTO ?", legacy).Error; err != nil {
		t.Fatal(err)
	}
	if err := store.PurgeSessions(legacy); err != nil {
		t.Fatal(err)
	}
	for _, filename := range []string{filepath.Join(paired, "cinefolio.db"), copyDB, legacy} {
		restored, err := store.ReadOnly(filename)
		if err != nil {
			t.Fatal(err)
		}
		restoredSQL, err := restored.DB()
		if err != nil {
			t.Fatal(err)
		}
		_, err = auth.New(restored).Current(ctx, token)
		restoredSQL.Close()
		if !errors.Is(err, auth.ErrUnauthorized) {
			t.Errorf("%s accepted saved session: %v", filename, err)
		}
	}
	if _, err := repo.Current(ctx, token); err != nil {
		t.Fatalf("backup changed live session: %v", err)
	}
}
