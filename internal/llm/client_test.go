package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/firasmosbehi/coddy/pkg/models"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:11434/v1", "qwen3-coder", "")

	if client.baseURL != "http://localhost:11434/v1" {
		t.Errorf("expected baseURL to be set, got %s", client.baseURL)
	}

	if client.model != "qwen3-coder" {
		t.Errorf("expected model to be set, got %s", client.model)
	}

	// Empty API key should default to "not-needed"
	if client.apiKey != "not-needed" {
		t.Errorf("expected apiKey to be 'not-needed', got %s", client.apiKey)
	}
}

func TestClient_Chat_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected /chat/completions, got %s", r.URL.Path)
		}

		// Return mock response
		response := ChatResponse{
			ID:      "test-id",
			Object:  "chat.completion",
			Model:   "test-model",
			Choices: []struct {
				Index   int            `json:"index"`
				Message models.Message `json:"message"`
			}{
				{
					Index: 0,
					Message: models.Message{
						Role:    "assistant",
						Content: "Hello, world!",
					},
				},
			},
			Usage: models.TokenUsage{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-model", "test-key")

	messages := []models.Message{
		{Role: "user", Content: "Hello"},
	}

	ctx := context.Background()
	resp, err := client.Chat(ctx, messages, 0.7)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if resp.Text != "Hello, world!" {
		t.Errorf("expected 'Hello, world!', got %s", resp.Text)
	}

	if resp.Usage.TotalTokens != 15 {
		t.Errorf("expected 15 tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestClient_Chat_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-model", "test-key")

	messages := []models.Message{
		{Role: "user", Content: "Hello"},
	}

	ctx := context.Background()
	_, err := client.Chat(ctx, messages, 0.7)
	if err == nil {
		t.Error("expected error for API failure")
	}
}

func TestClient_Chat_WithToolCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ChatResponse{
			ID:     "test-id",
			Object: "chat.completion",
			Choices: []struct {
				Index   int            `json:"index"`
				Message models.Message `json:"message"`
			}{
				{
					Index: 0,
					Message: models.Message{
						Role: "assistant",
						ToolCalls: []models.ToolCall{
							{
								ID:        "call-1",
								Name:      "execute_code",
								Arguments: []byte(`{"language": "python", "code": "print(1+1)"}`),
							},
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-model", "test-key")

	messages := []models.Message{
		{Role: "user", Content: "Calculate 1+1"},
	}

	ctx := context.Background()
	resp, err := client.Chat(ctx, messages, 0.7)
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}

	if !resp.HasToolCalls() {
		t.Error("expected tool calls in response")
	}

	if len(resp.ToolCalls) != 1 {
		t.Errorf("expected 1 tool call, got %d", len(resp.ToolCalls))
	}

	if resp.ToolCalls[0].Name != "execute_code" {
		t.Errorf("expected tool name 'execute_code', got %s", resp.ToolCalls[0].Name)
	}
}

func TestClient_Chat_NoChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ChatResponse{
			ID:      "test-id",
			Choices: []struct {
				Index   int            `json:"index"`
				Message models.Message `json:"message"`
			}{},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-model", "test-key")

	messages := []models.Message{
		{Role: "user", Content: "Hello"},
	}

	ctx := context.Background()
	_, err := client.Chat(ctx, messages, 0.7)
	if err == nil {
		t.Error("expected error for empty choices")
	}
}

func TestClient_SetHTTPClient(t *testing.T) {
	client := NewClient("http://localhost", "model", "key")

	customClient := &http.Client{}
	client.SetHTTPClient(customClient)

	if client.httpClient != customClient {
		t.Error("expected http client to be updated")
	}
}

func TestChatRequest_Marshal(t *testing.T) {
	req := ChatRequest{
		Model:       "test-model",
		Messages:    []models.Message{{Role: "user", Content: "Hello"}},
		Temperature: 0.7,
		MaxTokens:   100,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled ChatRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.Model != req.Model {
		t.Error("model mismatch")
	}

	if unmarshaled.Temperature != req.Temperature {
		t.Error("temperature mismatch")
	}
}
