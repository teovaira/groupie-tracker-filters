# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-07-15

### Added
- `FilterCriteria` struct and `Store.FilterArtists(query, criteria)` for combined free-text and structured artist filtering
- Range filters for creation date, first album year, and number of members
- Checkbox filter for concert locations, grouped by country and matched with substring semantics (e.g. a `washington-usa` selection also matches a more specific `seattle-washington-usa` slug)
- `Store.LocationGroups()` and `models.LocationGroup`/`models.HomePageData` to populate the location checkbox filter from the full, unfiltered location vocabulary
- `GET /api/filter` endpoint, mirroring `/api/search`'s JSON contract, accepting `q`, `creation_min`/`creation_max`, `first_album_min`/`first_album_max`, `members_min`/`members_max`, and repeated `locations` query parameters — all optional, so no parameters returns every artist
- Goroutine worker pool in `RealStore.FilterArtists`, fanning artist-predicate evaluation across `runtime.NumCPU()` workers, with results re-sorted back into deterministic input order after collection
- Filter panel UI on the home page: range inputs plus collapsible, country-grouped location checkboxes
- `filter.js` — debounced, asynchronous wiring of the filter panel and the existing search box to `/api/filter`, so search and filters combine live without a page reload
- Filter panel styling matching the existing dark-theme design system
- Full test coverage for the filtering feature: `matchesCriteria`, `RealStore`/`MockStore` `FilterArtists` and `LocationGroups`, `FilterHandler`, route wiring, and `filter.js`'s `buildFilterQuery`

### Changed
- Go module renamed from `groupie-tracker-geolocalization` to `groupie-tracker-filters`, along with the Nominatim `User-Agent` header, build binary name, and repository links, to match the new feature and repository
- `HomeHandler` now passes `models.HomePageData{Artists, LocationGroups}` to the home template instead of a raw artist slice
- `search.js` no longer self-initializes; `filter.js` now owns `#search-results` and reads the search box as part of the combined `/api/filter` request

### Fixed
- `FilterHandler` normalizes a nil `FilterArtists` result to an empty JSON array (`[]`) rather than `null`, matching `/api/search`'s existing contract

## [1.1.0] - 2026-07-12

### Added
- `Geocoder` interface with `RealGeocoder` (Nominatim-backed) and `MockGeocoder` implementations
- File-backed geocoding cache to avoid redundant API calls across restarts
- Rate-limit delay integrated into the startup geocoding loop, respecting Nominatim's usage policy
- Location slug normalization (e.g. `san_francisco-usa` → `san francisco, usa`) before geocoding
- Artist location markers generated at startup and attached to `ArtistPageData`
- Concert map section on the artist detail page, rendered with Leaflet
- Marker data marshalled to JSON and injected into the template for client-side map rendering
- Pre-warmed `data/geocache.json` committed to the repository for instant first-run startup
- Test coverage for `RealStore` marker/`MarkersJSON` attachment and handler marker rendering
- Unit tests for `initMap` (`map.test.js`)

### Fixed
- Nil pointer risk, response body leak, and missing `User-Agent` header in `RealGeocoder`
- `MarkersJSON` now uses `template.JS` to prevent HTML escaping of embedded JSON
- Rate-limit sleep moved before the error check in the geocoding loop, so it applies on both success and failure paths
- Added concurrency safety, parent directory creation, and wrapped errors to the geocoding `Cache`
- Added a 10-second timeout to `RealGeocoder`'s HTTP client to prevent startup hanging indefinitely on an unresponsive request
- Nil marker slices now marshal to `[]` instead of `null`
- `search.js` no longer crashes under Node — `window`/`document` access guarded so `search.test.js` can `require()` it (pre-existing bug, unrelated to geolocalization)

### Changed
- `RealGeocoder` godoc comments updated to reflect corrected `Geocode` behaviour
- `MarkersJSON` marshalling moved from the handler into `RealStore.ArtistPageDataByID`, so every caller of the store gets it automatically
- Go module renamed to `groupie-tracker-geolocalization` to match the repository and documentation
- Comments expanded and corrected across `cmd`, `store`, `models`, and `handlers`
- Indentation fixed in the Concert Map section of `artist.html`
- Redundant comments removed from `map.test.js`

## [1.0.0] - 2026-04-13

### Added
- Project bootstrap: Go module, folder structure, HTTP server
- API client fetching artists, locations, dates, and relations at startup
- In-memory data store with `Store` interface, `RealStore`, and `MockStore` for testing
- `GET /` handler rendering the full artist list as cards
- `GET /artist/{id}` handler rendering artist detail page
- `GET /api/search?q=` handler returning JSON for live search
- Live search via debounced client-side `fetch()` without full page reload
- Artist detail page: members, creation year, first album, locations, dates, and dates grouped by location
- Locations and Concert Dates displayed as pill badges on artist detail page
- Dates by Location displayed as a card grid on artist detail page
- Base layout, home template, artist template
- Styled error pages for 400 Bad Request, 404 Not Found, and 500 Internal Server Error
- Dark-theme CSS design system with responsive layout
- Static file server under `/static/`
- `RecoveryMiddleware` wrapping the entire mux — catches panics, renders 500, never crashes

### Fixed
- Relations data wired into `ArtistPageData` and displayed on the artist detail page
- Non-numeric and unknown artist IDs consistently return 404 Not Found
- Empty or missing `?q=` search parameter returns 400 Bad Request
- Renamed `AppData.Date` → `Dates` for consistency
- HTTP error strings lowercased to follow Go conventions

[1.2.0]: https://github.com/teovaira/groupie-tracker-filters.git
[1.1.0]: https://github.com/vxanthio/groupie-tracker-geolocalization.git
[1.0.0]: https://github.com/vxanthio/groupie-tracker.git