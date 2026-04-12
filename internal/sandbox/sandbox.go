// Package sandbox provides isolated code execution environments.
package sandbox

import (
	"context"
	"fmt"

	"github.com/firasmosbehi/coddy/pkg/models"
)

// Sandbox defines the interface for code execution environments.
type Sandbox interface {
	// Execute runs code and returns the result.
	Execute(ctx context.Context, code string, language models.Language) (models.ExecutionResult, error)

	// UploadFile copies a file into the sandbox filesystem.
	UploadFile(ctx context.Context, localPath, sandboxPath string) error

	// DownloadFile retrieves a file from the sandbox filesystem.
	DownloadFile(ctx context.Context, sandboxPath string) ([]byte, error)

	// ListFiles returns a list of files in the sandbox.
	ListFiles(ctx context.Context, path string) ([]string, error)

	// Reset clears the sandbox state.
	Reset(ctx context.Context) error

	// Close cleans up sandbox resources.
	Close() error
}

// Factory creates sandboxes based on configuration.
func Factory(cfg *Config) (Sandbox, error) {
	switch cfg.Type {
	case "subprocess":
		return NewSubprocessSandbox(cfg)
	case "docker":
		return NewDockerSandbox(cfg)
	default:
		return nil, fmt.Errorf("unsupported sandbox type: %s", cfg.Type)
	}
}

// Config holds sandbox configuration.
type Config struct {
	Type        string
	Image       string
	Timeout     int
	MemoryLimit string
	CPULimit    float64
	Network     models.NetworkMode
	WorkingDir  string
}
