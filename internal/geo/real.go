package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// RealGeocoder is the production implementation of the Geocoder interface.
// It resolves addresses by querying a Nominatim-compatible HTTP endpoint.
// BaseURL and Client are exported so that tests can inject a local server
// without making real network requests.
type RealGeocoder struct {
	BaseURL string
	Client  *http.Client
}

// placeResult maps the lat/lon fields returned by the Nominatim JSON response.
type placeResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// NewRealGeocoder returns a RealGeocoder pointed at baseURL with a default
// HTTP client. Pass the Nominatim search endpoint as baseURL
// (e.g. "https://nominatim.openstreetmap.org/search").
func NewRealGeocoder(baseURL string) *RealGeocoder {
	return &RealGeocoder{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// Geocode builds a request to BaseURL with the address as the q parameter,
// sets a User-Agent header, and parses the first result into a Coordinates pair.
// It returns an error if the HTTP request fails, the response status is not 200,
// no results are returned, or the response body cannot be decoded.
func (g *RealGeocoder) Geocode(address string) (Coordinates, error) {
	u, err := url.Parse(g.BaseURL)
	if err != nil {
		return Coordinates{}, err
	}
	q := u.Query()
	q.Set("q", address)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return Coordinates{}, err
	}
	req.Header.Set("User-Agent", "groupie-tracker-geolocalization")
	resp, err := g.Client.Do(req)
	if err != nil {
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("geocoding request for %q failed: status %d %s", address, resp.StatusCode, resp.Status)
	}
	var places []placeResult
	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		return Coordinates{}, err
	}
	if len(places) == 0 {
		return Coordinates{}, fmt.Errorf("location %q was not found", address)
	}
	latitude, err := strconv.ParseFloat(places[0].Lat, 64)
	if err != nil {
		return Coordinates{}, err
	}
	longitude, err := strconv.ParseFloat(places[0].Lon, 64)
	if err != nil {
		return Coordinates{}, err
	}

	return Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil

}
