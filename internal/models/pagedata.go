package models

import "html/template"

// ArtistPageData is the view model passed to the artist detail template.
// It combines the core Artist data with the resolved concert data from all four
// API endpoints. DatesLocations is sourced from the relations endpoint and maps
// each concert location to the dates the artist performed there.
// Markers holds the pre-geocoded coordinates for each concert location,
// used to render the map on the artist detail page.
type ArtistPageData struct {
	Artist         Artist
	Locations      []string
	Dates          []string
	DatesLocations map[string][]string
	Markers        []Marker
	MarkersJSON    template.JS
}
