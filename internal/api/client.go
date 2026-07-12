// Package api handles all communication with the Groupie Trackers external REST API.
// It fetches four datasets — artists, locations, dates, and relations — over HTTP,
// decodes the JSON responses into typed models, and holds the result in a
// package-level store that the rest of the application reads at request time.
// A shared http.Client with a 10-second timeout is used for all outbound requests
// to prevent the server from hanging if the upstream API is slow or unresponsive.
package api

import (
	"encoding/json"
	"fmt"
	"groupie-tracker-geolocalization/internal/models"
	"net/http"
	"time"
)

// AppData holds the complete dataset fetched from the external API at startup.
// It is populated once by LoadData and then read-only for the lifetime of the server.
type AppData struct {
	Artists   []models.Artist
	Locations models.LocationsResponse
	Dates     models.DatesResponse
	Relations models.RelationResponse
}

var data AppData

// LoadData is the public entry point for populating the in-memory data store.
// It calls loadDataFromURLs with the four production API endpoints and should
// be invoked exactly once during application startup, before the HTTP server
// begins accepting requests. If any endpoint fails to respond or returns
// malformed JSON, an error is returned and the application should not start.
func LoadData() error {
	return loadDataFromURLs(
		"https://groupietrackers.herokuapp.com/api/artists",
		"https://groupietrackers.herokuapp.com/api/locations",
		"https://groupietrackers.herokuapp.com/api/dates",
		"https://groupietrackers.herokuapp.com/api/relation",
	)
}

// GetData returns a copy of the AppData struct populated by LoadData.
// It is used by main to pass the loaded datasets into the RealStore,
// which the HTTP handlers then query for every incoming request.
func GetData() AppData {
	return data
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// fetchAndDecode fetches the given URL and JSON-decodes the response body into target.
// It closes the response body before returning.
func fetchAndDecode(url string, target any) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // deferred close, error unrecoverable
	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode failed: %w", err)
	}
	return nil
}

// loadDataFromURLs sequentially fetches and JSON-decodes each of the four API
// endpoints into the package-level data variable via fetchAndDecode.
// It is intentionally separate from LoadData so that tests can inject a local
// httptest.Server URL instead of hitting the real API.
// Errors from any step are wrapped with context and returned immediately —
// no partial data is used if any step fails.
func loadDataFromURLs(artistsURL, locationsURL, datesURL, relationsURL string) error {
	if err := fetchAndDecode(artistsURL, &data.Artists); err != nil {
		return fmt.Errorf("artists: %w", err)
	}
	if err := fetchAndDecode(locationsURL, &data.Locations); err != nil {
		return fmt.Errorf("locations: %w", err)
	}
	if err := fetchAndDecode(datesURL, &data.Dates); err != nil {
		return fmt.Errorf("dates: %w", err)
	}
	if err := fetchAndDecode(relationsURL, &data.Relations); err != nil {
		return fmt.Errorf("relations: %w", err)
	}
	return nil
}
