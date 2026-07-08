# AGENTS.md

This file provides guidance for AI coding agents working on this repository.

## Project

Groupie Tracker Geolocalization is a Go web application that extends Groupie Tracker with geocoded concert locations, rendered as map markers. Standard library only — no external Go packages. Leaflet (JS, loaded via CDN) is used client-side for map rendering.

## Commands

```bash
make fmt       # format Go code before committing
make test      # run all tests
make build     # build the binary
make check     # full pre-commit check
go build ./... # verify compilation
```

## Architecture

- Entry point: `cmd/main.go` — loads data, builds store, registers routes
- Handlers: `internal/handlers/` — one file per route, injectable store + template
- Store: `internal/store/` — `Store` interface, `RealStore`, `MockStore`
- Models: `internal/models/` — data structs matching the external API
- API client: `internal/api/` — fetches and caches all data on startup
- Geo: `internal/geo/` — `Geocoder` interface, `RealGeocoder` (Nominatim), `MockGeocoder`, and a file-backed `Cache` (`data/geocache.json`)
- Templates: `web/templates/` — Go `html/template`, base + page pattern
- Static: `web/static/` — CSS and JS served under `/static/`

## Rules

- TDD always — write the failing test before the implementation
- Table-driven tests in Go
- One logical change per commit, conventional commit messages
- Never add external Go packages
- Never modify another team member's files (see CONTRIBUTING.md)
- Never remove the rate-limit delay in the geocoding loop (`cmd/main.go`) — Nominatim's usage policy caps requests at 1/second; removing or misplacing the delay risks the app getting rate-limited or blocked
- Always set a `User-Agent` header on any request to the geocoding service
- Geocoding failures must never abort the whole loop — log and skip the failed location, continue with the rest

## Testing

```bash
go test ./...                          # run all tests
go test -coverprofile=coverage.out ./... # with coverage
```

## Key Interfaces

- `store.Store` is the central interface handlers depend on. Use `store.MockStore` or the local `testStore` in handler tests — never call the external API in tests.
- `geo.Geocoder` is the interface used for resolving addresses to coordinates. Use `geo.MockGeocoder` in tests — never call the real Nominatim endpoint in tests.