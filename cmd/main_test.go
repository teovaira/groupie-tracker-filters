package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"groupie-tracker-filters/internal/handlers"
	"groupie-tracker-filters/internal/models"
	"groupie-tracker-filters/internal/store"
)

type routeStore struct{}

func (r *routeStore) AllArtists() []models.Artist {
	return []models.Artist{
		{ID: 1, Name: "Test Artist", Image: "http://img/1.jpg", CreationDate: 2000},
	}
}

func (r *routeStore) ArtistByID(id int) (models.Artist, bool) {
	if id == 1 {
		return models.Artist{ID: 1, Name: "Test Artist"}, true
	}
	return models.Artist{}, false
}

func (r *routeStore) SearchArtists(query string) []models.Artist {
	return r.AllArtists()
}

func (r *routeStore) FilterArtists(query string, criteria store.FilterCriteria) []models.Artist {
	return r.AllArtists()
}

func (r *routeStore) ArtistPageDataByID(id int) (models.ArtistPageData, bool) {
	if id == 1 {
		return models.ArtistPageData{
			Artist:         models.Artist{ID: 1, Name: "Test Artist"},
			Locations:      []string{},
			Dates:          []string{},
			DatesLocations: map[string][]string{},
		}, true
	}
	return models.ArtistPageData{}, false
}

// Compile-time check: routeStore satisfies store.Store.
var _ store.Store = (*routeStore)(nil)

func buildMux(s store.Store) *http.ServeMux {
	homeTmpl := template.Must(template.New("base").Parse(`{{range .}}{{.Name}}{{end}}`))
	artistTmpl := template.Must(template.New("base").Parse(`{{.Artist.Name}}`))
	errorTmpl := template.Must(template.New("404.html").Parse(`Not Found`))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handlers.NotFoundHandler(errorTmpl)(w, r)
			return
		}
		handlers.NewHomeHandler(s, homeTmpl).ServeHTTP(w, r)
	})
	notFoundTmpl := template.Must(template.New("404.html").Parse(`Not Found`))
	badRequestTmpl := template.Must(template.New("400.html").Parse(`Bad Request`))
	mux.Handle("GET /artist/{id}", handlers.NewArtistHandler(s, artistTmpl, notFoundTmpl))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../web/static"))))
	searchHandler := &handlers.SearchHandler{Store: s, BadRequestTmpl: badRequestTmpl}
	mux.HandleFunc("GET /api/search", searchHandler.Search)
	filterHandler := &handlers.FilterHandler{Store: s, BadRequestTmpl: badRequestTmpl}
	mux.HandleFunc("GET /api/filter", filterHandler.Filter)
	return mux
}

func TestRoutes(t *testing.T) {
	s := &routeStore{}
	mux := buildMux(s)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	tests := []struct {
		name     string
		method   string
		path     string
		wantCode int
	}{
		{"home_returns_200", http.MethodGet, "/", http.StatusOK},
		{"artist_valid_id_returns_200", http.MethodGet, "/artist/1", http.StatusOK},
		{"artist_unknown_id_returns_404", http.MethodGet, "/artist/99", http.StatusNotFound},
		{"artist_non_numeric_id_returns_404", http.MethodGet, "/artist/abc", http.StatusNotFound},
		{"static_css_returns_200", http.MethodGet, "/static/css/styles.css", http.StatusOK},
		{"unknown_route_returns_404", http.MethodGet, "/nonexistent", http.StatusNotFound},
		{"search_with_query_returns_200", http.MethodGet, "/api/search?q=test", http.StatusOK},
		{"search_missing_q_returns_400", http.MethodGet, "/api/search", http.StatusBadRequest},
		{"filter_no_params_returns_200", http.MethodGet, "/api/filter", http.StatusOK},
		{"filter_malformed_param_returns_400", http.MethodGet, "/api/filter?members_min=abc", http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, srv.URL+tc.path, nil)
			if err != nil {
				t.Fatalf("could not build request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close() //nolint:errcheck // deferred close, error unrecoverable

			if resp.StatusCode != tc.wantCode {
				t.Errorf("GET %s = %d, want %d", tc.path, resp.StatusCode, tc.wantCode)
			}
		})
	}
}
