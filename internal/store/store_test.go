package store

import (
	"groupie-tracker/internal/models"
	"testing"
)

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
	}

	tests := []struct {
		name               string
		id                 int
		wantFound          bool
		wantLocations      []string
		wantDatesLocations map[string][]string
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
		},
		{
			name:      "unknown_id_returns_false",
			id:        99,
			wantFound: false,
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
