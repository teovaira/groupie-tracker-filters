package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"groupie-tracker-filters/internal/models"
	"groupie-tracker-filters/internal/store"
)

type testStore struct {
	artists        []models.Artist
	locationGroups []models.LocationGroup
	bounds         models.FilterBounds
}

func (s *testStore) AllArtists() []models.Artist {
	return s.artists
}

func (s *testStore) ArtistByID(id int) (models.Artist, bool) {
	for _, a := range s.artists {
		if a.ID == id {
			return a, true
		}
	}
	return models.Artist{}, false
}

func (s *testStore) SearchArtists(query string) []models.Artist {
	return s.artists
}

func (s *testStore) FilterArtists(query string, criteria store.FilterCriteria) []models.Artist {
	return s.artists
}

func (s *testStore) LocationGroups() []models.LocationGroup {
	return s.locationGroups
}

func (s *testStore) FilterBounds() models.FilterBounds {
	return s.bounds
}

func (s *testStore) ArtistPageDataByID(id int) (models.ArtistPageData, bool) {
	for _, a := range s.artists {
		if a.ID == id {
			return models.ArtistPageData{
				Artist:         a,
				Locations:      []string{},
				Dates:          []string{},
				DatesLocations: map[string][]string{},
				Markers:        []models.Marker{{Name: "london-uk", Lat: 51.5074, Lng: -0.1278}},
				MarkersJSON:    template.JS(`[{"Name":"london-uk","Lat":51.5074,"Lng":-0.1278}]`),
			}, true
		}
	}
	return models.ArtistPageData{}, false
}

// Compile-time check: testStore satisfies store.Store.
var _ store.Store = (*testStore)(nil)

func mustParseTemplate(src string) *template.Template {
	return template.Must(template.New("base").Parse(src))
}

func brokenTemplate() *template.Template {
	// Calling a nil value forces an execution error without panicking the handler.
	tmpl, _ := template.New("base").Parse(`{{call .}}`)
	return tmpl
}

func TestHomeHandler(t *testing.T) {
	twoArtists := []models.Artist{
		{ID: 1, Name: "Foo Fighters", Image: "http://img/1.jpg", CreationDate: 1994},
		{ID: 2, Name: "Queen", Image: "http://img/2.jpg", CreationDate: 1970},
	}

	homeTmpl := mustParseTemplate(`{{range .Artists}}{{.Name}}{{end}}`)

	tests := []struct {
		name             string
		path             string
		artists          []models.Artist
		tmpl             *template.Template
		wantStatusCode   int
		wantBodyContains []string
	}{
		{
			name:             "happy_path_two_artists",
			path:             "/",
			artists:          twoArtists,
			tmpl:             homeTmpl,
			wantStatusCode:   http.StatusOK,
			wantBodyContains: []string{"Foo Fighters", "Queen"},
		},
		{
			name:             "empty_store_returns_200",
			path:             "/",
			artists:          []models.Artist{},
			tmpl:             homeTmpl,
			wantStatusCode:   http.StatusOK,
			wantBodyContains: []string{},
		},
		{
			name:             "unknown_path_returns_404",
			path:             "/nonexistent",
			artists:          twoArtists,
			tmpl:             homeTmpl,
			wantStatusCode:   http.StatusNotFound,
			wantBodyContains: []string{},
		},
		{
			name:             "template_error_returns_500",
			path:             "/",
			artists:          twoArtists,
			tmpl:             brokenTemplate(),
			wantStatusCode:   http.StatusInternalServerError,
			wantBodyContains: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &testStore{artists: tc.artists}
			h := NewHomeHandler(s, tc.tmpl)

			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatusCode)
			}

			body := rec.Body.String()
			for _, want := range tc.wantBodyContains {
				if !strings.Contains(body, want) {
					t.Errorf("body does not contain %q\nbody: %s", want, body)
				}
			}
		})
	}
}

func TestHomeHandler_PassesLocationGroupsToTemplate(t *testing.T) {
	s := &testStore{
		locationGroups: []models.LocationGroup{
			{Country: "Usa", Locations: []string{"texas-usa"}},
		},
	}
	tmpl := mustParseTemplate(`{{range .LocationGroups}}{{.Country}}{{end}}`)
	h := NewHomeHandler(s, tmpl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Usa") {
		t.Errorf("body does not contain %q\nbody: %s", "Usa", rec.Body.String())
	}
}

func TestHomeHandler_PassesBoundsToTemplate(t *testing.T) {
	s := &testStore{
		bounds: models.FilterBounds{CreationMin: 1958, CreationMax: 2015},
	}
	tmpl := mustParseTemplate(`{{.Bounds.CreationMin}}-{{.Bounds.CreationMax}}`)
	h := NewHomeHandler(s, tmpl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "1958-2015") {
		t.Errorf("body does not contain bounds %q\nbody: %s", "1958-2015", rec.Body.String())
	}
}

func TestArtistHandler(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Foo Fighters", Image: "http://img/1.jpg", CreationDate: 1994, FirstAlbum: "04-07-1995"},
	}

	artistTmpl := mustParseTemplate(`{{.Artist.Name}}`)

	tests := []struct {
		name             string
		url              string
		pathID           string
		artists          []models.Artist
		tmpl             *template.Template
		wantStatusCode   int
		wantBodyContains []string
	}{
		{
			name:             "valid_id_returns_200",
			url:              "/artist/1",
			pathID:           "1",
			artists:          artists,
			tmpl:             artistTmpl,
			wantStatusCode:   http.StatusOK,
			wantBodyContains: []string{"Foo Fighters"},
		},
		{
			name:             "unknown_id_returns_404",
			url:              "/artist/99",
			pathID:           "99",
			artists:          artists,
			tmpl:             artistTmpl,
			wantStatusCode:   http.StatusNotFound,
			wantBodyContains: []string{},
		},
		{
			name:             "non_numeric_id_returns_404",
			url:              "/artist/abc",
			pathID:           "abc",
			artists:          artists,
			tmpl:             artistTmpl,
			wantStatusCode:   http.StatusNotFound,
			wantBodyContains: []string{},
		},
		{
			name:             "markers_json_present_in_body",
			url:              "/artist/1",
			pathID:           "1",
			artists:          artists,
			tmpl:             mustParseTemplate(`{{.MarkersJSON}}`),
			wantStatusCode:   http.StatusOK,
			wantBodyContains: []string{`london-uk`},
		},
		{
			name:             "template_error_returns_500",
			url:              "/artist/1",
			pathID:           "1",
			artists:          artists,
			tmpl:             brokenTemplate(),
			wantStatusCode:   http.StatusInternalServerError,
			wantBodyContains: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &testStore{artists: tc.artists}
			h := NewArtistHandler(s, tc.tmpl, mustParseTemplate(`not found`))

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			req.SetPathValue("id", tc.pathID)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatusCode)
			}

			body := rec.Body.String()
			for _, want := range tc.wantBodyContains {
				if !strings.Contains(body, want) {
					t.Errorf("body does not contain %q\nbody: %s", want, body)
				}
			}
		})
	}
}

func TestSearchHandler(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Queen"},
		{ID: 2, Name: "Billie Eilish"},
	}

	tests := []struct {
		name           string
		method         string
		url            string
		wantStatusCode int
		wantInBody     string
	}{
		{
			name:           "valid_query_returns_200",
			method:         http.MethodGet,
			url:            "/api/search?q=queen",
			wantStatusCode: http.StatusOK,
			wantInBody:     "Queen",
		},
		{
			name:           "missing_q_returns_400",
			method:         http.MethodGet,
			url:            "/api/search",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "empty_q_returns_400",
			method:         http.MethodGet,
			url:            "/api/search?q=",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "post_method_returns_405",
			method:         http.MethodPost,
			url:            "/api/search?q=queen",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
	}

	badReqTmpl := template.Must(template.New("400.html").Parse(`Bad Request`))

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := &SearchHandler{Store: &testStore{artists: artists}, BadRequestTmpl: badReqTmpl}
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rec := httptest.NewRecorder()

			h.Search(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatusCode)
			}
			if tc.wantInBody != "" && !strings.Contains(rec.Body.String(), tc.wantInBody) {
				t.Errorf("body does not contain %q\nbody: %s", tc.wantInBody, rec.Body.String())
			}
		})
	}
}

func TestSearchHandler_NoMatchReturnsEmptyArray(t *testing.T) {
	h := &SearchHandler{
		Store:          &store.MockStore{},
		BadRequestTmpl: template.Must(template.New("400.html").Parse(`Bad Request`)),
	}
	req := httptest.NewRequest(http.MethodGet, "/api/search?q=zzznomatchxyz", nil)
	rec := httptest.NewRecorder()

	h.Search(rec, req)

	body := strings.TrimSpace(rec.Body.String())
	if body != "[]" {
		t.Errorf("body = %q, want %q", body, "[]")
	}
}

func TestFilterHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		wantStatusCode int
		wantInBody     string
	}{
		{
			name:           "no_params_returns_all_artists",
			method:         http.MethodGet,
			url:            "/api/filter",
			wantStatusCode: http.StatusOK,
			wantInBody:     "Billie Eilish",
		},
		{
			name:           "query_only_matches_subset",
			method:         http.MethodGet,
			url:            "/api/filter?q=billie",
			wantStatusCode: http.StatusOK,
			wantInBody:     "Billie Eilish",
		},
		{
			name:           "members_max_matches_zero_member_fixtures",
			method:         http.MethodGet,
			url:            "/api/filter?members_max=0",
			wantStatusCode: http.StatusOK,
			wantInBody:     "System of a down",
		},
		{
			name:           "members_min_excludes_zero_member_fixtures",
			method:         http.MethodGet,
			url:            "/api/filter?members_min=1",
			wantStatusCode: http.StatusOK,
			wantInBody:     "[]",
		},
		{
			name:           "malformed_creation_min_returns_400",
			method:         http.MethodGet,
			url:            "/api/filter?creation_min=notanumber",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "malformed_creation_max_returns_400",
			method:         http.MethodGet,
			url:            "/api/filter?creation_max=abc",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "malformed_first_album_min_returns_400",
			method:         http.MethodGet,
			url:            "/api/filter?first_album_min=xyz",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "malformed_first_album_max_returns_400",
			method:         http.MethodGet,
			url:            "/api/filter?first_album_max=xyz",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "malformed_members_min_returns_400",
			method:         http.MethodGet,
			url:            "/api/filter?members_min=xyz",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "malformed_members_max_returns_400",
			method:         http.MethodGet,
			url:            "/api/filter?members_max=xyz",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "locations_param_matches_nothing_in_mock_store",
			method:         http.MethodGet,
			url:            "/api/filter?locations=texas-usa",
			wantStatusCode: http.StatusOK,
			wantInBody:     "[]",
		},
		{
			name:           "post_method_returns_405",
			method:         http.MethodPost,
			url:            "/api/filter",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
	}

	badReqTmpl := template.Must(template.New("400.html").Parse(`Bad Request`))

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := &FilterHandler{Store: &store.MockStore{}, BadRequestTmpl: badReqTmpl}
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rec := httptest.NewRecorder()

			h.Filter(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatusCode)
			}
			if tc.wantInBody != "" && !strings.Contains(rec.Body.String(), tc.wantInBody) {
				t.Errorf("body does not contain %q\nbody: %s", tc.wantInBody, rec.Body.String())
			}
		})
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	errTmpl := template.Must(template.New("500.html").Parse(`Internal Server Error`))

	tests := []struct {
		name           string
		handler        http.HandlerFunc
		wantStatusCode int
	}{
		{
			name: "panic_returns_500",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "no_panic_passes_through",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			RecoveryMiddleware(errTmpl, tc.handler).ServeHTTP(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatusCode)
			}
		})
	}
}

func TestErrorHandlers(t *testing.T) {
	tests := []struct {
		name           string
		templateName   string
		templateBody   string
		handler        func(*template.Template) http.HandlerFunc
		wantStatusCode int
	}{
		{
			name:           "bad_request_returns_400",
			templateName:   "400.html",
			templateBody:   `Bad Request`,
			handler:        BadRequestHandler,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "not_found_returns_404",
			templateName:   "404.html",
			templateBody:   `Not Found`,
			handler:        NotFoundHandler,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "internal_server_error_returns_500",
			templateName:   "500.html",
			templateBody:   `Internal Server Error`,
			handler:        StatusInternalServerError,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpl := template.Must(template.New(tc.templateName).Parse(tc.templateBody))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			tc.handler(tmpl)(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatusCode)
			}
		})
	}
}
