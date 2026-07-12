package store

import (
	"fmt"
	"strconv"
	"strings"
)

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
