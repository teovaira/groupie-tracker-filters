# Groupie Tracker Geolocalization & Filters

An extension of Groupie Tracker that geocodes each artist's concert locations onto an interactive map and lets users filter the artist list by creation date, first album date, number of members, and concert locations.

## About

This project builds on the base Groupie Tracker in two ways. First, it resolves each concert location string (e.g. "san_francisco-usa") into real-world coordinates using the Nominatim geocoding service, then renders them as markers on a per-artist map. Geocoded results are cached to disk so repeated runs don't re-query the geocoding service for locations already resolved.

Second, it lets users narrow the artist list with range filters (creation date, first album year, number of members) and a checkbox filter (concert locations, grouped by country), combined with the existing live search — all applied asynchronously, without a page reload.

Built with Go (standard library only) and plain HTML/CSS/JS, using Leaflet for map rendering.

## Quick Start

**Prerequisites:** Go 1.22+

```bash
git clone https://github.com/teovaira/groupie-tracker-filters.git
cd groupie-tracker-filters

# run directly (no build step)
go run ./cmd

# or build a binary first
make build
./bin/groupie-tracker-filters
```

Visit `http://localhost:8080`

First run is instant — `data/geocache.json` ships pre-populated. If you delete it, the app will re-geocode every location live against Nominatim (~1 request/second, so a few minutes depending on data size).

## Features

- Everything from the base Groupie Tracker project
- Concert locations geocoded into latitude/longitude coordinates
- Interactive map on each artist page with a marker per concert location
- Disk-backed geocoding cache — avoids re-querying already-resolved locations across restarts
- Rate-limited geocoding requests to respect the Nominatim usage policy
- Graceful handling of unresolvable locations — a failed geocode is logged and skipped without breaking the rest of the map
- Range filters for creation date, first album year, and number of members
- Checkbox filter for concert locations, grouped by country
- Filters combine with each other and with live search using logical AND
- Filtering is fully asynchronous — results update as filters change, with no page reload
- Filter matching runs across a goroutine worker pool for concurrent evaluation

## Project Structure

```
cmd/                  entry point
data/                 geocoding cache (geocache.json)
internal/
  api/                external API client
  geo/                geocoder interface, real/mock implementations, cache
  handlers/           HTTP handlers
  models/             data structs
  store/              data layer + Store interface
web/
  static/             CSS, JS
  templates/          HTML templates
```

## Development

```bash
make fmt        # format
make test       # run tests
make check      # full pre-PR check (fmt + lint + build + test)
```

## Team

- **Vasiliki** — Backend: API client, models, store, search handler, geo package
- **Krysta** — Frontend: templates, CSS
- **Theo** — Full-stack: handlers, search.js, docs, QA

## License

MIT
