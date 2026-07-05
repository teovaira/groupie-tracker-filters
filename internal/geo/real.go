package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type RealGeocoder struct {
	BaseURL string
	Client  *http.Client
}
type placeResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func NewRealGeocoder(baseURL string) *RealGeocoder {
	return &RealGeocoder{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}
func (g *RealGeocoder) Geocode(address string) (Coordinates, error) {
	u, err := url.Parse(g.BaseURL)
	if err != nil {
		return Coordinates{}, err
	}
	q := u.Query()
	q.Set("q", address)
	q.Set("format", "json")
	u.RawQuery = q.Encode()
	resp, err := g.Client.Get(u.String())
	if err != nil {
		return Coordinates{}, err
	}
	defer resp.Body.Close()
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
