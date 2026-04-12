// Package sandbox provides isolated code execution environments.
package sandbox

import (
	"context"
	"fmt"

	"github.com/firasmosbehi/coddy/pkg/models"
)

// DockerSandbox runs code in isolated Docker containers.
// NOTE: This is a stub implementation. Full implementation requires
// Docker SDK dependencies which have complex version requirements.
// For production use, implement using github.com/docker/docker/client.
type DockerSandbox struct {
	config *Config
}

// NewDockerSandbox creates a new Docker-based sandbox.
func NewDockerSandbox(cfg *Config) (*DockerSandbox, error) {
	return nil, fmt.Errorf("docker sandbox not yet fully implemented. Use SANDBOX_TYPE=subprocess for development")
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
	return nil
}
