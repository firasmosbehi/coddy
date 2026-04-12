// Package llm provides LLM client implementations.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/firasmosbehi/coddy/internal/tools"
	"github.com/firasmosbehi/coddy/pkg/models"
)

// Client provides an interface to LLM APIs.
type Client struct {
	baseURL    string
	model      string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new LLM client.
func NewClient(baseURL, model, apiKey string) *Client {
	if apiKey == "" {
		apiKey = "not-needed"
	}

	return &Client{
		baseURL: baseURL,
		model:   model,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChatRequest represents a chat completion request.
type ChatRequest struct {
	Model       string              `json:"model"`
	Messages    []models.Message    `json:"messages"`
	Tools       []models.ToolDefinition `json:"tools,omitempty"`
	Temperature float64             `json:"temperature,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
}

// ChatResponse represents a chat completion response.
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int            `json:"index"`
		Message models.Message `json:"message"`
	} `json:"choices"`
	Usage models.TokenUsage `json:"usage"`
}

// Chat sends a chat completion request to the LLM.
func (c *Client) Chat(ctx context.Context, messages []models.Message, temperature float64) (models.LLMResponse, error) {
	var response models.LLMResponse

	toolDefs := tools.GetDefinitions()

	reqBody := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		Tools:       toolDefs,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return response, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return response, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return response, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return response, fmt.Errorf("no choices in response")
	}

	choice := chatResp.Choices[0]
	msg := choice.Message

	response.Text = msg.Content
	response.Message = msg
	response.Usage = chatResp.Usage

	// Parse tool calls
	if len(msg.ToolCalls) > 0 {
		response.ToolCalls = make([]models.ToolCall, len(msg.ToolCalls))
		for i, tc := range msg.ToolCalls {
			response.ToolCalls[i] = models.ToolCall{
				ID:        tc.ID,
				Name:      tc.Name,
				Arguments: tc.Arguments,
			}
		}
	}

	return response, nil
}

// SetHTTPClient sets a custom HTTP client (useful for testing).
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}
