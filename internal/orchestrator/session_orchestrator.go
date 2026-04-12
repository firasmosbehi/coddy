// Package orchestrator manages LLM interactions and tool execution.
package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/firasmosbehi/coddy/internal/config"
	"github.com/firasmosbehi/coddy/internal/llm"
	"github.com/firasmosbehi/coddy/internal/sandbox"
	"github.com/firasmosbehi/coddy/internal/session"
	"github.com/firasmosbehi/coddy/pkg/models"
)

// SessionOrchestrator manages conversations within a session.
type SessionOrchestrator struct {
	config    *config.Config
	llm       *llm.Client
	sessionID string
	store     session.Store
	messages  []models.Message
}

// NewSessionOrchestrator creates a new session-aware orchestrator.
func NewSessionOrchestrator(cfg *config.Config, llmClient *llm.Client, store session.Store, sessionID string) (*SessionOrchestrator, error) {
	// Get session from store
	ctx := context.Background()
	sess, err := store.Get(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	o := &SessionOrchestrator{
		config:    cfg,
		llm:       llmClient,
		sessionID: sessionID,
		store:     store,
		messages: []models.Message{
			{
				Role:    string(models.RoleSystem),
				Content: defaultSystemPrompt(),
			},
		},
	}

	// Add existing messages from session
	o.messages = append(o.messages, sess.Messages...)

	return o, nil
}

// HandleMessage processes a user message and returns the response.
func (o *SessionOrchestrator) HandleMessage(ctx context.Context, userInput string) (string, error) {
	// Add user message
	o.messages = append(o.messages, models.Message{
		Role:    string(models.RoleUser),
		Content: userInput,
	})

	// Get session with sandbox
	sessWithSandbox, err := o.store.GetWithSandbox(ctx, o.sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get session with sandbox: %w", err)
	}

	// Loop until we get a final response or hit iteration limit
	for i := 0; i < o.config.MaxToolIterations; i++ {
		// Call LLM
		response, err := o.llm.Chat(ctx, o.messages, 0.7)
		if err != nil {
			return "", fmt.Errorf("LLM error: %w", err)
		}

		// Add assistant message to history
		o.messages = append(o.messages, response.Message)

		// Check for tool calls
		if !response.HasToolCalls() {
			// Save messages to session
			o.saveMessages(ctx)
			return response.Text, nil
		}

		// Execute tool calls
		for _, tc := range response.ToolCalls {
			result := o.executeTool(ctx, tc, sessWithSandbox)
			o.messages = append(o.messages, models.Message{
				Role:       string(models.RoleTool),
				ToolCallID: tc.ID,
				Content:    result,
			})
		}
	}

	// Save messages before returning error
	o.saveMessages(ctx)

	return "Error: Exceeded maximum tool call iterations.", nil
}

// executeTool executes a single tool call.
func (o *SessionOrchestrator) executeTool(ctx context.Context, tc models.ToolCall, sess *session.SessionWithSandbox) string {
	switch tc.Name {
	case "execute_code":
		return o.executeCode(ctx, tc, sess.Sandbox)
	case "read_file":
		return o.readFile(ctx, tc, sess.Sandbox)
	case "write_file":
		return o.writeFile(ctx, tc, sess.Sandbox)
	case "list_files":
		return o.listFiles(ctx, tc, sess.Sandbox)
	default:
		return fmt.Sprintf("Error: Unknown tool: %s", tc.Name)
	}
}

// executeCode handles the execute_code tool.
func (o *SessionOrchestrator) executeCode(ctx context.Context, tc models.ToolCall, sb sandbox.Sandbox) string {
	var args models.ExecuteCodeArgs
	if err := json.Unmarshal(tc.Arguments, &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	if !args.Language.IsValid() {
		return fmt.Sprintf("Error: Invalid language: %s", args.Language)
	}

	// Create timeout context
	execCtx, cancel := context.WithTimeout(ctx, o.config.SandboxTimeout)
	defer cancel()

	result, err := sb.Execute(execCtx, args.Code, args.Language)
	if err != nil {
		return fmt.Sprintf("Error executing code: %v", err)
	}

	return result.String()
}

// readFile handles the read_file tool.
func (o *SessionOrchestrator) readFile(ctx context.Context, tc models.ToolCall, sb sandbox.Sandbox) string {
	var args struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(tc.Arguments, &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	data, err := sb.DownloadFile(ctx, args.Path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}

	return string(data)
}

// writeFile handles the write_file tool.
func (o *SessionOrchestrator) writeFile(ctx context.Context, tc models.ToolCall, sb sandbox.Sandbox) string {
	var args struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(tc.Arguments, &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	// Write to temp file first, then upload
	tmpFile := fmt.Sprintf("/tmp/coddy_write_%d", time.Now().UnixNano())
	if err := os.WriteFile(tmpFile, []byte(args.Content), 0644); err != nil {
		return fmt.Sprintf("Error writing temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	if err := sb.UploadFile(ctx, tmpFile, args.Path); err != nil {
		return fmt.Sprintf("Error writing file: %v", err)
	}

	return fmt.Sprintf("File written successfully: %s", args.Path)
}

// listFiles handles the list_files tool.
func (o *SessionOrchestrator) listFiles(ctx context.Context, tc models.ToolCall, sb sandbox.Sandbox) string {
	var args struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(tc.Arguments, &args); err != nil {
		// Default to empty path if parsing fails
		args.Path = "."
	}

	if args.Path == "" {
		args.Path = "."
	}

	files, err := sb.ListFiles(ctx, args.Path)
	if err != nil {
		return fmt.Sprintf("Error listing files: %v", err)
	}

	if len(files) == 0 {
		return "No files found."
	}

	result := "Files:\n"
	for _, f := range files {
		result += fmt.Sprintf("  %s\n", f)
	}
	return result
}

// saveMessages persists messages to the session store.
func (o *SessionOrchestrator) saveMessages(ctx context.Context) {
	sess, err := o.store.Get(ctx, o.sessionID)
	if err != nil {
		return
	}

	// Skip system message when saving
	if len(o.messages) > 1 {
		sess.Messages = o.messages[1:]
	} else {
		sess.Messages = nil
	}

	o.store.Update(ctx, sess)
}

// ClearHistory clears the conversation history (except system prompt).
func (o *SessionOrchestrator) ClearHistory() {
	if len(o.messages) > 0 && o.messages[0].Role == string(models.RoleSystem) {
		o.messages = o.messages[:1]
	} else {
		o.messages = nil
	}

	// Save cleared history
	ctx := context.Background()
	o.saveMessages(ctx)
}

// GetHistory returns the conversation history.
func (o *SessionOrchestrator) GetHistory() []models.Message {
	return o.messages
}
