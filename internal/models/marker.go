package models

// Marker holds the geocoded position of a single concert location.
// Name is the human-readable location string (e.g. "london-uk") and
// Lat/Lng are the coordinates resolved by the Geocoder.
type Marker struct {
	Name string
	Lat  float64
	Lng  float64
}
