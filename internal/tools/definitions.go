// Package tools defines LLM tool schemas.
package tools

import (
	"encoding/json"

	"github.com/firasmosbehi/coddy/pkg/models"
)

// GetDefinitions returns all available tool definitions.
func GetDefinitions() []models.ToolDefinition {
	return []models.ToolDefinition{
		ExecuteCodeTool(),
		ReadFileTool(),
		WriteFileTool(),
		ListFilesTool(),
	}
}

// ExecuteCodeTool returns the execute_code tool definition.
func ExecuteCodeTool() models.ToolDefinition {
	params, _ := json.Marshal(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"language": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"python", "nodejs"},
				"description": "The programming language to use",
			},
			"code": map[string]interface{}{
				"type":        "string",
				"description": "The code to execute. Use print() in Python or console.log() in Node.js to show output.",
			},
		},
		"required": []string{"language", "code"},
	})

	return models.ToolDefinition{
		Type: "function",
		Function: models.FunctionDef{
			Name:        "execute_code",
			Description: "Execute code in a sandboxed environment. Use this for calculations, data processing, file creation, plotting, or any task that benefits from running code.",
			Parameters:  params,
		},
	}
}

// ReadFileTool returns the read_file tool definition.
func ReadFileTool() models.ToolDefinition {
	params, _ := json.Marshal(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The path to the file to read",
			},
		},
		"required": []string{"path"},
	})

	return models.ToolDefinition{
		Type: "function",
		Function: models.FunctionDef{
			Name:        "read_file",
			Description: "Read the contents of a file from the sandbox filesystem.",
			Parameters:  params,
		},
	}
}

// WriteFileTool returns the write_file tool definition.
func WriteFileTool() models.ToolDefinition {
	params, _ := json.Marshal(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The path to write the file to",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "The content to write to the file",
			},
		},
		"required": []string{"path", "content"},
	})

	return models.ToolDefinition{
		Type: "function",
		Function: models.FunctionDef{
			Name:        "write_file",
			Description: "Write content to a file in the sandbox filesystem.",
			Parameters:  params,
		},
	}
}

// ListFilesTool returns the list_files tool definition.
func ListFilesTool() models.ToolDefinition {
	params, _ := json.Marshal(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The directory path to list (defaults to /home/user)",
			},
		},
	})

	return models.ToolDefinition{
		Type: "function",
		Function: models.FunctionDef{
			Name:        "list_files",
			Description: "List files in a directory in the sandbox filesystem.",
			Parameters:  params,
		},
	}
}

// GetToolNames returns a list of available tool names.
func GetToolNames() []string {
	defs := GetDefinitions()
	names := make([]string, len(defs))
	for i, def := range defs {
		names[i] = def.Function.Name
	}
	return names
}

// FindTool finds a tool by name.
func FindTool(name string) (models.ToolDefinition, bool) {
	for _, def := range GetDefinitions() {
		if def.Function.Name == name {
			return def, true
		}
	}
	return models.ToolDefinition{}, false
}
