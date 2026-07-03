// Package models defines the data structures used throughout the application.
// Each type maps directly to the JSON shape returned by the Groupie Trackers API,
// with struct tags used for decoding. These types are shared across the api,
// store, and handlers packages and should not contain any business logic.
package models

// Artist represents a single music artist or band as returned by the /api/artists endpoint.
// The Locations, ConcertDates, and Relations fields hold URLs to the related sub-resources
// rather than the data itself — those are fetched separately and joined in the store layer.
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// LocationsResponse is the top-level wrapper returned by the /api/locations endpoint.
type LocationsResponse struct {
	Index []Locations `json:"index"`
}

// Locations holds the list of concert location strings for a single artist,
// identified by the artist ID. Location strings are formatted as "city-country".
type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

// DatesResponse is the top-level wrapper returned by the /api/dates endpoint.
type DatesResponse struct {
	Index []Dates `json:"index"`
}

// Dates holds the list of concert date strings for a single artist.
// Dates are formatted as DD-MM-YYYY. Past dates are prefixed with "*".
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// RelationResponse is the top-level wrapper returned by the /api/relation endpoint.
type RelationResponse struct {
	Index []Relation `json:"index"`
}

// Relation maps each concert location to the list of dates the artist performed there.
// The DatesLocations key is a location string (e.g. "london-uk") and the value
// is a slice of date strings (e.g. ["06-03-2020", "07-03-2020"]).
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
