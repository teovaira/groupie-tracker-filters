# Groupie Tracker

A web app that fetches data from a music API and displays artists, their members, concert locations, and dates.

## About

Groupie Tracker consumes a REST API with four endpoints — artists, locations, dates, and relations — and presents the data through a clean, browsable website. It includes a live search feature that queries the backend without a full page reload.

Built with Go (standard library only) and plain HTML/CSS/JS.

## Quick Start

**Prerequisites:** Go 1.22+

```bash
git clone https://github.com/vxanthio/groupie-tracker.git
cd groupie-tracker
go build -o groupie-tracker ./cmd
./groupie-tracker
```

Visit `http://localhost:8080`

## Features

- Browse all artists as cards on the home page
- View artist detail: members, creation year, first album, concert locations, dates, and dates grouped by location
- Live search — filter artists by name, member, location, or year without page reload
- Styled error pages for 400, 404, and 500 responses

## Project Structure

```
cmd/                  entry point
internal/
  api/                external API client
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

- **Vasiliki** — Backend: API client, models, store, search handler
- **Krysta** — Frontend: templates, CSS
- **Theo** — Full-stack: handlers, search.js, docs, QA

## License

MIT
