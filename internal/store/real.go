package store

import (
	"encoding/json"
	"groupie-tracker-filters/internal/models"
	"html/template"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// RealStore is the production implementation of the Store interface.
// It holds all data loaded from the external API at startup and serves it
// directly from memory on every request, avoiding repeated network calls.
type RealStore struct {
	Artists   []models.Artist
	Locations models.LocationsResponse
	Dates     models.DatesResponse
	Relations models.RelationResponse
	Markers   map[int][]models.Marker
}

// AllArtists returns the full list of artists held in memory.
func (r *RealStore) AllArtists() []models.Artist {
	return r.Artists
}

// ArtistByID performs a linear scan over all artists and returns the one
// matching the given ID. Returns false if no match is found.
func (r *RealStore) ArtistByID(id int) (models.Artist, bool) {
	for _, a := range r.AllArtists() {
		if a.ID == id {
			return a, true
		}
	}
	return models.Artist{}, false
}

// SearchArtists filters the artist list using matchesQuery and returns all
// artists that contain the query string in any searchable field.
func (r *RealStore) SearchArtists(query string) []models.Artist {
	var result []models.Artist
	for _, a := range r.AllArtists() {
		if matchesQuery(a, query) {
			result = append(result, a)
		}
	}
	return result
}

// filterMatch pairs a matching artist with its original index in AllArtists,
// so concurrent workers in FilterArtists can report results out of order
// and still let the caller restore a deterministic, input-order result.
type filterMatch struct {
	index  int
	artist models.Artist
}

// FilterArtists returns all artists matching both query (via matchesQuery)
// and criteria (via matchesCriteria), combined with a logical AND. Matching
// is fanned out across a bounded pool of goroutines — one per available CPU —
// since evaluating each artist is independent, read-only work over data that
// never changes after startup. Workers report matches on a channel in
// completion order, which is not guaranteed to match input order, so results
// are sorted back into the original AllArtists order before returning —
// this ordering step must not be removed, or repeated identical filter
// requests could return artists in a different order on each call.
func (r *RealStore) FilterArtists(query string, criteria FilterCriteria) []models.Artist {
	artists := r.AllArtists()

	jobs := make(chan int)
	matches := make(chan filterMatch)

	var wg sync.WaitGroup
	workerCount := runtime.NumCPU()
	if workerCount > len(artists) {
		workerCount = len(artists)
	}
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				a := artists[i]
				if matchesQuery(a, query) && matchesCriteria(a, r.locationsForArtist(a.ID), criteria) {
					matches <- filterMatch{index: i, artist: a}
				}
			}
		}()
	}

	go func() {
		for i := range artists {
			jobs <- i
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(matches)
	}()

	var result []filterMatch
	for m := range matches {
		result = append(result, m)
	}

	sort.Slice(result, func(i, j int) bool { return result[i].index < result[j].index })

	artistsResult := make([]models.Artist, len(result))
	for i, m := range result {
		artistsResult[i] = m.artist
	}
	return artistsResult
}

// LocationGroups scans every artist's concert locations, groups the unique
// slugs by their trailing country segment (e.g. "usa" from "texas-usa"),
// and returns one LocationGroup per country. Locations within a group and
// the groups themselves are both sorted alphabetically by country, so the
// checkbox filter renders in a stable, predictable order across requests.
func (r *RealStore) LocationGroups() []models.LocationGroup {
	byCountry := make(map[string]map[string]struct{})

	for _, entry := range r.Locations.Index {
		for _, loc := range entry.Locations {
			country := countrySlug(loc)
			if byCountry[country] == nil {
				byCountry[country] = make(map[string]struct{})
			}
			byCountry[country][loc] = struct{}{}
		}
	}

	groups := make([]models.LocationGroup, 0, len(byCountry))
	for country, locSet := range byCountry {
		locs := make([]string, 0, len(locSet))
		for loc := range locSet {
			locs = append(locs, loc)
		}
		sort.Strings(locs)
		groups = append(groups, models.LocationGroup{
			Country:   titleCaseCountry(country),
			Locations: locs,
		})
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i].Country < groups[j].Country })

	return groups
}

// FilterBounds scans every artist and returns the minimum and maximum
// creation year, first-album year, and member count across the dataset,
// used to seed the range slider endpoints. Artists whose FirstAlbum string
// cannot be parsed are skipped for the first-album bounds only, so one
// malformed date does not distort the slider range. Returns a zero-value
// FilterBounds if there are no artists.
func (r *RealStore) FilterBounds() models.FilterBounds {
	artists := r.AllArtists()
	if len(artists) == 0 {
		return models.FilterBounds{}
	}

	var b models.FilterBounds
	firstAlbumSeen := false

	for i, a := range artists {
		members := len(a.Members)

		if i == 0 {
			b.CreationMin, b.CreationMax = a.CreationDate, a.CreationDate
			b.MembersMin, b.MembersMax = members, members
		} else {
			b.CreationMin = min(b.CreationMin, a.CreationDate)
			b.CreationMax = max(b.CreationMax, a.CreationDate)
			b.MembersMin = min(b.MembersMin, members)
			b.MembersMax = max(b.MembersMax, members)
		}

		if year, err := firstAlbumYear(a.FirstAlbum); err == nil {
			if !firstAlbumSeen {
				b.FirstAlbumMin, b.FirstAlbumMax = year, year
				firstAlbumSeen = true
			} else {
				b.FirstAlbumMin = min(b.FirstAlbumMin, year)
				b.FirstAlbumMax = max(b.FirstAlbumMax, year)
			}
		}
	}

	return b
}

// countrySlug extracts the country segment from a "city-country" location
// slug — the substring after the last hyphen.
func countrySlug(location string) string {
	idx := strings.LastIndex(location, "-")
	if idx == -1 {
		return location
	}
	return location[idx+1:]
}

// titleCaseCountry converts a country slug (e.g. "new_zealand") into a
// human-readable label (e.g. "New Zealand") by replacing underscores with
// spaces and capitalising the first letter of each word.
func titleCaseCountry(slug string) string {
	words := strings.Split(slug, "_")
	for i, w := range words {
		if w == "" {
			continue
		}
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}

// locationsForArtist returns the concert location slugs for the artist with
// the given ID, or nil if the artist has no entry in the Locations index.
func (r *RealStore) locationsForArtist(id int) []string {
	for _, l := range r.Locations.Index {
		if l.ID == id {
			return l.Locations
		}
	}
	return nil
}

// matchesQuery performs a case-insensitive substring search across the artist's
// name, individual member names, creation date, and first album date.
// It returns true as soon as any field matches, without checking the rest.
func matchesQuery(a models.Artist, query string) bool {
	q := strings.ToLower(query)
	if strings.Contains(strings.ToLower(a.Name), q) {
		return true
	}
	for _, member := range a.Members {
		if strings.Contains(strings.ToLower(member), q) {
			return true
		}
	}
	if strings.Contains(strings.ToLower(strconv.Itoa(a.CreationDate)), q) {
		return true
	}
	if strings.Contains(strings.ToLower(a.FirstAlbum), q) {
		return true
	}

	return false
}

// ArtistPageDataByID looks up an artist by ID and assembles the ArtistPageData
// view model by joining locations, dates, and relations data from their respective
// index slices, and attaching the pre-geocoded markers from the Markers map.
// Returns false if no artist with the given ID exists.
func (r *RealStore) ArtistPageDataByID(id int) (models.ArtistPageData, bool) {
	for _, a := range r.AllArtists() {
		if a.ID == id {
			locations := r.locationsForArtist(id)
			var dates []string
			for _, d := range r.Dates.Index {
				if d.ID == id {
					dates = d.Dates
				}
			}
			var datesLocations map[string][]string
			for _, rel := range r.Relations.Index {
				if rel.ID == id {
					datesLocations = rel.DatesLocations
				}
			}
			markers := r.Markers[id]
			if markers == nil {
				markers = []models.Marker{}
			}
			markersJSON, err := json.Marshal(markers)
			if err != nil {
				markersJSON = []byte("[]")
			}
			return models.ArtistPageData{
				Artist:         a,
				Locations:      locations,
				Dates:          dates,
				DatesLocations: datesLocations,
				Markers:        markers,
				MarkersJSON:    template.JS(markersJSON),
			}, true
		}
	}
	return models.ArtistPageData{}, false
}
