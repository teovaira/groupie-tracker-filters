// Package store defines the Store interface and provides two implementations:
// RealStore, which holds data loaded from the external API at startup, and
// MockStore, which returns hardcoded data for use in tests. All HTTP handlers
// depend on the Store interface rather than a concrete type, keeping them
// decoupled from the data source and easy to test in isolation.
package store

import "groupie-tracker/internal/models"

// Store is the data access interface used by all HTTP handlers.
// It abstracts over the underlying data source so that handlers can be tested
// with a MockStore without needing a live API or filesystem.
type Store interface {
	// AllArtists returns the full list of artists.
	AllArtists() []models.Artist
	// ArtistByID returns the artist with the given ID and true, or a zero-value
	// Artist and false if no artist with that ID exists.
	ArtistByID(id int) (models.Artist, bool)
	// SearchArtists returns all artists whose name, members, creation date, or
	// first album contain the given query string (case-insensitive).
	SearchArtists(query string) []models.Artist
	// ArtistPageDataByID returns the combined artist, locations, and dates data
	// needed to render the artist detail page, and true. Returns false if the ID
	// does not match any artist.
	ArtistPageDataByID(id int) (models.ArtistPageData, bool)
}
