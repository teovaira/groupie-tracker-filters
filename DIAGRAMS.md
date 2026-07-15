# Diagrams

## Application Flow

```mermaid
flowchart TD
    A([Browser]) -->|GET /| B[HomeHandler]
    A -->|"GET /artist/{id}"| C[ArtistHandler]
    A -->|"GET /api/search?q="| D[SearchHandler]
    A -->|"GET /api/filter?..."| Y[FilterHandler]
    A -->|"GET /static/..."| E[FileServer]

    B --> F[store.AllArtists + store.LocationGroups]
    C -->|valid numeric ID| G[store.ArtistPageDataByID]
    C -->|non-numeric or unknown ID| Q[404 Not Found]
    D -->|empty query| P[400 Bad Request]
    D --> H[store.SearchArtists]
    Y -->|malformed numeric param| P
    Y --> Z[store.FilterArtists]

    F --> I[(RealStore)]
    G -->|ArtistPageData incl. MarkersJSON| I
    Z -->|goroutine worker pool| I

    I -->|loaded once at startup| J[api.LoadData]
    J -->|fetch| K([groupietrackers API])

    I -->|markers loaded once at startup| R[Geocoding Pipeline]
    R --> S{cache.Get}
    S -->|hit| T[Coordinates]
    S -->|miss| U[geocoder.Geocode]
    U -->|normalize + rate-limit| V([Nominatim API])
    U --> W[cache.Set]
    W --> T
    T --> R

    B -->|render| L[home.html]
    C -->|render| M[artist.html]

    L -->|extends| N[base.html]
    M -->|extends| N

    M -->|embeds MarkersJSON, loads| X[map.js + Leaflet]
```

## Startup Lifecycle — Geocoding

```mermaid
sequenceDiagram
    participant M as main.go
    participant API as api.LoadData
    participant C as geo.Cache
    participant G as geo.RealGeocoder
    participant N as Nominatim

    M->>API: LoadData()
    API-->>M: artists, locations, dates, relations
    M->>C: NewCache(geoCachePath)
    loop for each artist location
        M->>C: Get(location)
        alt cache hit
            C-->>M: Coordinates
        else cache miss
            M->>G: Geocode(location)
            G->>G: normalizeLocation(location)
            G->>N: GET /search?q=...
            N-->>G: [{lat, lon}] or []
            G-->>M: Coordinates or error
            M->>M: sleep 1.1s (rate limit)
            M->>C: Set(location, coords)
        end
    end
    M->>C: Save()
    M->>M: build RealStore with Markers
```

## Request Lifecycle — Search & Filter

`search.js` no longer calls `/api/search` on its own — `filter.js` owns
`#search-results` and reads the search box alongside the filter panel inputs,
combining both into a single `/api/filter` request with logical AND.

```mermaid
sequenceDiagram
    participant U as User
    participant JS as filter.js
    participant S as Server
    participant ST as Store

    U->>JS: types in search input or changes a filter
    JS->>JS: debounce 300ms
    JS->>JS: buildFilterQuery(state)
    JS->>S: GET /api/filter?q=...&creation_min=...&locations=...
    S->>S: parse + validate numeric params
    alt malformed numeric param
        S-->>JS: 400 Bad Request
    else valid
        S->>ST: FilterArtists(query, criteria)
        ST->>ST: fan out to goroutine worker pool
        ST->>ST: matchesQuery && matchesCriteria per artist
        ST->>ST: re-sort matches into input order
        ST-->>S: []Artist
        S-->>JS: JSON array (never null)
        JS->>U: render cards dynamically
    end
```