// Package main is the entry point for the groupie-tracker web server.
// It loads all artist data from the external API on startup, builds the
// in-memory store, registers all HTTP routes, and starts listening.
package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/geo"
	"groupie-tracker/internal/handlers"
	"groupie-tracker/internal/models"
	"groupie-tracker/internal/store"
)

// addr is the TCP address the HTTP server listens on.
const (
	addr        = ":8080"
	geocoderURL = "https://nominatim.openstreetmap.org/search"
)

// main loads all data from the external API, wires up the dependency graph
// (store, templates, handlers, middleware), and starts the HTTP server.
func main() {
	// Load all data from the external API once at startup.
	if err := api.LoadData(); err != nil {
		log.Fatalf("failed to load data: %v", err)
	}

	d := api.GetData()
	geocoder := geo.NewRealGeocoder(geocoderURL)
	markersByArtistID := make(map[int][]models.Marker)
	for _, entry := range d.Locations.Index {
		artistID := entry.ID
		var markers []models.Marker
		for _, location := range entry.Locations {
			coords, err := geocoder.Geocode(location)
			if err != nil {
				log.Printf("failed to geocode %q: %v", location, err)
				continue
			}
			marker := models.Marker{
				Name: location,
				Lat:  coords.Latitude,
				Lng:  coords.Longitude,
			}
			markers = append(markers, marker)
		}
		markersByArtistID[artistID] = markers
	}
	s := &store.RealStore{
		Artists:   d.Artists,
		Locations: d.Locations,
		Dates:     d.Dates,
		Relations: d.Relations,
		Markers:   markersByArtistID,
	}
	log.Printf("Data loaded: %d artists", len(d.Artists))
	homeTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/home.html",
	))

	artistTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/artist.html",
	))
	notFoundTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/404.html",
	))
	serverErrorTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/500.html",
	))
	badRequestTmpl := template.Must(template.ParseFiles(
		"web/templates/base.html",
		"web/templates/400.html",
	))

	homeHandler := handlers.NewHomeHandler(s, homeTmpl)
	notFoundHandler := handlers.NotFoundHandler(notFoundTmpl)
	artistHandler := handlers.NewArtistHandler(s, artistTmpl, notFoundTmpl)
	searchHandler := &handlers.SearchHandler{Store: s, BadRequestTmpl: badRequestTmpl}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFoundHandler(w, r)
			return
		}
		homeHandler.ServeHTTP(w, r)
	})
	mux.Handle("GET /artist/{id}", artistHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.HandleFunc("GET /api/search", searchHandler.Search)

	// mux.HandleFunc("GET /panic-test", func(w http.ResponseWriter, r *http.Request) {
	// 	panic("intentional test panic")
	// })

	log.Printf("server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handlers.RecoveryMiddleware(serverErrorTmpl, mux)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
