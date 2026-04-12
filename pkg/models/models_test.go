package models

import (
	"testing"
	"time"
)

func TestExecutionResult_Success(t *testing.T) {
	tests := []struct {
		name     string
		result   ExecutionResult
		expected bool
	}{
		{
			name:     "successful execution",
			result:   ExecutionResult{ExitCode: 0, TimedOut: false},
			expected: true,
		},
		{
			name:     "failed execution",
			result:   ExecutionResult{ExitCode: 1, TimedOut: false},
			expected: false,
		},
		{
			name:     "timed out execution",
			result:   ExecutionResult{ExitCode: 0, TimedOut: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.Success(); got != tt.expected {
				t.Errorf("Success() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLanguage_IsValid(t *testing.T) {
	tests := []struct {
		lang     Language
		expected bool
	}{
		{Python, true},
		{NodeJS, true},
		{Language("invalid"), false},
		{Language(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.lang), func(t *testing.T) {
			if got := tt.lang.IsValid(); got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSession_AddMessage(t *testing.T) {
	session := NewSession("test-id", "subprocess")

	initialActivity := session.LastActivity
	time.Sleep(10 * time.Millisecond)

	session.AddMessage(Message{Role: "user", Content: "Hello"})

	if len(session.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(session.Messages))
	}

	if session.Messages[0].Content != "Hello" {
		t.Errorf("expected 'Hello', got %s", session.Messages[0].Content)
	}

	if !session.LastActivity.After(initialActivity) {
		t.Error("expected LastActivity to be updated")
	}
}

func TestSession_IsExpired(t *testing.T) {
	session := NewSession("test-id", "subprocess")

	// Session should not be expired immediately
	if session.IsExpired(1 * time.Hour) {
		t.Error("new session should not be expired")
	}

	// Manually set last activity to be old
	session.LastActivity = time.Now().UTC().Add(-2 * time.Hour)

	if !session.IsExpired(1 * time.Hour) {
		t.Error("session should be expired after timeout")
	}
}

func TestLLMResponse_HasToolCalls(t *testing.T) {
	tests := []struct {
		name     string
		response LLMResponse
		expected bool
	}{
		{
			name:     "no tool calls",
			response: LLMResponse{Text: "Hello"},
			expected: false,
		},
		{
			name:     "with tool calls",
			response: LLMResponse{ToolCalls: []ToolCall{{ID: "1", Name: "test"}}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.HasToolCalls(); got != tt.expected {
				t.Errorf("HasToolCalls() = %v, want %v", got, tt.expected)
			}
		})
	}
}
