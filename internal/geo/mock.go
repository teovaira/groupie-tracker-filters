package geo

import "errors"

// MockGeocoder is a test implementation of the Geocoder interface.
// It returns hardcoded coordinates for a small set of known addresses and an
// error for any unknown address, allowing tests to run without network access.
type MockGeocoder struct{}

// Geocode returns a fixed Coordinates pair for known fixture addresses.
// It returns an error for any address not present in the fixture map,
// simulating a failed geocoding lookup.
func (m *MockGeocoder) Geocode(address string) (Coordinates, error) {
	fixtures := map[string]Coordinates{
		"london-uk":    {Latitude: 51.5074, Longitude: -0.1278},
		"paris-france": {Latitude: 48.8566, Longitude: 2.3522},
	}
	if c, ok := fixtures[address]; ok {
		return c, nil
	}
	return Coordinates{}, errors.New("address not found: " + address)
}
