package models

import (
	"encoding/json"
	"testing"
)

func TestArtistJSONRoundTrip(t *testing.T) {
	a := Artist{ID: 1, Name: "Test Band", Members: []string{"Alice", "Bob"}, CreationDate: 2000, FirstAlbum: "01-01-2000"}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var got Artist
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != a.ID || got.Name != a.Name || len(got.Members) != len(a.Members) {
		t.Errorf("got %+v, want %+v", got, a)
	}
}

func TestRelationJSONRoundTrip(t *testing.T) {
	r := Relation{ID: 1, DatesLocations: map[string][]string{"paris": {"2024-01-01"}}}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var got Relation
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != r.ID || len(got.DatesLocations["paris"]) != 1 {
		t.Errorf("got %+v", got)
	}
}

func TestLocationsJSONRoundTrip(t *testing.T) {
	l := Locations{ID: 2, Locations: []string{"london", "berlin"}}
	data, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	var got Locations
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != l.ID || len(got.Locations) != 2 {
		t.Errorf("got %+v", got)
	}
}

func TestDatesJSONRoundTrip(t *testing.T) {
	d := Dates{ID: 3, Dates: []string{"2024-06-01", "2024-07-15"}}
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	var got Dates
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != d.ID || len(got.Dates) != 2 {
		t.Errorf("got %+v", got)
	}
}
func TestArtistPageData_EmptySlices(t *testing.T) {
	p := ArtistPageData{
		Artist:    Artist{ID: 1, Name: "Test"},
		Locations: []string{},
		Dates:     []string{},
	}
	if len(p.Locations) != 0 {
		t.Errorf("expected empty locations")
	}
	if len(p.Dates) != 0 {
		t.Errorf("expected empty dates")
	}
}

func TestArtistPageData_NilSlices(t *testing.T) {
	p := ArtistPageData{}
	if p.Locations != nil {
		t.Errorf("expected nil locations")
	}
	if p.Dates != nil {
		t.Errorf("expected nil dates")
	}
}
