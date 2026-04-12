package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.LLMBaseURL != "http://localhost:11434/v1" {
		t.Errorf("expected default LLMBaseURL, got %s", cfg.LLMBaseURL)
	}

	if cfg.LLMModel != "qwen3-coder" {
		t.Errorf("expected default LLMModel, got %s", cfg.LLMModel)
	}

	if cfg.SandboxType != "subprocess" {
		t.Errorf("expected default SandboxType, got %s", cfg.SandboxType)
	}

	if cfg.SandboxTimeout != 30*time.Second {
		t.Errorf("expected default SandboxTimeout, got %v", cfg.SandboxTimeout)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid sandbox type",
			cfg: func() *Config {
				c := DefaultConfig()
				c.SandboxType = "invalid"
				return c
			}(),
			wantErr: true,
		},
		{
			name: "sandbox timeout too high",
			cfg: func() *Config {
				c := DefaultConfig()
				c.SandboxTimeout = 200 * time.Second
				return c
			}(),
			wantErr: true,
		},
		{
			name: "invalid port",
			cfg: func() *Config {
				c := DefaultConfig()
				c.Port = 0
				return c
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_IsLocalLLM(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		expected bool
	}{
		{
			name:     "localhost URL without API key",
			cfg:      &Config{LLMBaseURL: "http://localhost:11434/v1"},
			expected: true,
		},
		{
			name:     "remote URL with API key",
			cfg:      &Config{LLMBaseURL: "https://api.openai.com/v1", LLMAPIKey: "sk-test"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.IsLocalLLM(); got != tt.expected {
				t.Errorf("IsLocalLLM() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Set environment variables
	os.Setenv("LLM_MODEL", "test-model")
	os.Setenv("PORT", "9000")
	defer func() {
		os.Unsetenv("LLM_MODEL")
		os.Unsetenv("PORT")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.LLMModel != "test-model" {
		t.Errorf("expected LLMModel to be loaded from env, got %s", cfg.LLMModel)
	}

	if cfg.Port != 9000 {
		t.Errorf("expected Port to be loaded from env, got %d", cfg.Port)
	}
}
