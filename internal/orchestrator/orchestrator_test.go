package orchestrator

import (
	"encoding/json"
	"testing"

	"github.com/firasmosbehi/coddy/pkg/models"
)

func TestDefaultSystemPrompt(t *testing.T) {
	prompt := defaultSystemPrompt()

	if prompt == "" {
		t.Error("expected non-empty system prompt")
	}

	if len(prompt) < 100 {
		t.Error("expected prompt to be substantial")
	}
}

func TestToolCallParsing(t *testing.T) {
	// Test that tool calls are properly parsed
	tc := models.ToolCall{
		ID:        "call-1",
		Name:      "execute_code",
		Arguments: []byte(`{"language": "python", "code": "print(1+1)"}`),
	}

	if tc.ID != "call-1" {
		t.Errorf("expected ID 'call-1', got %s", tc.ID)
	}

	if tc.Name != "execute_code" {
		t.Errorf("expected Name 'execute_code', got %s", tc.Name)
	}

	// Verify arguments can be parsed
	var args map[string]interface{}
	if err := json.Unmarshal(tc.Arguments, &args); err != nil {
		t.Errorf("failed to unmarshal arguments: %v", err)
	}

	if args["language"] != "python" {
		t.Errorf("expected language 'python', got %v", args["language"])
	}
}

func TestLLMResponse_HasToolCalls(t *testing.T) {
	// No tool calls
	resp1 := models.LLMResponse{
		Text: "Hello!",
	}
	if resp1.HasToolCalls() {
		t.Error("expected HasToolCalls to be false")
	}

	// With tool calls
	resp2 := models.LLMResponse{
		ToolCalls: []models.ToolCall{
			{ID: "1", Name: "test"},
		},
	}
	if !resp2.HasToolCalls() {
		t.Error("expected HasToolCalls to be true")
	}
}

func TestExecutionResult_String(t *testing.T) {
	result := models.ExecutionResult{
		Stdout:          "Hello",
		Stderr:          "Warning",
		ExitCode:        0,
		OutputFiles:     []string{"file.txt"},
		ExecutionTimeMs: 100,
		TimedOut:        false,
	}

	str := result.String()
	if str == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestExecutionResult_Success(t *testing.T) {
	tests := []struct {
		name     string
		result   models.ExecutionResult
		expected bool
	}{
		{"success", models.ExecutionResult{ExitCode: 0, TimedOut: false}, true},
		{"failure", models.ExecutionResult{ExitCode: 1, TimedOut: false}, false},
		{"timeout", models.ExecutionResult{ExitCode: 0, TimedOut: true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.Success(); got != tt.expected {
				t.Errorf("Success() = %v, want %v", got, tt.expected)
			}
		})
	}
}
