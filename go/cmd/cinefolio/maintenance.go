package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"cinefolio/internal/media"
	"cinefolio/internal/store"
)

func maintenance(data string, args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	switch args[0] {
	case "verify-data":
		if len(args) != 1 {
			return errors.New("usage: cinefolio verify-data")
		}
		count, unused, err := verifyData(data)
		if err != nil {
			return err
		}
		fmt.Printf("Verified SQLite integrity and %d media files.\n", count)
		for _, name := range unused {
			fmt.Printf("Unreferenced media: %s\n", name)
		}
		return nil
	default:
		output := flags.String("output", "", "new directory for a verified database and media backup")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if *output == "" || flags.NArg() != 0 {
			return errors.New("usage: cinefolio backup --output DIRECTORY")
		}
		if err := backupData(data, *output); err != nil {
			return err
		}
		fmt.Printf("Verified paired backup: %s\n", *output)
		return nil
	}
}

func assetsAt(data string) ([]media.Asset, error) {
	db, err := store.ReadOnly(filepath.Join(data, "cinefolio.db"))
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()
	var integrity []string
	if err := db.Raw("PRAGMA integrity_check").Scan(&integrity).Error; err != nil {
		return nil, err
	}
	if len(integrity) != 1 || integrity[0] != "ok" {
		return nil, errors.New("SQLite integrity check failed")
	}
	var foreignKeys int
	if err := db.Raw("SELECT count(*) FROM pragma_foreign_key_check").Scan(&foreignKeys).Error; err != nil {
		return nil, err
	}
	if foreignKeys != 0 {
		return nil, fmt.Errorf("SQLite foreign key check found %d violations", foreignKeys)
	}
	var assets []media.Asset
	err = db.Order("id").Find(&assets).Error
	return assets, err
}

func verifyData(data string) (int, []string, error) {
	assets, err := assetsAt(data)
	if err != nil {
		return 0, nil, err
	}
	root := filepath.Join(data, "media")
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, nil, err
	}
	storage, err := media.New(root)
	if err != nil {
		return 0, nil, err
	}
	known := make(map[string]bool, len(assets))
	for _, asset := range assets {
		file, err := storage.Open(asset.StorageKey)
		if err != nil {
			return 0, nil, fmt.Errorf("asset %d: %w", asset.ID, err)
		}
		info, err := file.Stat()
		file.Close()
		if err != nil {
			return 0, nil, err
		}
		if !info.Mode().IsRegular() || info.Size() != asset.Size {
			return 0, nil, fmt.Errorf("asset %d has unexpected size or type", asset.ID)
		}
		known[asset.StorageKey] = true
	}
	var unused []string
	for _, entry := range entries {
		if !known[entry.Name()] {
			unused = append(unused, entry.Name())
		}
	}
	return len(assets), unused, nil
}

func backupData(data, output string) error {
	destination, err := filepath.Abs(output)
	if err != nil {
		return err
	}
	if _, err := os.Lstat(destination); !errors.Is(err, fs.ErrNotExist) {
		if err != nil {
			return err
		}
		return errors.New("backup destination already exists")
	}
	scratch, err := os.MkdirTemp(filepath.Dir(destination), ".cinefolio-backup-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(scratch)
	if err := store.BackupFile(filepath.Join(data, "cinefolio.db"), filepath.Join(scratch, "cinefolio.db")); err != nil {
		return err
	}
	assets, err := assetsAt(scratch)
	if err != nil {
		return err
	}
	sourceRoot := filepath.Join(data, "media")
	source, err := media.New(sourceRoot)
	if err != nil {
		return err
	}
	target := filepath.Join(scratch, "media")
	if err := os.Mkdir(target, 0700); err != nil {
		return err
	}
	for _, asset := range assets {
		if err := backupAsset(source, sourceRoot, target, asset); err != nil {
			return fmt.Errorf("back up asset %d (retry if it was just replaced): %w", asset.ID, err)
		}
	}
	if _, _, err := verifyData(scratch); err != nil {
		return err
	}
	return os.Rename(scratch, destination)
}

func backupAsset(source *media.Storage, sourceRoot, target string, asset media.Asset) error {
	original, err := source.Open(asset.StorageKey)
	if err != nil {
		return err
	}
	defer original.Close()
	output := filepath.Join(target, asset.StorageKey)
	// Assets are immutable. A hard link keeps the snapshot copy available after a replacement.
	if err := os.Link(filepath.Join(sourceRoot, asset.StorageKey), output); err == nil {
		return nil
	}
	snapshot, err := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	size, err := io.Copy(snapshot, original)
	if err == nil && size != asset.Size {
		err = errors.New("media size changed")
	}
	if err == nil {
		err = snapshot.Sync()
	}
	return errors.Join(err, snapshot.Close())
}
