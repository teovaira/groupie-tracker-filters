// Package handlers implements all HTTP handlers and middleware for the
// groupie-tracker web server. Each handler receives its dependencies —
// store and template — at construction time via injection, keeping
// handlers stateless and independently testable.
package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"groupie-tracker-filters/internal/models"
	"groupie-tracker-filters/internal/store"
)

// HomeHandler handles GET / requests by rendering the artist list page.
type HomeHandler struct {
	store store.Store
	tmpl  *template.Template
}

// NewHomeHandler constructs a HomeHandler with the given store and template.
// The template is parsed once at construction time and reused across requests.
func NewHomeHandler(s store.Store, tmpl *template.Template) http.Handler {
	return &HomeHandler{store: s, tmpl: tmpl}
}

// ServeHTTP retrieves all artists, the location filter vocabulary, and the
// range-slider bounds from the store and renders the home template with them.
// LocationGroups and Bounds are always the complete, unfiltered values — they
// seed the filter panel's checkboxes and sliders and do not depend on which
// artists are currently displayed.
// Returns 500 if template execution fails — the page is never partially written.
func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := models.HomePageData{
		Artists:        h.store.AllArtists(),
		LocationGroups: h.store.LocationGroups(),
		Bounds:         h.store.FilterBounds(),
	}

	var buf bytes.Buffer
	if err := h.tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w) //nolint:errcheck // response write errors are unrecoverable
}
