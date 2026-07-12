package store

import (
	"fmt"
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
