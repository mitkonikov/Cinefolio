package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"cinefolio/internal/auth"
	"cinefolio/internal/catalog"
	"cinefolio/internal/httpapi"
	"cinefolio/internal/media"
	"cinefolio/internal/store"
	"golang.org/x/term"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(args []string) error {
	data := valueOr("CINEFOLIO_DATA", "../data")
	database := filepath.Join(data, "cinefolio.db")
	if len(args) > 0 && args[0] == "backup-db" {
		options := flag.NewFlagSet("backup-db", flag.ContinueOnError)
		output := options.String("output", "", "new filename for the database snapshot")
		if err := options.Parse(args[1:]); err != nil {
			return err
		}
		if *output == "" || options.NArg() != 0 {
			return errors.New("usage: cinefolio backup-db --output FILE")
		}
		return store.BackupFile(database, *output)
	}
	db, err := store.Open(database)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	if len(args) > 0 {
		if args[0] != "set-password" {
			return fmt.Errorf("unknown command %q", args[0])
		}
		options := flag.NewFlagSet("set-password", flag.ContinueOnError)
		email := options.String("email", "", "administrator email")
		if err := options.Parse(args[1:]); err != nil {
			return err
		}
		if *email == "" || options.NArg() != 0 {
			return errors.New("usage: cinefolio set-password --email EMAIL")
		}
		fd := int(os.Stdin.Fd())
		if !term.IsTerminal(fd) {
			return errors.New("password entry requires a terminal")
		}
		fmt.Fprint(os.Stderr, "Password: ")
		password, err := term.ReadPassword(fd)
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return err
		}
		return auth.New(db).SetPassword(context.Background(), *email, string(password))
	}
	addr := valueOr("CINEFOLIO_ADDR", "127.0.0.1:5002")
	origin := valueOr("CINEFOLIO_ORIGIN", "http://127.0.0.1:5002")
	parsed, err := url.Parse(origin)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil || parsed.Path != "" || parsed.ForceQuery || parsed.RawQuery != "" || parsed.Fragment != "" {
		return errors.New("CINEFOLIO_ORIGIN must be an HTTP origin with no path")
	}
	secure, err := strconv.ParseBool(valueOr("CINEFOLIO_SECURE_COOKIE", strconv.FormatBool(parsed.Scheme == "https")))
	if err != nil {
		return fmt.Errorf("CINEFOLIO_SECURE_COOKIE: %w", err)
	}
	storage, err := media.New(filepath.Join(data, "media"))
	if err != nil {
		return err
	}
	site := valueOr("CINEFOLIO_SITE", "cinefolio")
	if err := catalog.New(db, site).CheckSite(context.Background()); err != nil {
		return err
	}
	server := &http.Server{Addr: addr, Handler: httpapi.New(db, storage, site, httpapi.Config{
		WebDir: valueOr("CINEFOLIO_WEB", "../frontend/build"), Origin: origin, SecureCookie: secure,
	}), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	done := make(chan error, 1)
	go func() { slog.Info("cinefolio serving", "address", addr); done <- server.ListenAndServe() }()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
}

func valueOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
