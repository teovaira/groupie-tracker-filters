package store

import (
	"fmt"
	"groupie-tracker-filters/internal/models"
	"strconv"
	"strings"
)

// FilterCriteria holds the optional constraints used by Store.FilterArtists.
// Each bound is a pointer so that "unset" (nil) can be distinguished from a
// genuine zero value — e.g. a nil MembersMax means no upper bound, whereas
// MembersMax pointing at 0 would mean "at most zero members". A nil or empty
// Locations means no location constraint; a non-empty Locations matches an
// artist if any of its concert locations matches any entry (logical OR).
type FilterCriteria struct {
	CreationDateMin *int
	CreationDateMax *int
	FirstAlbumMin   *int
	FirstAlbumMax   *int
	MembersMin      *int
	MembersMax      *int
	Locations       []string
}

// matchesCriteria reports whether an artist satisfies every constraint set on
// c. locations holds the artist's resolved concert location slugs (joined by
// the caller from the store's Locations index, since models.Artist itself
// carries no location data) and is only consulted when c.Locations is
// non-empty. Every dimension in FilterCriteria is combined with a logical
// AND — an artist must pass all set constraints, not just one.
// An artist whose FirstAlbum string cannot be parsed fails any FirstAlbum
// bound rather than panicking or being silently included.
func matchesCriteria(a models.Artist, locations []string, c FilterCriteria) bool {
	if c.CreationDateMin != nil && a.CreationDate < *c.CreationDateMin {
		return false
	}
	if c.CreationDateMax != nil && a.CreationDate > *c.CreationDateMax {
		return false
	}

	if c.FirstAlbumMin != nil || c.FirstAlbumMax != nil {
		year, err := firstAlbumYear(a.FirstAlbum)
		if err != nil {
			return false
		}
		if c.FirstAlbumMin != nil && year < *c.FirstAlbumMin {
			return false
		}
		if c.FirstAlbumMax != nil && year > *c.FirstAlbumMax {
			return false
		}
	}

	memberCount := len(a.Members)
	if c.MembersMin != nil && memberCount < *c.MembersMin {
		return false
	}
	if c.MembersMax != nil && memberCount > *c.MembersMax {
		return false
	}

	if len(c.Locations) > 0 && !anyLocationMatches(locations, c.Locations) {
		return false
	}

	return true
}

// anyLocationMatches reports whether any of an artist's resolved concert
// locations matches any of the wanted location slugs. Matching is substring-based
// (via strings.Contains) rather than exact-equality so that a hierarchical
// wanted value like "washington-usa" also matches a more specific artist
// location such as "seattle-washington-usa", per the location-hierarchy
// behaviour the spec requires.
func anyLocationMatches(artistLocations []string, wanted []string) bool {
	for _, loc := range artistLocations {
		for _, w := range wanted {
			if strings.Contains(loc, w) {
				return true
			}
		}
	}
	return false
}

// firstAlbumYear extracts the release year from an Artist.FirstAlbum string
// formatted as "DD-MM-YYYY" (the format used by the Groupie Trackers API).
// It returns an error if the string does not have exactly three hyphen-separated
// parts or the year part is not numeric.
func firstAlbumYear(firstAlbum string) (int, error) {
	parts := strings.Split(firstAlbum, "-")
	if len(parts) != 3 {
		return 0, fmt.Errorf("firstAlbumYear: %q is not in DD-MM-YYYY format", firstAlbum)
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("firstAlbumYear: %q has a non-numeric year: %w", firstAlbum, err)
	}
	return year, nil
}
