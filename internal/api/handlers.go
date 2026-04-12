// Package api provides HTTP handlers for the REST API.
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/firasmosbehi/coddy/internal/config"
	"github.com/firasmosbehi/coddy/internal/llm"
	"github.com/firasmosbehi/coddy/internal/session"
)

// Handlers holds all HTTP handlers.
type Handlers struct {
	config   *config.Config
	llm      *llm.Client
	sessions *session.SessionManager
}

// NewHandlers creates a new handlers instance.
func NewHandlers(cfg *config.Config, llmClient *llm.Client, sessionManager *session.SessionManager) *Handlers {
	return &Handlers{
		config:   cfg,
		llm:      llmClient,
		sessions: sessionManager,
	}
}

// CreateSessionRequest represents a request to create a session.
type CreateSessionRequest struct {
	SandboxType string `json:"sandbox_type,omitempty"`
}

// CreateSessionResponse represents the response from creating a session.
type CreateSessionResponse struct {
	ID          string    `json:"id"`
	SandboxType string    `json:"sandbox_type"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateSession handles session creation.
func (h *Handlers) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Use default sandbox type
		req.SandboxType = h.config.SandboxType
	}

	if req.SandboxType == "" {
		req.SandboxType = "subprocess"
	}

	ctx := r.Context()
	sess, err := h.sessions.Store().Create(ctx, req.SandboxType)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create session: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, CreateSessionResponse{
		ID:          sess.ID,
		SandboxType: sess.SandboxType,
		CreatedAt:   sess.CreatedAt,
	})
}

// GetSession retrieves a session.
func (h *Handlers) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Path[len("/sessions/"):]
	if sessionID == "" || sessionID == "/" {
		respondError(w, http.StatusBadRequest, "Session ID required")
		return
	}

	ctx := r.Context()
	sess, err := h.sessions.Store().Get(ctx, sessionID)
	if err != nil {
		if err == session.ErrSessionNotFound {
			respondError(w, http.StatusNotFound, "Session not found")
			return
		}
		if err == session.ErrSessionExpired {
			respondError(w, http.StatusGone, "Session expired")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, sess)
}

// DeleteSession deletes a session.
func (h *Handlers) DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Path[len("/sessions/"):]
	if sessionID == "" || sessionID == "/" {
		respondError(w, http.StatusBadRequest, "Session ID required")
		return
	}

	ctx := r.Context()
	if err := h.sessions.Store().Delete(ctx, sessionID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListSessions lists all active sessions.
func (h *Handlers) ListSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ids, err := h.sessions.Store().List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"sessions": ids,
		"count":    len(ids),
	})
}

// GetStats returns server statistics.
func (h *Handlers) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := h.sessions.Stats()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"sessions":  stats,
		"timestamp": time.Now().UTC(),
	})
}

// Health handles health check requests.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// respondJSON sends a JSON response.
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response.
func respondError(w http.ResponseWriter, status int, err string) {
	respondJSON(w, status, map[string]string{"error": err})
}
