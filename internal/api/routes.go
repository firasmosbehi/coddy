// Package api provides HTTP handlers for the REST API.
package api

import (
	"encoding/json"
	"net/http"
)

// Handler holds dependencies for API handlers.
type Handler struct {
	// TODO: Add orchestrator, session manager, etc.
}

// NewHandler creates a new API handler.
func NewHandler() *Handler {
	return &Handler{}
}

// HealthResponse represents a health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health handles health check requests.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{Status: "healthy"})
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// respondJSON sends a JSON response.
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response.
func respondError(w http.ResponseWriter, status int, err string) {
	respondJSON(w, status, ErrorResponse{Error: err})
}
