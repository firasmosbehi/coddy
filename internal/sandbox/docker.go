// Package sandbox provides isolated code execution environments.
package sandbox

import (
	"context"
	"fmt"
	"io"

	"github.com/firasmosbehi/coddy/pkg/models"
)

// DockerSandbox is a placeholder for Docker-based sandbox.
// Full implementation will be added in Phase 2.
type DockerSandbox struct {
	config *Config
}

// NewDockerSandbox creates a new Docker-based sandbox.
func NewDockerSandbox(cfg *Config) (*DockerSandbox, error) {
	return nil, fmt.Errorf("docker sandbox not yet implemented")
}

// Execute runs code in a Docker container.
func (s *DockerSandbox) Execute(ctx context.Context, code string, language models.Language) (models.ExecutionResult, error) {
	return models.ExecutionResult{}, fmt.Errorf("not implemented")
}

// UploadFile copies a file into the container.
func (s *DockerSandbox) UploadFile(ctx context.Context, localPath, sandboxPath string) error {
	return fmt.Errorf("not implemented")
}

// DownloadFile retrieves a file from the container.
func (s *DockerSandbox) DownloadFile(ctx context.Context, sandboxPath string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

// ListFiles returns files in the container.
func (s *DockerSandbox) ListFiles(ctx context.Context, path string) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// Reset clears the container state.
func (s *DockerSandbox) Reset(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

// Close stops and removes the container.
func (s *DockerSandbox) Close() error {
	return fmt.Errorf("not implemented")
}

// CopyToContainer copies data to a path in the container.
func CopyToContainer(ctx context.Context, client interface{}, containerID, path string, content io.Reader) error {
	return fmt.Errorf("not implemented")
}
