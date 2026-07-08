package geo

import "strings"

// normalizeLocation converts a Groupie Trackers location slug (e.g.
// "san_francisco-usa") into a human-readable address string
// (e.g. "san francisco, usa") that Nominatim's free-text search can resolve.
// Underscores become spaces (multi-word city names), and the city/country
// separator hyphen becomes a comma.
func normalizeLocation(location string) string {
	location = strings.ReplaceAll(location, "_", " ")
	location = strings.ReplaceAll(location, "-", ", ")
	return location
}
