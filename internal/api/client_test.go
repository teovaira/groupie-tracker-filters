package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeTestServer(artistsJSON, locationsJSON, datesJSON, relationsJSON string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/artists", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(artistsJSON))
	})
	mux.HandleFunc("/api/locations", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(locationsJSON))
	})
	mux.HandleFunc("/api/dates", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(datesJSON))
	})
	mux.HandleFunc("/api/relation", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(relationsJSON))
	})
	return httptest.NewServer(mux)
}

func TestLoadData_Success(t *testing.T) {
	data = AppData{}
	srv := makeTestServer(
		`[{"id":1,"name":"Test Artist","members":["Alice"],"creationDate":2000,"firstAlbum":"01-01-2000"}]`,
		`{"index":[{"id":1,"locations":["paris"]}]}`,
		`{"index":[{"id":1,"dates":["2024-01-01"]}]}`,
		`{"index":[{"id":1,"datesLocations":{"paris":["2024-01-01"]}}]}`,
	)
	defer srv.Close()

	err := loadDataFromURLs(srv.URL+"/api/artists", srv.URL+"/api/locations", srv.URL+"/api/dates", srv.URL+"/api/relation")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data.Artists) != 1 || data.Artists[0].Name != "Test Artist" {
		t.Errorf("unexpected artists: %+v", data.Artists)
	}
	if len(data.Locations.Index) != 1 {
		t.Errorf("unexpected locations: %+v", data.Locations)
	}
	if len(data.Dates.Index) != 1 {
		t.Errorf("unexpected dates: %+v", data.Dates)
	}
	if len(data.Relations.Index) != 1 {
		t.Errorf("unexpected relations: %+v", data.Relations)
	}
}

func TestLoadData_ArtistFetchFail(t *testing.T) {
	data = AppData{}
	err := loadDataFromURLs("http://invalid.invalid", "", "", "")
	if err == nil {
		t.Error("expected error for bad artist URL")
	}
}

func TestLoadData_InvalidJSON(t *testing.T) {
	data = AppData{}
	srv := makeTestServer("not-json", `{"index":[]}`, `{"index":[]}`, `{"index":[]}`)
	defer srv.Close()

	err := loadDataFromURLs(srv.URL+"/api/artists", srv.URL+"/api/locations", srv.URL+"/api/dates", srv.URL+"/api/relation")
	if err == nil {
		t.Error("expected error for invalid artist JSON")
	}
}
