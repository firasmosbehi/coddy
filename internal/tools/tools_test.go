package tools

import (
	"testing"
)

func TestGetDefinitions(t *testing.T) {
	defs := GetDefinitions()

	if len(defs) != 4 {
		t.Errorf("expected 4 tool definitions, got %d", len(defs))
	}

	expectedNames := map[string]bool{
		"execute_code": false,
		"read_file":    false,
		"write_file":   false,
		"list_files":   false,
	}

	for _, def := range defs {
		expectedNames[def.Function.Name] = true
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("expected tool %s not found", name)
		}
	}
}

func TestExecuteCodeTool(t *testing.T) {
	tool := ExecuteCodeTool()

	if tool.Type != "function" {
		t.Errorf("expected type 'function', got %s", tool.Type)
	}

	if tool.Function.Name != "execute_code" {
		t.Errorf("expected name 'execute_code', got %s", tool.Function.Name)
	}

	if tool.Function.Description == "" {
		t.Error("expected non-empty description")
	}
}

func TestGetToolNames(t *testing.T) {
	names := GetToolNames()

	if len(names) != 4 {
		t.Errorf("expected 4 tool names, got %d", len(names))
	}

	for _, name := range names {
		if name == "" {
			t.Error("expected non-empty tool name")
		}
	}
}

func TestFindTool(t *testing.T) {
	tests := []struct {
		name      string
		wantFound bool
	}{
		{"execute_code", true},
		{"read_file", true},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool, found := FindTool(tt.name)
			if found != tt.wantFound {
				t.Errorf("FindTool() found = %v, want %v", found, tt.wantFound)
			}
			if found && tool.Function.Name != tt.name {
				t.Errorf("FindTool() name = %s, want %s", tool.Function.Name, tt.name)
			}
		})
	}
}
