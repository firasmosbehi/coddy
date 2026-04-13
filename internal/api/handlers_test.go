package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/firasmosbehi/coddy/internal/config"
	"github.com/firasmosbehi/coddy/internal/llm"
	"github.com/firasmosbehi/coddy/internal/sandbox"
	"github.com/firasmosbehi/coddy/internal/session"
)

func setupTestHandlers(t *testing.T) (*Handlers, func()) {
	cfg := &config.Config{
		SandboxType:    "subprocess",
		SessionTimeout: 1 * time.Hour,
	}

	llmClient := llm.NewClient("http://localhost", "test-model", "")

	sessionConfig := &session.StoreConfig{
		SessionTimeout: cfg.SessionTimeout,
		SandboxConfig: &sandbox.Config{
			Type: "subprocess",
		},
	}

	sessionManager := session.NewSessionManager(sessionConfig)
	sessionManager.Start()

	handlers := NewHandlers(cfg, llmClient, sessionManager)

	cleanup := func() {
		sessionManager.Stop()
	}

	return handlers, cleanup
}

func TestHandlers_Health(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handlers.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %s", response["status"])
	}
}

func TestHandlers_CreateSession(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	reqBody := CreateSessionRequest{SandboxType: "subprocess"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/sessions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handlers.CreateSession(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response CreateSessionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.ID == "" {
		t.Error("expected session ID to be set")
	}

	if response.SandboxType != "subprocess" {
		t.Errorf("expected sandbox type 'subprocess', got %s", response.SandboxType)
	}
}

func TestHandlers_CreateSession_DefaultType(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/sessions", bytes.NewReader([]byte("{}")))
	w := httptest.NewRecorder()

	handlers.CreateSession(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandlers_GetSession(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	// First create a session
	ctx := req.Context()
	sess, err := handlers.sessions.Store().Create(ctx, "subprocess")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Now get it
	req := httptest.NewRequest(http.MethodGet, "/sessions/"+sess.ID, nil)
	w := httptest.NewRecorder()

	handlers.GetSession(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandlers_GetSession_NotFound(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/sessions/non-existent", nil)
	w := httptest.NewRecorder()

	handlers.GetSession(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandlers_GetSession_MissingID(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/sessions/", nil)
	w := httptest.NewRecorder()

	handlers.GetSession(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlers_DeleteSession(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	// Create a session first
	ctx := req.Context()
	sess, _ := handlers.sessions.Store().Create(ctx, "subprocess")

	req := httptest.NewRequest(http.MethodDelete, "/sessions/"+sess.ID, nil)
	w := httptest.NewRecorder()

	handlers.DeleteSession(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify it's deleted
	_, err := handlers.sessions.Store().Get(ctx, sess.ID)
	if err != session.ErrSessionNotFound {
		t.Error("expected session to be deleted")
	}
}

func TestHandlers_DeleteSession_MissingID(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodDelete, "/sessions/", nil)
	w := httptest.NewRecorder()

	handlers.DeleteSession(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlers_ListSessions(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	// Create a few sessions
	ctx := req.Context()
	handlers.sessions.Store().Create(ctx, "subprocess")
	handlers.sessions.Store().Create(ctx, "subprocess")

	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	w := httptest.NewRecorder()

	handlers.ListSessions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	count, ok := response["count"].(float64)
	if !ok {
		t.Fatal("expected count in response")
	}

	if count != 2 {
		t.Errorf("expected 2 sessions, got %v", count)
	}
}

func TestHandlers_GetStats(t *testing.T) {
	handlers, cleanup := setupTestHandlers(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()

	handlers.GetStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if _, ok := response["sessions"]; !ok {
		t.Error("expected sessions in response")
	}
}

// Helper to get context from request
var req = httptest.NewRequest(http.MethodGet, "/", nil)
