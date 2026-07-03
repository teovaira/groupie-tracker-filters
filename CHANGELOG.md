# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[1.0.0]: https://github.com/vxanthio/groupie-tracker.git
