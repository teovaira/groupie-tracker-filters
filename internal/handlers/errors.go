// Package handlers implements all HTTP handlers and middleware for the
// groupie-tracker web server. Each handler receives its dependencies —
// store and template — at construction time via injection, keeping
// handlers stateless and independently testable.
package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// BadRequestHandler returns an http.HandlerFunc that renders the 400.html template
// and writes a 400 Bad Request status. Used when the client sends a request
// with a missing or empty required query parameter, such as ?q= in /api/search.
func BadRequestHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		err := tmpl.ExecuteTemplate(w, "400.html", nil)
		if err != nil {
			log.Print(err)
			http.Error(w, "bad request", http.StatusBadRequest)
		}
	}
}

// NotFoundHandler returns an http.HandlerFunc that renders the 404.html template
// and writes a 404 Not Found status. If the template fails to execute, it falls
// back to a plain-text http.Error to ensure the client always receives a response.
func NotFoundHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		err := tmpl.ExecuteTemplate(w, "404.html", nil)
		if err != nil {
			log.Print(err)
			http.Error(w, "not found", http.StatusNotFound)
		}
	}
}

// StatusInternalServerError returns an http.HandlerFunc that renders the 500.html
// template and writes a 500 Internal Server Error status. It is used both as a
// direct handler and called by RecoveryMiddleware when a panic is caught.
// If the template fails to execute, it falls back to a plain-text http.Error.
func StatusInternalServerError(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		err := tmpl.ExecuteTemplate(w, "500.html", nil)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

// RecoveryMiddleware wraps the provided handler with panic recovery logic.
// If any handler in the chain panics, the deferred recover catches it,
// logs the panic value, and delegates to StatusInternalServerError to send
// a 500 response. This prevents a single unhandled panic from crashing
// the entire server process.
func RecoveryMiddleware(tmpl *template.Template, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic recovered", err)
				StatusInternalServerError(tmpl)(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
