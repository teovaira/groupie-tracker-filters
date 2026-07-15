package store

import (
	"groupie-tracker-filters/internal/models"
	"html/template"
)

// MockStore is a test implementation of the Store interface.
// It returns a small set of hardcoded artists and uses the same matchesQuery
// logic as RealStore so that search behaviour is consistent in tests.
// It has no external dependencies and requires no setup beyond instantiation.
type MockStore struct{}

// AllArtists returns a fixed list of two artists used across all handler tests.
func (m *MockStore) AllArtists() []models.Artist {
	return []models.Artist{
		{ID: 1, Name: "Billie Eilish"},
		{ID: 2, Name: "System of a down"},
	}
}

// ArtistByID scans the hardcoded artist list and returns the matching artist.
// Returns false if the ID does not match either fixture artist.
func (m *MockStore) ArtistByID(id int) (models.Artist, bool) {
	for _, a := range m.AllArtists() {
		if a.ID == id {
			return a, true
		}
	}
	return models.Artist{}, false
}

// SearchArtists filters the hardcoded artist list using the same matchesQuery
// function as RealStore, ensuring search tests reflect real filtering behaviour.
func (m *MockStore) SearchArtists(query string) []models.Artist {
	var results []models.Artist
	for _, a := range m.AllArtists() {
		if matchesQuery(a, query) {
			results = append(results, a)
		}
	}
	return results
}

// FilterArtists filters the hardcoded artist list by query and criteria using
// the same matchesQuery and matchesCriteria logic as RealStore. The fixture
// artists carry no location data, so any non-empty criteria.Locations never
// matches. Filtering runs sequentially rather than through a worker pool —
// unlike RealStore, the fixture set is tiny and tests benefit more from
// straightforward, deterministic execution than from concurrency.
func (m *MockStore) FilterArtists(query string, criteria FilterCriteria) []models.Artist {
	var results []models.Artist
	for _, a := range m.AllArtists() {
		if matchesQuery(a, query) && matchesCriteria(a, nil, criteria) {
			results = append(results, a)
		}
	}
	return results
}

// ArtistPageDataByID returns an ArtistPageData for the matching fixture artist
// with empty locations, dates, datesLocations, markers, and MarkersJSON,
// since the mock holds no concert or geocoding data.
func (m *MockStore) ArtistPageDataByID(id int) (models.ArtistPageData, bool) {
	for _, a := range m.AllArtists() {
		if a.ID == id {
			return models.ArtistPageData{
				Artist:         a,
				Locations:      []string{},
				Dates:          []string{},
				DatesLocations: map[string][]string{},
				Markers:        []models.Marker{},
				MarkersJSON:    template.JS("[]"),
			}, true
		}
	}
	return models.ArtistPageData{}, false
}
