// Package config handles application configuration.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/firasmosbehi/coddy/pkg/models"
)

// Config holds all application configuration.
type Config struct {
	// LLM Configuration
	LLMBaseURL string `json:"llm_base_url"`
	LLMModel   string `json:"llm_model"`
	LLMAPIKey  string `json:"llm_api_key"`

	// Sandbox Configuration
	SandboxType    string              `json:"sandbox_type"`
	SandboxImage   string              `json:"sandbox_image"`
	SandboxTimeout time.Duration       `json:"sandbox_timeout"`
	SandboxMemory  string              `json:"sandbox_memory"`
	SandboxCPU     float64             `json:"sandbox_cpu"`
	SandboxNetwork models.NetworkMode  `json:"sandbox_network"`

	// Orchestrator Configuration
	MaxToolIterations int           `json:"max_tool_iterations"`
	SessionTimeout    time.Duration `json:"session_timeout"`

	// Server Configuration
	Host      string `json:"host"`
	Port      int    `json:"port"`
	LogLevel  string `json:"log_level"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		LLMBaseURL: "http://localhost:11434/v1",
		LLMModel:   "qwen3-coder",
		LLMAPIKey:  "",

		SandboxType:    "subprocess",
		SandboxImage:   "coddy-sandbox:latest",
		SandboxTimeout: 30 * time.Second,
		SandboxMemory:  "512m",
		SandboxCPU:     1.0,
		SandboxNetwork: models.NetworkNone,

		MaxToolIterations: 10,
		SessionTimeout:    1 * time.Hour,

		Host:     "0.0.0.0",
		Port:     8000,
		LogLevel: "info",
	}
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// LLM Configuration
	if v := os.Getenv("LLM_BASE_URL"); v != "" {
		cfg.LLMBaseURL = v
	}
	if v := os.Getenv("LLM_MODEL"); v != "" {
		cfg.LLMModel = v
	}
	if v := os.Getenv("LLM_API_KEY"); v != "" {
		cfg.LLMAPIKey = v
	}

	// Sandbox Configuration
	if v := os.Getenv("SANDBOX_TYPE"); v != "" {
		cfg.SandboxType = v
	}
	if v := os.Getenv("SANDBOX_IMAGE"); v != "" {
		cfg.SandboxImage = v
	}
	if v := os.Getenv("SANDBOX_TIMEOUT"); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid SANDBOX_TIMEOUT: %w", err)
		}
		cfg.SandboxTimeout = time.Duration(d) * time.Second
	}
	if v := os.Getenv("SANDBOX_MEMORY_LIMIT"); v != "" {
		cfg.SandboxMemory = v
	}
	if v := os.Getenv("SANDBOX_CPU_LIMIT"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid SANDBOX_CPU_LIMIT: %w", err)
		}
		cfg.SandboxCPU = f
	}
	if v := os.Getenv("SANDBOX_NETWORK"); v != "" {
		cfg.SandboxNetwork = models.NetworkMode(v)
	}

	// Orchestrator Configuration
	if v := os.Getenv("MAX_TOOL_ITERATIONS"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_TOOL_ITERATIONS: %w", err)
		}
		cfg.MaxToolIterations = i
	}
	if v := os.Getenv("SESSION_TIMEOUT"); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid SESSION_TIMEOUT: %w", err)
		}
		cfg.SessionTimeout = time.Duration(d) * time.Second
	}

	// Server Configuration
	if v := os.Getenv("HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT: %w", err)
		}
		cfg.Port = p
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}

	return cfg, cfg.Validate()
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	// Validate sandbox type
	switch c.SandboxType {
	case "subprocess", "docker":
		// Valid
	default:
		return fmt.Errorf("invalid sandbox type: %s (must be 'subprocess' or 'docker')", c.SandboxType)
	}

	// Validate network mode
	switch c.SandboxNetwork {
	case models.NetworkNone, models.NetworkBridge:
		// Valid
	default:
		return fmt.Errorf("invalid network mode: %s", c.SandboxNetwork)
	}

	// Validate timeouts
	if c.SandboxTimeout < 1*time.Second || c.SandboxTimeout > 120*time.Second {
		return fmt.Errorf("sandbox timeout must be between 1 and 120 seconds")
	}

	if c.MaxToolIterations < 1 || c.MaxToolIterations > 50 {
		return fmt.Errorf("max tool iterations must be between 1 and 50")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	return nil
}

// IsLocalLLM returns true if using a local LLM.
func (c *Config) IsLocalLLM() bool {
	return c.LLMAPIKey == "" || strings.Contains(c.LLMBaseURL, "localhost")
}

// Address returns the server address.
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
