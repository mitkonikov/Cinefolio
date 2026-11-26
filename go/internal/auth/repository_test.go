package auth

import (
	"errors"
	"path/filepath"
	"testing"

	"cinefolio/internal/store"
)

func TestPasswordResetRevokesSessions(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	r := New(db)
	ctx := t.Context()
	email := "owner@example.com"
	if err := r.SetPassword(ctx, email, "first password"); err != nil {
		t.Fatal(err)
	}
	token, _, err := r.Login(ctx, email, "first password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Current(ctx, token); err != nil {
		t.Fatal(err)
	}
	if err := r.SetPassword(ctx, email, "second password"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Current(ctx, token); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("old session: %v", err)
	}
	if _, _, err := r.Login(ctx, email, "second password"); err != nil {
		t.Fatal(err)
	}
}
