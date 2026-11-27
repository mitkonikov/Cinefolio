# Cinefolio

A film portfolio with a static SvelteKit frontend, a Go API and editor, and a SQLite catalog. Go serves the public site, /admin/, and media on the same origin.

## Local development

Use Go 1.27.4 and Node.js 24. Run these commands from the repository root.

```sh
npm --prefix frontend ci
npm --prefix frontend run build
go -C go run ./cmd/cinefolio set-password --email owner@example.com
go -C go run ./cmd/cinefolio
```

Open `http://127.0.0.1:5002` or `http://127.0.0.1:5002/admin/`. The password command prompts in the terminal and also resets existing accounts and revokes their sessions. Data defaults to the repository's ignored `data/` directory. To use live Svelte reload, start Go with `CINEFOLIO_ORIGIN=http://127.0.0.1:5173`, then run `npm --prefix frontend run dev` and open that origin.

```sh
go -C go test ./...
npm --prefix frontend run check
npm --prefix frontend test -- --run
npm --prefix frontend run build
```
