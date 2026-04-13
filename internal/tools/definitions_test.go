package tools

import (
	"encoding/json"
	"testing"

	"github.com/firasmosbehi/coddy/pkg/models"
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

	// Verify parameters
	var params map[string]interface{}
	if err := json.Unmarshal(tool.Function.Parameters, &params); err != nil {
		t.Fatalf("failed to unmarshal parameters: %v", err)
	}

	if params["type"] != "object" {
		t.Error("expected parameters type to be 'object'")
	}

	properties, ok := params["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("expected properties to be an object")
	}

	if _, ok := properties["language"]; !ok {
		t.Error("expected 'language' property")
	}

	if _, ok := properties["code"]; !ok {
		t.Error("expected 'code' property")
	}

	// Check required fields
	required, ok := params["required"].([]interface{})
	if !ok {
		t.Fatal("expected required to be an array")
	}

	if len(required) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(required))
	}
}

func TestReadFileTool(t *testing.T) {
	tool := ReadFileTool()

	if tool.Function.Name != "read_file" {
		t.Errorf("expected name 'read_file', got %s", tool.Function.Name)
	}

	var params map[string]interface{}
	json.Unmarshal(tool.Function.Parameters, &params)

	properties, _ := params["properties"].(map[string]interface{})
	if _, ok := properties["path"]; !ok {
		t.Error("expected 'path' property")
	}

	required, _ := params["required"].([]interface{})
	if len(required) != 1 || required[0] != "path" {
		t.Error("expected 'path' to be required")
	}
}

func TestWriteFileTool(t *testing.T) {
	tool := WriteFileTool()

	if tool.Function.Name != "write_file" {
		t.Errorf("expected name 'write_file', got %s", tool.Function.Name)
	}

	var params map[string]interface{}
	json.Unmarshal(tool.Function.Parameters, &params)

	properties, _ := params["properties"].(map[string]interface{})
	if _, ok := properties["path"]; !ok {
		t.Error("expected 'path' property")
	}
	if _, ok := properties["content"]; !ok {
		t.Error("expected 'content' property")
	}

	required, _ := params["required"].([]interface{})
	if len(required) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(required))
	}
}

func TestListFilesTool(t *testing.T) {
	tool := ListFilesTool()

	if tool.Function.Name != "list_files" {
		t.Errorf("expected name 'list_files', got %s", tool.Function.Name)
	}

	var params map[string]interface{}
	json.Unmarshal(tool.Function.Parameters, &params)

	// path is optional, so no required check
	properties, _ := params["properties"].(map[string]interface{})
	if _, ok := properties["path"]; !ok {
		t.Error("expected 'path' property")
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
		{"write_file", true},
		{"list_files", true},
		{"nonexistent", false},
		{"", false},
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

func TestToolDefinition_ParametersValidJSON(t *testing.T) {
	defs := GetDefinitions()

	for _, def := range defs {
		t.Run(def.Function.Name, func(t *testing.T) {
			var params map[string]interface{}
			if err := json.Unmarshal(def.Function.Parameters, &params); err != nil {
				t.Errorf("invalid JSON in parameters: %v", err)
			}

			// Check type is present
			if _, ok := params["type"]; !ok {
				t.Error("parameters missing 'type' field")
			}

			// Check properties is present
			if _, ok := params["properties"]; !ok {
				t.Error("parameters missing 'properties' field")
			}
		})
	}
}

func TestToolDefinition_StructTags(t *testing.T) {
	def := models.ToolDefinition{
		Type: "function",
		Function: models.FunctionDef{
			Name:        "test",
			Description: "test description",
			Parameters:  json.RawMessage(`{}`),
		},
	}

	// Verify struct can be marshaled and unmarshaled
	data, err := json.Marshal(def)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled models.ToolDefinition
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.Type != def.Type {
		t.Error("type mismatch after unmarshal")
	}

	if unmarshaled.Function.Name != def.Function.Name {
		t.Error("name mismatch after unmarshal")
	}
}
