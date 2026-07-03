package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"

	"groupie-tracker/internal/store"
)

// ArtistHandler handles GET /artist/{id} requests by looking up a single artist
// and rendering the artist detail page. The store and templates are injected at
// construction time so the handler is stateless and safe for concurrent use.
type ArtistHandler struct {
	store        store.Store
	tmpl         *template.Template
	notFoundTmpl *template.Template
}

// NewArtistHandler constructs an ArtistHandler with the given store, page template,
// and not-found template. All templates are parsed once at construction time and
// reused across all requests, avoiding repeated filesystem reads on every page load.
func NewArtistHandler(s store.Store, tmpl *template.Template, notFoundTmpl *template.Template) http.Handler {
	return &ArtistHandler{store: s, tmpl: tmpl, notFoundTmpl: notFoundTmpl}
}

// ServeHTTP extracts the artist ID from the URL path, validates it is a positive
// integer, and retrieves the matching ArtistPageData from the store. It renders
// the result into a buffer before writing to the response so that a template
// execution error does not result in a partially written 200 response.
// Returns a styled 404 page for non-numeric or unknown IDs, 500 if template execution fails.
func (h *ArtistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NotFoundHandler(h.notFoundTmpl)(w, r)
		return
	}

	data, ok := h.store.ArtistPageDataByID(id)
	if !ok {
		NotFoundHandler(h.notFoundTmpl)(w, r)
		return
	}

	var buf bytes.Buffer
	if err := h.tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w) //nolint:errcheck // response write errors are unrecoverable
}
