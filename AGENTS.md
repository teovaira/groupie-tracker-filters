# AGENTS.md

This file provides guidance for AI coding agents working on this repository.

## Project

Groupie Tracker Geolocalization & Filters is a Go web application that extends Groupie Tracker with geocoded concert locations rendered as map markers, and with asynchronous range/checkbox filtering of the artist list. Standard library only — no external Go packages. Leaflet (JS, loaded via CDN) is used client-side for map rendering.

## Commands

```bash
make fmt       # format Go code before committing
make test      # run all tests
make build     # build the binary
make check     # full pre-commit check
go build ./... # verify compilation
```

## Architecture

- Entry point: `cmd/main.go` — loads data, geocodes concert locations, builds store, registers routes
- Handlers: `internal/handlers/` — one file per route, injectable store + template
- Store: `internal/store/` — `Store` interface, `RealStore`, `MockStore`, plus `filter.go` (`FilterCriteria`, `matchesCriteria`, `firstAlbumYear`, location-grouping helpers)
- Models: `internal/models/` — data structs matching the external API, plus view models (`ArtistPageData`, `HomePageData`, `LocationGroup`)
- API client: `internal/api/` — fetches and caches all data on startup
- Geo: `internal/geo/` — `Geocoder` interface, `RealGeocoder` (Nominatim), `MockGeocoder`, and a file-backed `Cache` (`data/geocache.json`)
- Templates: `web/templates/` — Go `html/template`, base + page pattern
- Static: `web/static/` — CSS and JS served under `/static/`, including `filter.js` (async filter panel + search wiring)

## Rules

- TDD always — write the failing test before the implementation
- Table-driven tests in Go
- One logical change per commit, conventional commit messages
- Never add external Go packages
- Never modify another team member's files (see CONTRIBUTING.md)
- Never remove the rate-limit delay in the geocoding loop (`cmd/main.go`) — Nominatim's usage policy caps requests at 1/second; removing or misplacing the delay risks the app getting rate-limited or blocked
- Never remove the HTTP client timeout on `RealGeocoder` — an unbounded request can hang startup indefinitely with no error logged
- Always set a `User-Agent` header on any request to the geocoding service
- Geocoding failures must never abort the whole loop — log and skip the failed location, continue with the rest
- `data/geocache.json` is intentionally committed (not gitignored) — it's a pre-warmed cache so the app starts instantly without live-geocoding on first run. Don't delete it or re-add it to `.gitignore` without discussing with the team.
- Never remove the result re-sorting step in `RealStore.FilterArtists` — matching runs across a goroutine worker pool, so results arrive on the channel in completion order, not input order; without re-sorting by original index, identical filter requests could return artists in a different order on each call.

## Testing

```bash
go test ./...                          # run all tests
go test -coverprofile=coverage.out ./... # with coverage
```

## Key Interfaces

- `store.Store` is the central interface handlers depend on. Use `store.MockStore` or the local `testStore` in handler tests — never call the external API in tests.
- `geo.Geocoder` is the interface used for resolving addresses to coordinates. Use `geo.MockGeocoder` in tests — never call the real Nominatim endpoint in tests.
- `store.FilterCriteria` bounds are `*int` so "unset" (nil) is distinguishable from a real zero value — don't switch these to plain `int` with a sentinel like `0` or `-1`, it silently breaks "no constraint" semantics.
- `store.MockStore` has no location data, so any test asserting on `criteria.Locations` matches must use `RealStore` with a populated `Locations` index instead.