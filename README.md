# Cinefolio

![Go](https://img.shields.io/badge/Go-backend-007D9C?style=flat-square&logo=go&logoColor=white)
![SvelteKit](https://img.shields.io/badge/SvelteKit-frontend-D53C13?style=flat-square&logo=svelte&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-catalog-356C9B?style=flat-square&logo=sqlite&logoColor=white)

Cinefolio is a cinematography portfolio for presenting films in curated playlists. A featured trailer can introduce the portfolio and lead viewers to the full film.

The built-in editor lets you add films, link to external videos or upload your own, and choose which playlists appear on the site. Upload thumbnails and trailers from the editor, or capture a thumbnail from a local video frame.

## Run locally

Install Go 1.27.4 and Node.js 24, the versions used by this project’s release build. Run these commands from the repository root in an interactive terminal.

```sh
npm --prefix frontend ci
npm --prefix frontend run build
go -C go run ./cmd/cinefolio set-password --email owner@example.com
go -C go run ./cmd/cinefolio
```

Use your own email for the owner account and enter a password when prompted. Open [the portfolio](http://127.0.0.1:5002/) or [the editor](http://127.0.0.1:5002/admin/). To put your first film on the homepage, add it to a playlist marked **Visible on portfolio**.

Go serves the built SvelteKit frontend, API, and uploaded media. The SQLite database and uploads are stored in the repository’s ignored `data/` directory by default.

## Develop and verify

For frontend live reload, start Go with `CINEFOLIO_ORIGIN` set to `http://127.0.0.1:5173`. In a second terminal, run `npm --prefix frontend run dev` and open [the Vite development site](http://127.0.0.1:5173/). Vite forwards API and media requests to Go.

```sh
go -C go test ./...
npm --prefix frontend run check
npm --prefix frontend test -- --run
npm --prefix frontend run build
```
