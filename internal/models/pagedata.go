package models

// ArtistPageData is the view model passed to the artist detail template.
// It combines the core Artist data with the resolved concert data from all four
// API endpoints. DatesLocations is sourced from the relations endpoint and maps
// each concert location to the dates the artist performed there.
type ArtistPageData struct {
	Artist         Artist
	Locations      []string
	Dates          []string
	DatesLocations map[string][]string
}
