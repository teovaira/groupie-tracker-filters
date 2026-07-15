package handlers

import (
	"encoding/json"
	"groupie-tracker-filters/internal/models"
	"groupie-tracker-filters/internal/store"
	"html/template"
	"net/http"
	"strconv"
)

// FilterHandler handles structured artist filtering requests and holds the
// store dependency needed to query artists, and the bad-request template for
// malformed numeric query parameters. It is registered as a GET-only route
// and responds with JSON, making it suitable for consumption by the frontend
// filter.js script without a full page reload.
type FilterHandler struct {
	Store          store.Store
	BadRequestTmpl *template.Template
}

// Filter handles GET /api/filter requests by parsing the free-text query and
// structured range/checkbox parameters, delegating to the store's
// FilterArtists method, and encoding the result as a JSON array.
//
// All parameters are optional — a request with none returns every artist,
// satisfying the requirement that filters can be cleared to show all results.
// Query parameters:
//
//	q               — free-text search, same semantics as /api/search
//	creation_min/max — CreationDate range (int, inclusive)
//	first_album_min/max — FirstAlbum release year range (int, inclusive)
//	members_min/max — member count range (int, inclusive)
//	locations       — repeated param, e.g. ?locations=a&locations=b
//
// A malformed numeric parameter returns a styled 400 page immediately, before
// the store is queried, so the store never receives an unparseable bound.
// The response Content-Type is set to application/json before encoding.
func (h *FilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	criteria := store.FilterCriteria{
		Locations: query["locations"],
	}

	var ok bool
	if criteria.CreationDateMin, ok = parseOptionalInt(query, "creation_min"); !ok {
		BadRequestHandler(h.BadRequestTmpl)(w, r)
		return
	}
	if criteria.CreationDateMax, ok = parseOptionalInt(query, "creation_max"); !ok {
		BadRequestHandler(h.BadRequestTmpl)(w, r)
		return
	}
	if criteria.FirstAlbumMin, ok = parseOptionalInt(query, "first_album_min"); !ok {
		BadRequestHandler(h.BadRequestTmpl)(w, r)
		return
	}
	if criteria.FirstAlbumMax, ok = parseOptionalInt(query, "first_album_max"); !ok {
		BadRequestHandler(h.BadRequestTmpl)(w, r)
		return
	}
	if criteria.MembersMin, ok = parseOptionalInt(query, "members_min"); !ok {
		BadRequestHandler(h.BadRequestTmpl)(w, r)
		return
	}
	if criteria.MembersMax, ok = parseOptionalInt(query, "members_max"); !ok {
		BadRequestHandler(h.BadRequestTmpl)(w, r)
		return
	}

	result := h.Store.FilterArtists(query.Get("q"), criteria)
	// FilterArtists may return a nil slice when nothing matches, which encodes
	// as JSON null; the API contract guarantees an empty array instead, same
	// as /api/search, so the frontend can always safely call .length on it.
	if result == nil {
		result = []models.Artist{}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// parseOptionalInt reads the named query parameter and parses it as an int.
// It returns (nil, true) if the parameter is absent — meaning "no constraint"
// — and (nil, false) if present but not a valid integer, signalling the
// caller to reject the request with a 400 rather than silently ignoring a
// malformed bound.
func parseOptionalInt(query map[string][]string, name string) (*int, bool) {
	values, present := query[name]
	if !present || len(values) == 0 || values[0] == "" {
		return nil, true
	}
	v, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, false
	}
	return &v, true
}
