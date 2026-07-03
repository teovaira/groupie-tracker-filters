// Package handlers implements all HTTP handlers and middleware for the
// groupie-tracker web server. Each handler receives its dependencies —
// store and template — at construction time via injection, keeping
// handlers stateless and independently testable.
package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"groupie-tracker/internal/store"
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

// ServeHTTP retrieves all artists from the store and renders the home template.
// Returns 500 if template execution fails — the page is never partially written.
func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	artists := h.store.AllArtists()

	var buf bytes.Buffer
	if err := h.tmpl.ExecuteTemplate(&buf, "base", artists); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w) //nolint:errcheck // response write errors are unrecoverable
}
