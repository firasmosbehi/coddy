// Package orchestrator manages the LLM conversation loop and tool execution.
package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/firasmosbehi/coddy/internal/config"
	"github.com/firasmosbehi/coddy/internal/llm"
	"github.com/firasmosbehi/coddy/internal/sandbox"
	"github.com/firasmosbehi/coddy/pkg/models"
)

// Orchestrator manages the conversation flow between user, LLM, and tools.
type Orchestrator struct {
	config   *config.Config
	llm      *llm.Client
	sandbox  sandbox.Sandbox
	messages []models.Message
}

// New creates a new orchestrator.
func New(cfg *config.Config, llmClient *llm.Client, sb sandbox.Sandbox) *Orchestrator {
	o := &Orchestrator{
		config:  cfg,
		llm:     llmClient,
		sandbox: sb,
		messages: []models.Message{
			{
				Role:    string(models.RoleSystem),
				Content: defaultSystemPrompt(),
			},
		},
	}
	return o
}

// HandleMessage processes a user message and returns the response.
func (o *Orchestrator) HandleMessage(ctx context.Context, userInput string) (string, error) {
	// Add user message
	o.messages = append(o.messages, models.Message{
		Role:    string(models.RoleUser),
		Content: userInput,
	})

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
			return response.Text, nil
		}

		// Execute tool calls
		for _, tc := range response.ToolCalls {
			result := o.executeTool(ctx, tc)
			o.messages = append(o.messages, models.Message{
				Role:       string(models.RoleTool),
				ToolCallID: tc.ID,
				Content:    result,
			})
		}
	}

	return "Error: Exceeded maximum tool call iterations.", nil
}

// executeTool executes a single tool call.
func (o *Orchestrator) executeTool(ctx context.Context, tc models.ToolCall) string {
	switch tc.Name {
	case "execute_code":
		return o.executeCode(ctx, tc)
	case "read_file":
		return o.readFile(ctx, tc)
	case "write_file":
		return o.writeFile(ctx, tc)
	case "list_files":
		return o.listFiles(ctx, tc)
	default:
		return fmt.Sprintf("Error: Unknown tool: %s", tc.Name)
	}
}

// executeCode handles the execute_code tool.
func (o *Orchestrator) executeCode(ctx context.Context, tc models.ToolCall) string {
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

	result, err := o.sandbox.Execute(execCtx, args.Code, args.Language)
	if err != nil {
		return fmt.Sprintf("Error executing code: %v", err)
	}

	return result.String()
}

// readFile handles the read_file tool.
func (o *Orchestrator) readFile(ctx context.Context, tc models.ToolCall) string {
	var args struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(tc.Arguments, &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	data, err := o.sandbox.DownloadFile(ctx, args.Path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}

	return string(data)
}

// writeFile handles the write_file tool.
func (o *Orchestrator) writeFile(ctx context.Context, tc models.ToolCall) string {
	var args struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(tc.Arguments, &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}

	// Write to temp file first, then upload
	tmpFile := fmt.Sprintf("/tmp/coddy_write_%d", time.Now().UnixNano())
	if err := o.sandbox.UploadFile(ctx, tmpFile, args.Path); err != nil {
		return fmt.Sprintf("Error writing file: %v", err)
	}

	return fmt.Sprintf("File written successfully: %s", args.Path)
}

// listFiles handles the list_files tool.
func (o *Orchestrator) listFiles(ctx context.Context, tc models.ToolCall) string {
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

	files, err := o.sandbox.ListFiles(ctx, args.Path)
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

// ClearHistory clears the conversation history (except system prompt).
func (o *Orchestrator) ClearHistory() {
	if len(o.messages) > 0 && o.messages[0].Role == string(models.RoleSystem) {
		o.messages = o.messages[:1]
	} else {
		o.messages = nil
	}
}

// GetHistory returns the conversation history.
func (o *Orchestrator) GetHistory() []models.Message {
	return o.messages
}

// defaultSystemPrompt returns the default system prompt.
func defaultSystemPrompt() string {
	return `You are Coddy, an AI assistant with access to a code execution environment.

You can run Python 3.12 and Node.js 20 code to help users with calculations, data analysis, file processing, and more.

## When to Use Code

You SHOULD use code for:
- Mathematical calculations and computations
- Data analysis and visualization
- File processing (CSV, JSON, images, etc.)
- Text processing and transformation
- Any task that requires precise computation

You should NOT use code for:
- Simple factual questions
- Opinions or recommendations
- General conversation

## How to Write Code

### Python
- Use print() to output results
- The environment includes: numpy, pandas, matplotlib, scipy, scikit-learn

### Node.js
- Use console.log() to output results

## Important Notes

1. **Stateful Environment**: Variables persist between code executions.
2. **File Paths**: Use paths relative to the working directory.
3. **Error Handling**: If code fails, read the error and retry with fixes.
4. **Timeouts**: Code execution has a time limit (30 seconds default).
5. **Output**: Always use print/console.log to show results.

After executing code, explain what it did in natural language.`
}
