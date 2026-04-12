// Package sandbox provides isolated code execution environments.
package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/firasmosbehi/coddy/pkg/models"
)

// SubprocessSandbox runs code using local subprocesses.
// WARNING: This provides NO security isolation and should only be used
// for local development with trusted code.
type SubprocessSandbox struct {
	config     *Config
	workingDir string
}

// NewSubprocessSandbox creates a new subprocess-based sandbox.
func NewSubprocessSandbox(cfg *Config) (*SubprocessSandbox, error) {
	workingDir := cfg.WorkingDir
	if workingDir == "" {
		var err error
		workingDir, err = os.MkdirTemp("", "coddy_sandbox_")
		if err != nil {
			return nil, fmt.Errorf("failed to create working directory: %w", err)
		}
	}

	return &SubprocessSandbox{
		config:     cfg,
		workingDir: workingDir,
	}, nil
}

// Execute runs code in a subprocess.
func (s *SubprocessSandbox) Execute(ctx context.Context, code string, language models.Language) (models.ExecutionResult, error) {
	start := time.Now()
	result := models.ExecutionResult{
		ExitCode: -1,
	}

	// Track files before execution
	filesBefore, err := s.listFiles()
	if err != nil {
		return result, fmt.Errorf("failed to list files: %w", err)
	}

	// Prepare command based on language
	var cmd *exec.Cmd
	switch language {
	case models.Python:
		cmd = exec.CommandContext(ctx, "python3", "-c", code)
	case models.NodeJS:
		cmd = exec.CommandContext(ctx, "node", "-e", code)
	default:
		return result, fmt.Errorf("unsupported language: %s", language)
	}

	cmd.Dir = s.workingDir

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()
	result.ExecutionTimeMs = time.Since(start).Milliseconds()

	// Handle timeout
	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.Stderr = "Execution timed out"
		return result, nil
	}

	// Get exit code
	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	// Detect new files
	filesAfter, err := s.listFiles()
	if err == nil {
		result.OutputFiles = diffFiles(filesBefore, filesAfter)
	}

	return result, nil
}

// UploadFile copies a file into the sandbox.
func (s *SubprocessSandbox) UploadFile(ctx context.Context, localPath, sandboxPath string) error {
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %w", err)
	}

	sandboxPath = filepath.Join(s.workingDir, sandboxPath)
	if err := os.MkdirAll(filepath.Dir(sandboxPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(sandboxPath, data, 0644)
}

// DownloadFile retrieves a file from the sandbox.
func (s *SubprocessSandbox) DownloadFile(ctx context.Context, sandboxPath string) ([]byte, error) {
	sandboxPath = filepath.Join(s.workingDir, sandboxPath)
	return os.ReadFile(sandboxPath)
}

// ListFiles returns files in the sandbox.
func (s *SubprocessSandbox) ListFiles(ctx context.Context, path string) ([]string, error) {
	return s.listFiles()
}

// Reset clears the sandbox working directory.
func (s *SubprocessSandbox) Reset(ctx context.Context) error {
	entries, err := os.ReadDir(s.workingDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(s.workingDir, entry.Name())
		if entry.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}

// Close cleans up the sandbox.
func (s *SubprocessSandbox) Close() error {
	return os.RemoveAll(s.workingDir)
}

func (s *SubprocessSandbox) listFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(s.workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(s.workingDir, path)
			if err != nil {
				return err
			}
			files = append(files, rel)
		}
		return nil
	})

	return files, err
}

func diffFiles(before, after []string) []string {
	beforeSet := make(map[string]bool)
	for _, f := range before {
		beforeSet[f] = true
	}

	var newFiles []string
	for _, f := range after {
		if !beforeSet[f] {
			newFiles = append(newFiles, f)
		}
	}
	return newFiles
}

// ReadFileContent reads a file and returns its content as string.
func ReadFileContent(r io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	return buf.String(), err
}
