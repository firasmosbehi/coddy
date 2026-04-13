package sandbox

import (
	"context"
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
			name:       "python multiline",
			code:       "x = 10\ny = 20\nprint(x + y)",
			language:   models.Python,
			wantStdout: "30\n",
			wantErr:    false,
		},
		{
			name:       "nodejs hello world",
			code:       "console.log('Hello from Node!')",
			language:   models.NodeJS,
			wantStdout: "Hello from Node!\n",
			wantErr:    false,
		},
		{
			name:       "nodejs calculation",
			code:       "console.log(2 + 2)",
			language:   models.NodeJS,
			wantStdout: "4\n",
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
			if result.TimedOut {
				t.Error("Execute() should not timeout")
			}
		})
	}
}

func TestSubprocessSandbox_Execute_Error(t *testing.T) {
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

	// Python error
	result, err := sb.Execute(ctx, "print(undefined_var)", models.Python)
	if err != nil {
		t.Errorf("Execute() should not return error for code failure: %v", err)
	}

	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code for error")
	}

	if result.Stderr == "" {
		t.Error("expected stderr for error")
	}
}

func TestSubprocessSandbox_Execute_InvalidLanguage(t *testing.T) {
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

	result, err := sb.Execute(ctx, "print('hello')", models.Language("invalid"))
	if err == nil {
		t.Error("expected error for invalid language")
	}

	if result.ExitCode != -1 {
		t.Error("expected exit code -1 for invalid language")
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

func TestSubprocessSandbox_UploadFile_NonExistent(t *testing.T) {
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

	err = sb.UploadFile(ctx, "/non/existent/file.txt", "dest.txt")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestSubprocessSandbox_DownloadFile_NonExistent(t *testing.T) {
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

	_, err = sb.DownloadFile(ctx, "non_existent.txt")
	if err == nil {
		t.Error("expected error for non-existent file")
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

func TestSubprocessSandbox_OutputFileDetection(t *testing.T) {
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

	// Create a file
	code := `
with open('output.txt', 'w') as f:
    f.write('Hello World')
print("File created")
`
	result, err := sb.Execute(ctx, code, models.Python)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(result.OutputFiles) != 1 {
		t.Errorf("expected 1 output file, got %d", len(result.OutputFiles))
	}

	if len(result.OutputFiles) > 0 && result.OutputFiles[0] != "output.txt" {
		t.Errorf("expected output.txt, got %s", result.OutputFiles[0])
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
		{"1gb", 1 * 1024 * 1024 * 1024},
		{"100k", 100 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseMemory(tt.input)
			if result != tt.expected {
				t.Errorf("ParseMemory(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDiffStringSlices(t *testing.T) {
	before := []string{"a.txt", "b.txt"}
	after := []string{"a.txt", "b.txt", "c.txt", "d.txt"}

	result := DiffStringSlices(before, after)

	if len(result) != 2 {
		t.Errorf("expected 2 new files, got %d", len(result))
	}

	foundC := false
	foundD := false
	for _, f := range result {
		if f == "c.txt" {
			foundC = true
		}
		if f == "d.txt" {
			foundD = true
		}
	}

	if !foundC || !foundD {
		t.Error("expected c.txt and d.txt in result")
	}
}

func TestDiffStringSlices_NoChange(t *testing.T) {
	before := []string{"a.txt", "b.txt"}
	after := []string{"a.txt", "b.txt"}

	result := DiffStringSlices(before, after)

	if len(result) != 0 {
		t.Errorf("expected 0 new files, got %d", len(result))
	}
}

func TestSubprocessSandbox_Factory(t *testing.T) {
	cfg := &Config{
		Type:    "subprocess",
		Timeout: 30,
	}

	sb, err := Factory(cfg)
	if err != nil {
		t.Fatalf("Factory() error = %v", err)
	}

	if sb == nil {
		t.Error("expected sandbox to be created")
	}

	// Verify it's the right type
	_, ok := sb.(*SubprocessSandbox)
	if !ok {
		t.Error("expected *SubprocessSandbox type")
	}

	sb.Close()
}

func TestSubprocessSandbox_Factory_InvalidType(t *testing.T) {
	cfg := &Config{
		Type: "invalid",
	}

	_, err := Factory(cfg)
	if err == nil {
		t.Error("expected error for invalid sandbox type")
	}

	if !strings.Contains(err.Error(), "unsupported sandbox type") {
		t.Errorf("expected 'unsupported sandbox type' error, got %v", err)
	}
}

func TestExecutionResult_Success(t *testing.T) {
	tests := []struct {
		name     string
		result   models.ExecutionResult
		expected bool
	}{
		{
			name:     "success",
			result:   models.ExecutionResult{ExitCode: 0, TimedOut: false},
			expected: true,
		},
		{
			name:     "failure exit code",
			result:   models.ExecutionResult{ExitCode: 1, TimedOut: false},
			expected: false,
		},
		{
			name:     "timeout",
			result:   models.ExecutionResult{ExitCode: 0, TimedOut: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.Success(); got != tt.expected {
				t.Errorf("Success() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExecutionResult_String(t *testing.T) {
	result := models.ExecutionResult{
		Stdout:          "Hello",
		Stderr:          "Warning",
		ExitCode:        0,
		OutputFiles:     []string{"file.txt"},
		ExecutionTimeMs: 100,
		TimedOut:        false,
	}

	str := result.String()

	if !strings.Contains(str, "STDOUT:") {
		t.Error("expected STDOUT in string representation")
	}

	if !strings.Contains(str, "STDERR:") {
		t.Error("expected STDERR in string representation")
	}

	if !strings.Contains(str, "file.txt") {
		t.Error("expected output file in string representation")
	}

	if !strings.Contains(str, "Exit code: 0") {
		t.Error("expected exit code in string representation")
	}
}
