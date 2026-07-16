package models

import "html/template"

// ArtistPageData is the view model passed to the artist detail template.
// It combines the core Artist data with the resolved concert data from the
// API endpoints. DatesLocations is sourced from the relations endpoint and maps
// each concert location to the dates the artist performed there.
// Markers holds the pre-geocoded coordinates for each concert location and
// MarkersJSON is its JSON-encoded form, safe for direct embedding in templates,
// used to render the map on the artist detail page.
type ArtistPageData struct {
	Artist         Artist
	Locations      []string
	Dates          []string
	DatesLocations map[string][]string
	Markers        []Marker
	MarkersJSON    template.JS
}

// LocationGroup collects every concert location slug that shares a common
// country, for rendering as a labelled group of checkboxes in the location
// filter. Country is a human-readable label derived from the slug's country
// segment (e.g. "New Zealand" from "new_zealand"); Locations holds the raw,
// unmodified slugs used as filter values (e.g. "auckland-new_zealand").
type LocationGroup struct {
	Country   string
	Locations []string
}

// FilterBounds holds the minimum and maximum values for each range filter,
// derived from the full artist dataset. These seed the endpoints of the
// range sliders on the home page so that every slider position corresponds
// to real data rather than empty space beyond the actual range.
type FilterBounds struct {
	CreationMin   int
	CreationMax   int
	FirstAlbumMin int
	FirstAlbumMax int
	MembersMin    int
	MembersMax    int
}

// HomePageData is the view model passed to the home page template.
// Artists is the full or filtered artist list rendered as cards; LocationGroups
// is the complete, country-grouped vocabulary of concert locations available
// across all artists, used to populate the location checkbox filter; Bounds
// holds the min/max endpoints for the range sliders. None of these change
// based on which artists are currently displayed.
type HomePageData struct {
	Artists        []Artist
	LocationGroups []LocationGroup
	Bounds         FilterBounds
}
