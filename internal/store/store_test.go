package store

import (
	"groupie-tracker-geolocalization/internal/models"
	"strings"
	"testing"
)

func TestFirstAlbumYear(t *testing.T) {
	tests := []struct {
		name      string
		firstAlbum string
		wantYear  int
		wantErr   bool
	}{
		{"valid_date", "14-12-1973", 1973, false},
		{"valid_date_single_digit_day_month", "4-7-1995", 1995, false},
		{"empty_string", "", 0, true},
		{"missing_parts", "14-12", 0, true},
		{"non_numeric_year", "14-12-abcd", 0, true},
		{"extra_parts", "14-12-1973-extra", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			year, err := firstAlbumYear(tc.firstAlbum)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && year != tc.wantYear {
				t.Errorf("year = %d, want %d", year, tc.wantYear)
			}
		})
	}
}

func TestFilterCriteria_ZeroValueHasNoConstraints(t *testing.T) {
	var c FilterCriteria
	if c.CreationDateMin != nil || c.CreationDateMax != nil {
		t.Error("expected zero-value FilterCriteria to have nil CreationDate bounds")
	}
	if c.FirstAlbumMin != nil || c.FirstAlbumMax != nil {
		t.Error("expected zero-value FilterCriteria to have nil FirstAlbum bounds")
	}
	if c.MembersMin != nil || c.MembersMax != nil {
		t.Error("expected zero-value FilterCriteria to have nil Members bounds")
	}
	if c.Locations != nil {
		t.Error("expected zero-value FilterCriteria to have nil Locations")
	}
}

func TestMockStore_FilterArtists_NoConstraintsReturnsAll(t *testing.T) {
	m := &MockStore{}
	got := m.FilterArtists("", FilterCriteria{})
	if len(got) != len(m.AllArtists()) {
		t.Errorf("expected all %d artists, got %d", len(m.AllArtists()), len(got))
	}
}

func TestMockStore_AllArtists(t *testing.T) {
	m := &MockStore{}
	artists := m.AllArtists()
	if len(artists) != 2 {
		t.Errorf("expected 2 artists, got %d", len(artists))
	}
}

func TestMockStore_ArtistByID_Found(t *testing.T) {
	m := &MockStore{}
	a, ok := m.ArtistByID(1)
	if !ok {
		t.Fatal("expected artist to be found")
	}
	if a.Name != "Billie Eilish" {
		t.Errorf("expected Billie Eilish, got %s", a.Name)
	}
}

func TestMockStore_ArtistByID_NotFound(t *testing.T) {
	m := &MockStore{}
	_, ok := m.ArtistByID(99)
	if ok {
		t.Error("expected false for missing ID")
	}
}

func TestMockStore_SearchArtists(t *testing.T) {
	m := &MockStore{}

	results := m.SearchArtists("Billie")
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	results = m.SearchArtists("Queen")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}

	results = m.SearchArtists("system")
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}
func TestRealStore_ArtistPageDataByID(t *testing.T) {
	s := &RealStore{
		Artists: []models.Artist{
			{ID: 1, Name: "Queen", Members: []string{"Freddie Mercury"}, CreationDate: 1970, FirstAlbum: "14-12-1973"},
			{ID: 2, Name: "Billie Eilish", Members: []string{"Billie Eilish"}, CreationDate: 2015, FirstAlbum: "26-03-2017"},
		},
		Locations: models.LocationsResponse{
			Index: []models.Locations{{ID: 1, Locations: []string{"london-uk", "paris-france"}}},
		},
		Dates: models.DatesResponse{
			Index: []models.Dates{{ID: 1, Dates: []string{"*06-03-2020", "07-03-2020"}}},
		},
		Relations: models.RelationResponse{
			Index: []models.Relation{{ID: 1, DatesLocations: map[string][]string{
				"london-uk":    {"*06-03-2020", "07-03-2020"},
				"paris-france": {"10-03-2020"},
			}}},
		},
		Markers: map[int][]models.Marker{
			1: {
				{Name: "london-uk", Lat: 51.5074, Lng: -0.1278},
				{Name: "paris-france", Lat: 48.8566, Lng: 2.3522},
			},
		},
	}

	tests := []struct {
		name               string
		id                 int
		wantFound          bool
		wantLocations      []string
		wantDatesLocations map[string][]string
		wantMarkersJSON    string
	}{
		{
			name:          "known_id_returns_data",
			id:            1,
			wantFound:     true,
			wantLocations: []string{"london-uk", "paris-france"},
			wantDatesLocations: map[string][]string{
				"london-uk":    {"*06-03-2020", "07-03-2020"},
				"paris-france": {"10-03-2020"},
			},
			wantMarkersJSON: "london-uk",
		},
		{
			name:      "unknown_id_returns_false",
			id:        99,
			wantFound: false,
		},
		{
			name:               "known_id_no_markers_returns_empty_json",
			id:                 2,
			wantFound:          true,
			wantLocations:      []string{},
			wantDatesLocations: map[string][]string{},
			wantMarkersJSON:    "[]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, ok := s.ArtistPageDataByID(tc.id)
			if ok != tc.wantFound {
				t.Fatalf("found = %v, want %v", ok, tc.wantFound)
			}
			if !tc.wantFound {
				return
			}
			if len(data.Locations) != len(tc.wantLocations) {
				t.Errorf("locations = %v, want %v", data.Locations, tc.wantLocations)
			}
			if len(data.DatesLocations) != len(tc.wantDatesLocations) {
				t.Errorf("datesLocations len = %d, want %d", len(data.DatesLocations), len(tc.wantDatesLocations))
			}
			for loc, dates := range tc.wantDatesLocations {
				got, exists := data.DatesLocations[loc]
				if !exists {
					t.Errorf("missing location %q in DatesLocations", loc)
					continue
				}
				if len(got) != len(dates) {
					t.Errorf("DatesLocations[%q] = %v, want %v", loc, got, dates)
				}
			}
			if len(data.MarkersJSON) == 0 {
				t.Error("expected MarkersJSON to be non-empty")
			}
			if !strings.Contains(string(data.MarkersJSON), tc.wantMarkersJSON) {
				t.Errorf("MarkersJSON does not contain %q: %s", tc.wantMarkersJSON, data.MarkersJSON)
			}
		})
	}
}

func TestRealStore_SearchArtists(t *testing.T) {
	s := &RealStore{
		Artists: []models.Artist{
			{ID: 1, Name: "Queen", Members: []string{"Freddie Mercury", "Brian May"}, CreationDate: 1970, FirstAlbum: "14-12-1973"},
			{ID: 2, Name: "Billie Eilish", Members: []string{"Billie Eilish"}, CreationDate: 2015, FirstAlbum: "26-03-2017"},
		},
	}

	tests := []struct {
		name      string
		query     string
		wantNames []string
	}{
		{"match_by_name", "queen", []string{"Queen"}},
		{"match_by_name_case_insensitive", "BILLIE", []string{"Billie Eilish"}},
		{"match_by_member", "freddie", []string{"Queen"}},
		{"match_by_creation_date", "1970", []string{"Queen"}},
		{"match_by_first_album", "26-03-2017", []string{"Billie Eilish"}},
		{"no_match_returns_empty", "zzznomatch", []string{}},
		{"match_multiple", "2", []string{"Queen", "Billie Eilish"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results := s.SearchArtists(tc.query)
			if len(results) != len(tc.wantNames) {
				t.Fatalf("got %d results, want %d", len(results), len(tc.wantNames))
			}
			for i, name := range tc.wantNames {
				if results[i].Name != name {
					t.Errorf("result[%d] = %q, want %q", i, results[i].Name, name)
				}
			}
		})
	}
}
