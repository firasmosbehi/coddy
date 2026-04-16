// Package api provides WebSocket support for real-time streaming.
package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/firasmosbehi/coddy/internal/orchestrator"
	"github.com/firasmosbehi/coddy/pkg/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, configure this properly
		return true
	},
}

// WSMessage represents a WebSocket message.
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// WSChatRequest represents a chat request over WebSocket.
type WSChatRequest struct {
	Message string `json:"message"`
}

// WSChatResponse represents a chat response over WebSocket.
type WSChatResponse struct {
	Type      string `json:"type"` // "chunk", "tool_call", "tool_result", "done", "error"
	Content   string `json:"content,omitempty"`
	ToolCall  *models.ToolCall `json:"tool_call,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	Error     string `json:"error,omitempty"`
}

// HandleWebSocket handles WebSocket connections for streaming chat.
func (h *Handlers) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from path
	sessionID := r.URL.Path[len("/ws/sessions/"):]
	if sessionID == "" || sessionID == "/" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	// Verify session exists
	ctx := r.Context()
	_, err := h.sessions.Store().Get(ctx, sessionID)
	if err != nil {
		if err.Error() == "session not found" {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connection established for session %s", sessionID)

	// Handle messages
	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "chat":
			var chatReq WSChatRequest
			if err := json.Unmarshal(msg.Payload, &chatReq); err != nil {
				sendWSError(conn, "Invalid chat request")
				continue
			}
			h.handleWSChat(conn, sessionID, chatReq)
		default:
			sendWSError(conn, "Unknown message type: "+msg.Type)
		}
	}

	log.Printf("WebSocket connection closed for session %s", sessionID)
}

func (h *Handlers) handleWSChat(conn *websocket.Conn, sessionID string, req WSChatRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create session orchestrator
	orch, err := orchestrator.NewSessionOrchestrator(h.config, h.llm, h.sessions.Store(), sessionID)
	if err != nil {
		sendWSError(conn, "Failed to create orchestrator: "+err.Error())
		return
	}

	// Handle with streaming simulation
	// Note: For true streaming, we'd need to modify the LLM client and orchestrator
	// to support streaming responses. For now, we send the full response.
	response, err := orch.HandleMessage(ctx, req.Message)
	if err != nil {
		sendWSError(conn, err.Error())
		return
	}

	// Send response in chunks for better UX
	chunkSize := 50
	for i := 0; i < len(response); i += chunkSize {
		end := i + chunkSize
		if end > len(response) {
			end = len(response)
		}

		chunk := WSChatResponse{
			Type:      "chunk",
			Content:   response[i:end],
			SessionID: sessionID,
		}

		if err := conn.WriteJSON(chunk); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}

		// Small delay for streaming effect
		time.Sleep(10 * time.Millisecond)
	}

	// Send done message
	done := WSChatResponse{
		Type:      "done",
		SessionID: sessionID,
	}
	conn.WriteJSON(done)
}

func sendWSError(conn *websocket.Conn, errMsg string) {
	resp := WSChatResponse{
		Type:  "error",
		Error: errMsg,
	}
	conn.WriteJSON(resp)
}
