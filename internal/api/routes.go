// Package api provides HTTP handlers for the REST API.
package api

import (
	"net/http"
)

// SetupRoutes configures all HTTP routes.
func SetupRoutes(mux *http.ServeMux, handlers *Handlers) {
	// Health check
	mux.HandleFunc("/health", handlers.Health)

	// Stats
	mux.HandleFunc("/stats", handlers.GetStats)

	// Session management
	mux.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListSessions(w, r)
		case http.MethodPost:
			handlers.CreateSession(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/sessions/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetSession(w, r)
		case http.MethodDelete:
			handlers.DeleteSession(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
