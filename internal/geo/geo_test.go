package geo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestGeocoder(handler http.HandlerFunc) (*RealGeocoder, func()) {
	srv := httptest.NewServer(handler)
	g := NewRealGeocoder(srv.URL)
	return g, srv.Close
}

func TestRealGeocoder_Success(t *testing.T) {
	g, close := newTestGeocoder(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"lat":"51.5074","lon":"-0.1278"}]`)) //nolint:errcheck // test fixture response write, error unrecoverable
	})
	defer close()

	got, err := g.Geocode("london-uk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Latitude != 51.5074 || got.Longitude != -0.1278 {
		t.Errorf("got %+v, want {51.5074 -0.1278}", got)
	}
}

func TestRealGeocoder_EmptyResults(t *testing.T) {
	g, close := newTestGeocoder(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`)) //nolint:errcheck // test fixture response write, error unrecoverable
	})
	defer close()

	_, err := g.Geocode("unknown-place")
	if err == nil {
		t.Error("expected error for empty results, got nil")
	}
}

func TestRealGeocoder_BadJSON(t *testing.T) {
	g, close := newTestGeocoder(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json`)) //nolint:errcheck // test fixture response write, error unrecoverable
	})
	defer close()

	_, err := g.Geocode("london-uk")
	if err == nil {
		t.Error("expected error for bad JSON, got nil")
	}
}

func TestRealGeocoder_QueryParams(t *testing.T) {
	g, close := newTestGeocoder(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "paris, france" {
			t.Errorf("unexpected q param: %s", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("format") != "json" {
			t.Errorf("unexpected format param: %s", r.URL.Query().Get("format"))
		}
		w.Write([]byte(`[{"lat":"48.8566","lon":"2.3522"}]`)) //nolint:errcheck // test fixture response write, error unrecoverable
	})
	defer close()

	// Return values are unused — this test only asserts on the query
	// parameters the handler above receives, not the geocode result.
	_, _ = g.Geocode("paris-france")
}

func TestNormalizeLocation(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"paris-france", "paris, france"},
		{"san_francisco-usa", "san francisco, usa"},
		{"north_carolina-usa", "north carolina, usa"},
		{"london-uk", "london, uk"},
	}
	for _, tc := range tests {
		got := normalizeLocation(tc.input)
		if got != tc.want {
			t.Errorf("normalizeLocation(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestRealGeocoder_ImplementsInterface(t *testing.T) {
	var _ Geocoder = &RealGeocoder{}
}
