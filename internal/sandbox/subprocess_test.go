package sandbox

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/firasmosbehi/coddy/pkg/models"
)

func TestSubprocessSandbox_Execute(t *testing.T) {
	cfg := &Config{
		Type:    "subprocess",
		Timeout: 30,
	}

	sb, err := NewSubprocessSandbox(cfg)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}
	defer sb.Close()

	ctx := context.Background()

	tests := []struct {
		name       string
		code       string
		language   models.Language
		wantStdout string
		wantErr    bool
	}{
		{
			name:       "python hello world",
			code:       "print('Hello, World!')",
			language:   models.Python,
			wantStdout: "Hello, World!\n",
			wantErr:    false,
		},
		{
			name:       "python calculation",
			code:       "print(2**10)",
			language:   models.Python,
			wantStdout: "1024\n",
			wantErr:    false,
		},
		{
			name:       "nodejs hello world",
			code:       "console.log('Hello from Node!')",
			language:   models.NodeJS,
			wantStdout: "Hello from Node!\n",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sb.Execute(ctx, tt.code, tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result.Stdout != tt.wantStdout {
				t.Errorf("Execute() stdout = %q, want %q", result.Stdout, tt.wantStdout)
			}
			if result.ExitCode != 0 {
				t.Errorf("Execute() exit code = %d, want 0", result.ExitCode)
			}
		})
	}
}

func TestSubprocessSandbox_ExecuteTimeout(t *testing.T) {
	cfg := &Config{
		Type:    "subprocess",
		Timeout: 1,
	}

	sb, err := NewSubprocessSandbox(cfg)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}
	defer sb.Close()

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Code that takes longer than timeout
	code := "import time; time.sleep(10)"
	result, err := sb.Execute(ctx, code, models.Python)

	if err != nil {
		t.Logf("Execute returned error (expected for timeout): %v", err)
	}

	if !result.TimedOut && result.ExitCode == 0 {
		t.Error("expected timeout or non-zero exit code")
	}
}

func TestSubprocessSandbox_FileOperations(t *testing.T) {
	cfg := &Config{
		Type:    "subprocess",
		Timeout: 30,
	}

	sb, err := NewSubprocessSandbox(cfg)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}
	defer sb.Close()

	ctx := context.Background()

	// Create a temp file to upload
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello from test")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Upload file
	err = sb.UploadFile(ctx, testFile, "uploaded.txt")
	if err != nil {
		t.Errorf("UploadFile() error = %v", err)
	}

	// List files
	files, err := sb.ListFiles(ctx, ".")
	if err != nil {
		t.Errorf("ListFiles() error = %v", err)
	}

	found := false
	for _, f := range files {
		if f == "uploaded.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Error("uploaded file not found in list")
	}

	// Download file
	downloaded, err := sb.DownloadFile(ctx, "uploaded.txt")
	if err != nil {
		t.Errorf("DownloadFile() error = %v", err)
	}
	if string(downloaded) != string(content) {
		t.Errorf("DownloadFile() content = %q, want %q", downloaded, content)
	}
}

func TestSubprocessSandbox_Reset(t *testing.T) {
	cfg := &Config{
		Type:    "subprocess",
		Timeout: 30,
	}

	sb, err := NewSubprocessSandbox(cfg)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}
	defer sb.Close()

	ctx := context.Background()

	// Create a file by executing code
	code := "with open('test.txt', 'w') as f: f.write('test')"
	sb.Execute(ctx, code, models.Python)

	// Verify file exists
	files, _ := sb.ListFiles(ctx, ".")
	if len(files) == 0 {
		t.Fatal("expected file to exist before reset")
	}

	// Reset
	err = sb.Reset(ctx)
	if err != nil {
		t.Errorf("Reset() error = %v", err)
	}

	// Verify file is gone
	files, _ = sb.ListFiles(ctx, ".")
	if len(files) != 0 {
		t.Error("expected no files after reset")
	}
}

func TestParseMemory(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"512m", 512 * 1024 * 1024},
		{"512M", 512 * 1024 * 1024},
		{"1g", 1024 * 1024 * 1024},
		{"2G", 2 * 1024 * 1024 * 1024},
		{"1024", 1024},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseMemory(tt.input)
			if result != tt.expected {
				t.Errorf("parseMemory(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// Helper functions for subprocess

func parseMemory(limit string) int64 {
	limit = strings.ToLower(strings.TrimSpace(limit))

	multiplier := int64(1)
	switch {
	case strings.HasSuffix(limit, "gb"), strings.HasSuffix(limit, "g"):
		multiplier = 1024 * 1024 * 1024
		limit = strings.TrimSuffix(strings.TrimSuffix(limit, "b"), "g")
	case strings.HasSuffix(limit, "mb"), strings.HasSuffix(limit, "m"):
		multiplier = 1024 * 1024
		limit = strings.TrimSuffix(strings.TrimSuffix(limit, "b"), "m")
	case strings.HasSuffix(limit, "kb"), strings.HasSuffix(limit, "k"):
		multiplier = 1024
		limit = strings.TrimSuffix(strings.TrimSuffix(limit, "b"), "k")
	}

	var value int64
	fmt.Sscanf(limit, "%d", &value)
	return value * multiplier
}
