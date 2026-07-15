package store

import (
	"groupie-tracker-filters/internal/models"
	"strings"
	"testing"
)

func TestFirstAlbumYear(t *testing.T) {
	tests := []struct {
		name       string
		firstAlbum string
		wantYear   int
		wantErr    bool
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

func TestMockStore_FilterArtists_AppliesCriteria(t *testing.T) {
	m := &MockStore{}

	tests := []struct {
		name     string
		query    string
		criteria FilterCriteria
		wantLen  int
	}{
		{"members_min_excludes_zero_member_fixtures", "", FilterCriteria{MembersMin: intPtr(1)}, 0},
		{"members_max_zero_matches_fixtures", "", FilterCriteria{MembersMax: intPtr(0)}, 2},
		{"location_constraint_matches_nothing", "", FilterCriteria{Locations: []string{"texas-usa"}}, 0},
		{"query_and_criteria_combined", "billie", FilterCriteria{MembersMax: intPtr(0)}, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := m.FilterArtists(tc.query, tc.criteria)
			if len(got) != tc.wantLen {
				t.Errorf("got %d results, want %d", len(got), tc.wantLen)
			}
		})
	}
}

func intPtr(v int) *int { return &v }

func TestMatchesCriteria(t *testing.T) {
	artist := models.Artist{
		ID:           1,
		Name:         "Foo Fighters",
		Members:      []string{"Dave", "Nate", "Pat", "Chris", "Rami", "Josh"},
		CreationDate: 1994,
		FirstAlbum:   "04-07-1995",
	}
	artistLocations := []string{"texas-usa", "washington-usa"}

	tests := []struct {
		name      string
		locations []string
		criteria  FilterCriteria
		want      bool
	}{
		{"zero_value_matches_everything", artistLocations, FilterCriteria{}, true},
		{"creation_date_within_range", artistLocations, FilterCriteria{CreationDateMin: intPtr(1990), CreationDateMax: intPtr(2000)}, true},
		{"creation_date_below_min", artistLocations, FilterCriteria{CreationDateMin: intPtr(1995)}, false},
		{"creation_date_above_max", artistLocations, FilterCriteria{CreationDateMax: intPtr(1993)}, false},
		{"first_album_within_range", artistLocations, FilterCriteria{FirstAlbumMin: intPtr(1990), FirstAlbumMax: intPtr(2000)}, true},
		{"first_album_below_min", artistLocations, FilterCriteria{FirstAlbumMin: intPtr(1996)}, false},
		{"first_album_above_max", artistLocations, FilterCriteria{FirstAlbumMax: intPtr(1994)}, false},
		{"members_exact_match", artistLocations, FilterCriteria{MembersMin: intPtr(6), MembersMax: intPtr(6)}, true},
		{"members_below_min", artistLocations, FilterCriteria{MembersMin: intPtr(7)}, false},
		{"members_above_max", artistLocations, FilterCriteria{MembersMax: intPtr(5)}, false},
		{"location_match", artistLocations, FilterCriteria{Locations: []string{"texas-usa"}}, true},
		{"location_no_match", artistLocations, FilterCriteria{Locations: []string{"california-usa"}}, false},
		{"location_match_any_of_multiple", artistLocations, FilterCriteria{Locations: []string{"california-usa", "washington-usa"}}, true},
		{"broad_filter_matches_more_specific_slug", []string{"seattle-washington-usa"}, FilterCriteria{Locations: []string{"washington-usa"}}, true},
		{"location_constraint_no_locations_data", nil, FilterCriteria{Locations: []string{"texas-usa"}}, false},
		{
			"all_criteria_combined_match",
			artistLocations,
			FilterCriteria{
				CreationDateMin: intPtr(1990), CreationDateMax: intPtr(2000),
				FirstAlbumMin: intPtr(1990), FirstAlbumMax: intPtr(2000),
				MembersMin: intPtr(6), MembersMax: intPtr(6),
				Locations: []string{"texas-usa"},
			},
			true,
		},
		{
			"all_criteria_combined_one_fails",
			artistLocations,
			FilterCriteria{
				CreationDateMin: intPtr(1990), CreationDateMax: intPtr(2000),
				MembersMin: intPtr(7),
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := matchesCriteria(artist, tc.locations, tc.criteria)
			if got != tc.want {
				t.Errorf("matchesCriteria() = %v, want %v", got, tc.want)
			}
		})
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

func filterFixtureStore() *RealStore {
	return &RealStore{
		Artists: []models.Artist{
			{ID: 1, Name: "SOJA", Members: []string{"a", "b", "c", "d", "e", "f", "g", "h"}, CreationDate: 1997, FirstAlbum: "01-01-2002"},
			{ID: 2, Name: "Pearl Jam", Members: []string{"a", "b", "c", "d", "e"}, CreationDate: 1990, FirstAlbum: "03-07-1992"},
			{ID: 3, Name: "Red Hot Chili Peppers", Members: []string{"a", "b", "c", "d"}, CreationDate: 1982, FirstAlbum: "10-02-1990"},
			{ID: 4, Name: "Pink Floyd", Members: []string{"a", "b", "c", "d", "e", "f"}, CreationDate: 1965, FirstAlbum: "05-08-1967"},
			{ID: 5, Name: "Foo Fighters", Members: []string{"a", "b", "c", "d", "e", "f"}, CreationDate: 1994, FirstAlbum: "04-07-1995"},
			{ID: 6, Name: "Bobby McFerrins", Members: []string{"a"}, CreationDate: 1977, FirstAlbum: "01-01-1979"},
			{ID: 7, Name: "Eminem", Members: []string{"a"}, CreationDate: 1996, FirstAlbum: "12-11-1996"},
			{ID: 8, Name: "Twenty One Pilots", Members: []string{"a", "b"}, CreationDate: 2009, FirstAlbum: "05-12-2009"},
			{ID: 9, Name: "The Rolling Stones", Members: []string{"a", "b", "c", "d"}, CreationDate: 1962, FirstAlbum: "01-01-1964"},
			{ID: 10, Name: "Metallica", Members: []string{"a", "b", "c", "d"}, CreationDate: 1981, FirstAlbum: "22-08-1987"},
			{ID: 11, Name: "Post Malone", Members: []string{"a"}, CreationDate: 2013, FirstAlbum: "09-12-2016"},
			{ID: 12, Name: "Early Debut Band", Members: []string{"a", "b"}, CreationDate: 2012, FirstAlbum: "01-01-2009"},
		},
		Locations: models.LocationsResponse{
			Index: []models.Locations{
				{ID: 8, Locations: []string{"texas-usa", "oklahoma-usa"}},
				{ID: 9, Locations: []string{"washington-usa", "california-usa"}},
			},
		},
	}
}

func TestRealStore_FilterArtists(t *testing.T) {
	s := filterFixtureStore()

	tests := []struct {
		name      string
		query     string
		criteria  FilterCriteria
		wantNames []string
	}{
		{
			name:      "creation_date_range_only",
			criteria:  FilterCriteria{CreationDateMin: intPtr(1995), CreationDateMax: intPtr(2000)},
			wantNames: []string{"SOJA", "Eminem"},
		},
		{
			name:      "first_album_range_only",
			criteria:  FilterCriteria{FirstAlbumMin: intPtr(1990), FirstAlbumMax: intPtr(1992)},
			wantNames: []string{"Pearl Jam", "Red Hot Chili Peppers"},
		},
		{
			name:      "exact_member_count",
			criteria:  FilterCriteria{MembersMin: intPtr(6), MembersMax: intPtr(6)},
			wantNames: []string{"Pink Floyd", "Foo Fighters"},
		},
		{
			name:      "location_only",
			criteria:  FilterCriteria{Locations: []string{"texas-usa"}},
			wantNames: []string{"Twenty One Pilots"},
		},
		{
			name: "creation_date_range_and_solo_artist",
			criteria: FilterCriteria{
				CreationDateMin: intPtr(1970), CreationDateMax: intPtr(2000),
				MembersMin: intPtr(1), MembersMax: intPtr(1),
			},
			wantNames: []string{"Bobby McFerrins", "Eminem"},
		},
		{
			name: "location_and_members_above",
			criteria: FilterCriteria{
				Locations:  []string{"washington-usa"},
				MembersMin: intPtr(4),
			},
			wantNames: []string{"The Rolling Stones"},
		},
		{
			name: "first_album_range_and_max_members",
			criteria: FilterCriteria{
				FirstAlbumMin: intPtr(1980), FirstAlbumMax: intPtr(1990),
				MembersMax: intPtr(4),
			},
			wantNames: []string{"Red Hot Chili Peppers", "Metallica"},
		},
		{
			name: "creation_after_2010_and_first_album_after_2010",
			criteria: FilterCriteria{
				CreationDateMin: intPtr(2011),
				FirstAlbumMin:   intPtr(2011),
			},
			wantNames: []string{"Post Malone"},
		},
		{
			name:      "no_constraints_returns_all",
			criteria:  FilterCriteria{},
			wantNames: []string{"SOJA", "Pearl Jam", "Red Hot Chili Peppers", "Pink Floyd", "Foo Fighters", "Bobby McFerrins", "Eminem", "Twenty One Pilots", "The Rolling Stones", "Metallica", "Post Malone", "Early Debut Band"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results := s.FilterArtists(tc.query, tc.criteria)
			gotNames := make(map[string]bool, len(results))
			for _, a := range results {
				gotNames[a.Name] = true
			}
			if len(results) != len(tc.wantNames) {
				t.Fatalf("got %d results (%v), want %d (%v)", len(results), gotNames, len(tc.wantNames), tc.wantNames)
			}
			for _, name := range tc.wantNames {
				if !gotNames[name] {
					t.Errorf("expected %q in results, got %v", name, results)
				}
			}
		})
	}
}

func TestRealStore_FilterArtists_StableOrder(t *testing.T) {
	s := filterFixtureStore()

	first := s.FilterArtists("", FilterCriteria{})
	for i := 0; i < 20; i++ {
		got := s.FilterArtists("", FilterCriteria{})
		if len(got) != len(first) {
			t.Fatalf("run %d: got %d results, want %d", i, len(got), len(first))
		}
		for j := range first {
			if got[j].ID != first[j].ID {
				t.Fatalf("run %d: order mismatch at index %d: got ID %d, want ID %d", i, j, got[j].ID, first[j].ID)
			}
		}
	}
}

func TestRealStore_LocationGroups(t *testing.T) {
	s := &RealStore{
		Locations: models.LocationsResponse{
			Index: []models.Locations{
				{ID: 1, Locations: []string{"texas-usa", "washington-usa"}},
				{ID: 2, Locations: []string{"london-uk", "california-usa"}},
				{ID: 3, Locations: []string{"auckland-new_zealand"}},
			},
		},
	}

	groups := s.LocationGroups()

	byCountry := make(map[string][]string, len(groups))
	for _, g := range groups {
		byCountry[g.Country] = g.Locations
	}

	if len(groups) != 3 {
		t.Fatalf("got %d groups, want 3", len(groups))
	}

	usa, ok := byCountry["Usa"]
	if !ok {
		t.Fatal("expected a Usa group")
	}
	wantUSA := []string{"california-usa", "texas-usa", "washington-usa"}
	if len(usa) != len(wantUSA) {
		t.Fatalf("Usa locations = %v, want %v", usa, wantUSA)
	}
	for i, loc := range wantUSA {
		if usa[i] != loc {
			t.Errorf("Usa[%d] = %q, want %q", i, usa[i], loc)
		}
	}

	uk, ok := byCountry["Uk"]
	if !ok || len(uk) != 1 || uk[0] != "london-uk" {
		t.Errorf("Uk group = %v, want [london-uk]", uk)
	}

	nz, ok := byCountry["New Zealand"]
	if !ok || len(nz) != 1 || nz[0] != "auckland-new_zealand" {
		t.Errorf("New Zealand group = %v, want [auckland-new_zealand]", nz)
	}
}

func TestRealStore_LocationGroups_SortedByCountry(t *testing.T) {
	s := &RealStore{
		Locations: models.LocationsResponse{
			Index: []models.Locations{
				{ID: 1, Locations: []string{"paris-france", "tokyo-japan", "london-uk"}},
			},
		},
	}

	groups := s.LocationGroups()

	wantOrder := []string{"France", "Japan", "Uk"}
	if len(groups) != len(wantOrder) {
		t.Fatalf("got %d groups, want %d", len(groups), len(wantOrder))
	}
	for i, country := range wantOrder {
		if groups[i].Country != country {
			t.Errorf("groups[%d].Country = %q, want %q", i, groups[i].Country, country)
		}
	}
}

func TestMockStore_LocationGroups_ReturnsEmpty(t *testing.T) {
	m := &MockStore{}
	groups := m.LocationGroups()
	if len(groups) != 0 {
		t.Errorf("expected no location groups from MockStore, got %v", groups)
	}
}
