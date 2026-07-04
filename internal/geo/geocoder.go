// Package geo defines the Geocoder interface and the Coordinates type used to
// represent geographic positions. Concrete implementations of Geocoder resolve
// a human-readable address to a latitude/longitude pair; handlers depend on
// the interface rather than a concrete type so that geocoding can be swapped
// or mocked in tests without network access.
package geo

// Coordinates holds the geographic position returned by a Geocoder.
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// Geocoder is the geocoding interface used by handlers that need to plot
// concert locations on a map. It abstracts over the underlying geocoding
// service so that a fake implementation can be injected during testing.
type Geocoder interface {
	// Geocode resolves address to a Coordinates pair. It returns an error if
	// the address cannot be resolved or the underlying service is unavailable.
	Geocode(address string) (Coordinates, error)
}
