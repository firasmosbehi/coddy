// Package models defines the core data types for Coddy.
package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// ExecutionResult represents the result of code execution in a sandbox.
type ExecutionResult struct {
	Stdout          string   `json:"stdout"`
	Stderr          string   `json:"stderr"`
	ExitCode        int      `json:"exit_code"`
	OutputFiles     []string `json:"output_files"`
	ExecutionTimeMs int64    `json:"execution_time_ms"`
	TimedOut        bool     `json:"timed_out"`
}

// Success returns true if execution completed successfully.
func (r ExecutionResult) Success() bool {
	return r.ExitCode == 0 && !r.TimedOut
}

// String returns a formatted string representation of the result.
func (r ExecutionResult) String() string {
	var parts []string
	if r.Stdout != "" {
		parts = append(parts, fmt.Sprintf("STDOUT:\n%s", r.Stdout))
	}
	if r.Stderr != "" {
		parts = append(parts, fmt.Sprintf("STDERR:\n%s", r.Stderr))
	}
	if r.TimedOut {
		parts = append(parts, "EXECUTION TIMED OUT")
	}
	if len(r.OutputFiles) > 0 {
		parts = append(parts, fmt.Sprintf("Files created: %s", r.OutputFiles))
	}
	parts = append(parts, fmt.Sprintf("Exit code: %d | Time: %dms", r.ExitCode, r.ExecutionTimeMs))

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += "\n"
		}
		result += part
	}
	return result
}

// Language represents supported programming languages.
type Language string

const (
	Python  Language = "python"
	NodeJS  Language = "nodejs"
)

// IsValid checks if the language is supported.
func (l Language) IsValid() bool {
	switch l {
	case Python, NodeJS:
		return true
	}
	return false
}

// ToolCall represents a function call from the LLM.
type ToolCall struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// LLMResponse represents a response from the LLM.
type LLMResponse struct {
	Text       string     `json:"text"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	Message    Message    `json:"message"`
	Usage      TokenUsage `json:"usage"`
}

// HasToolCalls returns true if the response contains tool calls.
func (r LLMResponse) HasToolCalls() bool {
	return len(r.ToolCalls) > 0
}

// TokenUsage tracks token consumption.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Message represents a chat message.
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

// Role represents the role of a message sender.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Session represents a conversation session.
type Session struct {
	ID            string        `json:"id"`
	SandboxType   string        `json:"sandbox_type"`
	Messages      []Message     `json:"messages"`
	CreatedAt     time.Time     `json:"created_at"`
	LastActivity  time.Time     `json:"last_activity"`
}

// NewSession creates a new session with the given ID.
func NewSession(id string, sandboxType string) *Session {
	now := time.Now().UTC()
	return &Session{
		ID:           id,
		SandboxType:  sandboxType,
		Messages:     make([]Message, 0),
		CreatedAt:    now,
		LastActivity: now,
	}
}

// AddMessage adds a message to the session.
func (s *Session) AddMessage(msg Message) {
	s.Messages = append(s.Messages, msg)
	s.Touch()
}

// Touch updates the last activity timestamp.
func (s *Session) Touch() {
	s.LastActivity = time.Now().UTC()
}

// IsExpired checks if the session has expired.
func (s *Session) IsExpired(timeout time.Duration) bool {
	return time.Since(s.LastActivity) > timeout
}

// ToolDefinition represents an OpenAI-compatible tool definition.
type ToolDefinition struct {
	Type     string       `json:"type"`
	Function FunctionDef  `json:"function"`
}

// FunctionDef represents a function definition.
type FunctionDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// ExecuteCodeArgs represents arguments for the execute_code tool.
type ExecuteCodeArgs struct {
	Language Language `json:"language"`
	Code     string   `json:"code"`
}

// NetworkMode represents the network mode for sandboxes.
type NetworkMode string

const (
	NetworkNone   NetworkMode = "none"
	NetworkBridge NetworkMode = "bridge"
)
