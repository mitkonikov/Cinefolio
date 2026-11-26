package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"os"

	"cinefolio/internal/legacyimport"
	"cinefolio/internal/media"
	"gorm.io/gorm"
)

func importLegacy(db *gorm.DB, storage *media.Storage, args []string) error {
	flags := flag.NewFlagSet("import-legacy", flag.ContinueOnError)
	var options legacyimport.Options
	flags.StringVar(&options.Export, "export", "", "JSON produced by tools/export-legacy.js")
	flags.StringVar(&options.MediaRoot, "source-media", "", "legacy public directory containing files, thumbs and trailers")
	flags.BoolVar(&options.DryRun, "dry-run", false, "validate and roll back rows and copied files")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || options.Export == "" || options.MediaRoot == "" {
		return errors.New("usage: cinefolio import-legacy --export FILE --source-media DIRECTORY [--dry-run]")
	}
	report, err := legacyimport.Run(context.Background(), db, storage, options)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
