# Contributing

## Workflow

1. Branch off `develop`
2. Write your tests first (TDD)
3. Implement the feature or fix
4. Run `make check` — all checks must pass
5. Open a PR against `main`

## Commits

Use conventional commits — one logical change per commit:

```
feat:     new feature
fix:      bug fix
test:     tests only
docs:     documentation
refactor: restructuring without behaviour change
chore:    tooling, config, cleanup
```

Example: `feat(handlers): implement GET /artist/{id} handler`

## Tests

- Go: table-driven tests, one `_test.go` per package
- JS: unit tests colocated with the file under test (e.g. `search.test.js`, `map.test.js`)
- Run: `make test`

## Code Style

- Go: follow standard conventions (`gofmt`, exported names PascalCase, errors lowercase); only standard library packages — no external Go dependencies (audited requirement)
- Frontend: vanilla JS/CSS, except Leaflet (loaded via CDN, not a Go package) for map rendering
- Functions under 50 lines

## Package Ownership

Each package has a primary owner. Discuss in a PR before modifying a package you did not author.
