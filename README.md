# Groupie Tracker Geolocalization

An extension of Groupie Tracker that geocodes each artist's concert locations and plots them on an interactive map.

## About

Groupie Tracker Geolocalization builds on the base project by resolving each concert location string (e.g. "san_francisco-usa") into real-world coordinates using the Nominatim geocoding service, then rendering them as markers on a per-artist map. Geocoded results are cached to disk so repeated runs don't re-query the geocoding service for locations already resolved.

Built with Go (standard library only) and plain HTML/CSS/JS, using Leaflet for map rendering.

## Quick Start

**Prerequisites:** Go 1.22+

```bash
git clone https://github.com/vxanthio/groupie-tracker-geolocalization.git
cd groupie-tracker-geolocalization
go build -o groupie-tracker-geolocalization ./cmd
./groupie-tracker-geolocalization
```

Visit `http://localhost:8080`

## Features

- Everything from the base Groupie Tracker project
- Concert locations geocoded into latitude/longitude coordinates
- Interactive map on each artist page with a marker per concert location
- Disk-backed geocoding cache — avoids re-querying already-resolved locations across restarts
- Rate-limited geocoding requests to respect the Nominatim usage policy
- Graceful handling of unresolvable locations — a failed geocode is logged and skipped without breaking the rest of the map

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

- **Vasiliki** — Backend: API client, models, store, search handler,geo package
- **Krysta** — Frontend: templates, CSS
- **Theo** — Full-stack: handlers, search.js, docs, QA

## License

MIT
