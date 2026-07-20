// Package handlers implements all HTTP handlers and middleware for the
// groupie-tracker web server. Each handler receives its dependencies —
// store and template — at construction time via injection, keeping
// handlers stateless and independently testable.
package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

// BadRequestHandler returns an http.HandlerFunc that renders the 400.html template
// and writes a 400 Bad Request status. Used when the client sends a request
// with a missing or empty required query parameter, such as ?q= in /api/search.
func BadRequestHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "400.html", nil); err != nil {
			log.Print(err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		buf.WriteTo(w)
	}
}

// NotFoundHandler returns an http.HandlerFunc that renders the 404.html template
// and writes a 404 Not Found status. If the template fails to execute, it falls
// back to a plain-text http.Error to ensure the client always receives a response.
func NotFoundHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "404.html", nil); err != nil {
			log.Print(err)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		buf.WriteTo(w)
	}
}

// StatusInternalServerError returns an http.HandlerFunc that renders the 500.html
// template and writes a 500 Internal Server Error status. It is used both as a
// direct handler and called by RecoveryMiddleware when a panic is caught.
// If the template fails to execute, it falls back to a plain-text http.Error.
func StatusInternalServerError(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "500.html", nil); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		buf.WriteTo(w)
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
