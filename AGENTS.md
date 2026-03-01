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
- In `cinefolio2`, use `npm install . --force` when peer conflicts block installation.
- This is expected because of old peer ranges in some dependencies.

## Known Compatibility Notes
### `cinefolio2`
- Peer warning expected: `rxfire@6.1.0` vs `firebase@12.x`.
- Peer warning expected: `svelte-material-icons@3.0.5` vs `svelte@5.x`.
- These warnings do not necessarily block install when using `--force`, but should be revisited.

## Validation Checklist
After dependency or migration work:
- Frontend (`cinefolio2`):
  - `npm run check`
  - `npm run build`
  - optionally `npm run dev`
- Backend:
  - `npm run build`
