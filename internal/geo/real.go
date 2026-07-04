package geo

import "net/http"

type RealGeocoder struct {
	BaseURL string
	Client  *http.Client
}

func NewRealGeocoder(baseURL string) *RealGeocoder {
	return &RealGeocoder{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}
func (g *RealGeocoder) Geocoder(adress string) (Coordinates, error) {

}
