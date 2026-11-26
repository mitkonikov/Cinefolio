package legacyimport

import (
	"os"
	"path/filepath"
	"testing"

	"cinefolio/internal/catalog"
	"cinefolio/internal/media"
	"cinefolio/internal/store"
)

func TestDryRunAndRepeat(t *testing.T) {
	root := t.TempDir()
	db, err := store.Open(filepath.Join(root, "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	importMedia := filepath.Join(root, "legacy")
	if err := os.MkdirAll(filepath.Join(importMedia, "thumbs"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(importMedia, "thumbs", "sample.jpg"), []byte{0xff, 0xd8, 0xff, 0xd9}, 0600); err != nil {
		t.Fatal(err)
	}
	targetDir := filepath.Join(root, "media")
	target, err := media.New(targetDir)
	if err != nil {
		t.Fatal(err)
	}
	export := filepath.Join(root, "sample.json")
	input := `{
		"sites":[{"_id":"000000000000000000000001","name":"Cinefolio"}],
		"thumbnails":[{"_id":"000000000000000000000002","title":"Sample","filepath":"sample.jpg","crop":true}],
		"files":[{"_id":"000000000000000000000003","title":"Short","published":"2024-06-01T22:00:00.000Z","thumbnail":"000000000000000000000002","link":"https://example.com/short"}],
		"playlists":[{"_id":"000000000000000000000004","site":"000000000000000000000001","name":"Shorts","order":0}],
		"filetoplaylists":[{"_id":"000000000000000000000005","fileId":"000000000000000000000003","playlistId":"000000000000000000000004"}],
		"trailers":[]
	}`
	if err := os.WriteFile(export, []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	ctx := t.Context()
	options := Options{Export: export, MediaRoot: importMedia, DryRun: true}
	for _, dry := range []bool{true, false, false} {
		options.DryRun = dry
		report, err := Run(ctx, db, target, options)
		if err != nil {
			t.Fatal(err)
		}
		items, err := os.ReadDir(targetDir)
		if err != nil {
			t.Fatal(err)
		}
		portfolio, err := catalog.New(db, "cinefolio").Public(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if dry {
			if report.Created["films"] != 1 || len(items) != 0 || len(portfolio.Playlists) != 0 {
				t.Fatalf("dirty dry run: %+v", report)
			}
			continue
		}
		if len(items) != 1 || len(portfolio.Playlists) != 1 || len(portfolio.Playlists[0].Films) != 1 {
			t.Fatalf("unexpected import: %+v", report)
		}
		if publication := portfolio.Playlists[0].Films[0].Published; publication != "2024-06-02" {
			t.Errorf("publication: %s", publication)
		}
		if report.Created["films"]+report.Existing["films"] != 1 {
			t.Fatalf("unexpected repeat report: %+v", report)
		}
	}
}
