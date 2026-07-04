package store

import (
	"groupie-tracker/internal/models"
	"strconv"
	"strings"
)

// RealStore is the production implementation of the Store interface.
// It holds all data loaded from the external API at startup and serves it
// directly from memory on every request, avoiding repeated network calls.
type RealStore struct {
	Artists   []models.Artist
	Locations models.LocationsResponse
	Dates     models.DatesResponse
	Relations models.RelationResponse
	Markers   map[int][]models.Marker
}

// AllArtists returns the full list of artists held in memory.
func (r *RealStore) AllArtists() []models.Artist {
	return r.Artists
}

// ArtistByID performs a linear scan over all artists and returns the one
// matching the given ID. Returns false if no match is found.
func (r *RealStore) ArtistByID(id int) (models.Artist, bool) {
	for _, a := range r.AllArtists() {
		if a.ID == id {
			return a, true
		}
	}
	return models.Artist{}, false
}

// SearchArtists filters the artist list using matchesQuery and returns all
// artists that contain the query string in any searchable field.
func (r *RealStore) SearchArtists(query string) []models.Artist {
	var result []models.Artist
	for _, a := range r.AllArtists() {
		if matchesQuery(a, query) {
			result = append(result, a)
		}
	}
	return result
}

// matchesQuery performs a case-insensitive substring search across the artist's
// name, individual member names, creation date, and first album date.
// It returns true as soon as any field matches, without checking the rest.
func matchesQuery(a models.Artist, query string) bool {
	q := strings.ToLower(query)
	if strings.Contains(strings.ToLower(a.Name), q) {
		return true
	}
	for _, member := range a.Members {
		if strings.Contains(strings.ToLower(member), q) {
			return true
		}
	}
	if strings.Contains(strings.ToLower(strconv.Itoa(a.CreationDate)), q) {
		return true
	}
	if strings.Contains(strings.ToLower(a.FirstAlbum), q) {
		return true
	}

	return false
}

// ArtistPageDataByID looks up an artist by ID and assembles the ArtistPageData
// view model by joining locations, dates, and relations data from their respective
// index slices. Returns false if no artist with the given ID exists.
func (r *RealStore) ArtistPageDataByID(id int) (models.ArtistPageData, bool) {
	for _, a := range r.AllArtists() {
		if a.ID == id {
			var locations []string
			for _, l := range r.Locations.Index {
				if l.ID == id {
					locations = l.Locations
				}
			}
			var dates []string
			for _, d := range r.Dates.Index {
				if d.ID == id {
					dates = d.Dates
				}
			}
			var datesLocations map[string][]string
			for _, rel := range r.Relations.Index {
				if rel.ID == id {
					datesLocations = rel.DatesLocations
				}
			}
			markers := r.Markers[id]

			return models.ArtistPageData{
				Artist:         a,
				Locations:      locations,
				Dates:          dates,
				DatesLocations: datesLocations,
				Markers:        markers,
			}, true
		}
	}
	return models.ArtistPageData{}, false
}
