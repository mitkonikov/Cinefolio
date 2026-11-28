# AGENTS.md

## Purpose
Operational guidance for coding agents working in this repository (`Cinefolio`) to keep migrations and dependency management consistent and reproducible.

## Core Rules
- Do **not** manually edit dependency versions in `package.json`.
- Use CLI/package-manager commands for dependency changes (`npm install`, `npm update`, `npx npm-check-updates -u`, etc.).
- After dependency updates, always run install and validation commands in the relevant project folder.
- Prefer minimal, targeted code changes; avoid unrelated refactors.

## Dependency Management Policy
### 1) Update flow (CLI-first)
2. `npm outdated`
3. `npx npm-check-updates -u`
4. `npm install . --force` (required in this repo due to legacy/peer mismatches)
5. `npm outdated` (verify no remaining outdated packages)

### 2) Force install guidance
- In `frontend`, use `npm install . --force` when peer conflicts block installation.
- This is expected because of old peer ranges in some dependencies.

## Known Compatibility Notes
### `frontend`
- Peer warning expected: `svelte-material-icons@3.0.5` vs `svelte@5.x`.
- These warnings do not necessarily block install when using `--force`, but should be revisited.

## Validation Checklist
After dependency or migration work:
- Frontend (`frontend`):
  - `npm run check`
  - `npm run build`
  - optionally `npm run dev`
- Backend (`go`):
  - `go test ./...`
  - `go build ./cmd/cinefolio`

## Release compatibility

- Use stable `vMAJOR.MINOR.PATCH` tags. Rollback is allowed only between releases with the same major and minor versions, such as `v1.4.3` and `v1.4.2`.
- Every patch release, including `v0.x` patches, must leave SQLite migrations, the database schema, and the public API and ABI unchanged. The previous patch must be able to run against data and persistent state written by the new patch.
- Use at least a minor version bump for SQLite migrations or compatibility changes. Use a major version bump for breaking public API or ABI changes in stable major versions.
- Plan recovery separately for minor and major updates. Switching the application version does not restore its data.
