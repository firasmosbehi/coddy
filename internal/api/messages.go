// Package api provides message/chat history endpoints.
package api

import (
	"net/http"
	"strings"

	"github.com/firasmosbehi/coddy/internal/session"
	"github.com/firasmosbehi/coddy/pkg/models"
)

// GetMessages retrieves chat history for a session.
func (h *Handlers) GetMessages(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Path[len("/sessions/"):]
	sessionID = strings.TrimSuffix(sessionID, "/messages")
	sessionID = strings.TrimSuffix(sessionID, "/")

	if sessionID == "" {
		respondError(w, http.StatusBadRequest, "Session ID required")
		return
	}

	// Get session
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

	// Filter out system messages for the response
	var messages []map[string]interface{}
	for _, msg := range sess.Messages {
		if msg.Role == "system" {
			continue
		}

		m := map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}

		if msg.ToolCallID != "" {
			m["tool_call_id"] = msg.ToolCallID
		}

		if len(msg.ToolCalls) > 0 {
			m["tool_calls"] = msg.ToolCalls
		}

		messages = append(messages, m)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"session_id": sessionID,
		"messages":   messages,
		"count":      len(messages),
	})
}

// ClearMessages clears chat history for a session (keeps system message).
func (h *Handlers) ClearMessages(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Path[len("/sessions/"):]
	sessionID = strings.TrimSuffix(sessionID, "/messages")
	sessionID = strings.TrimSuffix(sessionID, "/")

	if sessionID == "" {
		respondError(w, http.StatusBadRequest, "Session ID required")
		return
	}

	// Get session
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

	// Clear messages (keep only system message)
	var systemMsg *models.Message
	for _, msg := range sess.Messages {
		if msg.Role == "system" {
			systemMsg = &msg
			break
		}
	}

	sess.Messages = nil
	if systemMsg != nil {
		sess.Messages = append(sess.Messages, *systemMsg)
	}

	// Update session
	if err := h.sessions.Store().Update(ctx, sess); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to clear messages: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":    "Messages cleared",
		"session_id": sessionID,
	})
}
