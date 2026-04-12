"""Tool definitions for LLM tool calling.

These definitions follow the OpenAI function calling format.
"""

from typing import Any


def get_all_tools() -> list[dict[str, Any]]:
    """Get all available tool definitions.
    
    Returns:
        List of tool definitions in OpenAI format.
    """
    return [
        EXECUTE_CODE_TOOL,
        READ_FILE_TOOL,
        WRITE_FILE_TOOL,
        LIST_FILES_TOOL,
    ]


# Primary tool: Execute code in the sandbox
EXECUTE_CODE_TOOL: dict[str, Any] = {
    "type": "function",
    "function": {
        "name": "execute_code",
        "description": (
            "Execute code in a sandboxed environment. Use this for calculations, "
            "data processing, file creation, plotting, or any task that benefits "
            "from running code. The environment is stateful: variables and files "
            "persist between calls within the same session."
        ),
        "parameters": {
            "type": "object",
            "properties": {
                "language": {
                    "type": "string",
                    "enum": ["python", "nodejs"],
                    "description": (
                        "The programming language to use. Use 'python' for Python 3.12 "
                        "or 'nodejs' for Node.js 20."
                    ),
                },
                "code": {
                    "type": "string",
                    "description": (
                        "The code to execute. For Python, use print() to show output. "
                        "For Node.js, use console.log(). Code should be complete and "
                        "self-contained."
                    ),
                },
            },
            "required": ["language", "code"],
        },
    },
}

# Read a file from the sandbox
READ_FILE_TOOL: dict[str, Any] = {
    "type": "function",
    "function": {
        "name": "read_file",
        "description": "Read the contents of a file from the sandbox filesystem.",
        "parameters": {
            "type": "object",
            "properties": {
                "path": {
                    "type": "string",
                    "description": (
                        "The path to the file to read, relative to /home/user or "
                        "absolute. Examples: 'data.csv', '/home/user/output.txt'"
                    ),
                },
            },
            "required": ["path"],
        },
    },
}

# Write a file to the sandbox
WRITE_FILE_TOOL: dict[str, Any] = {
    "type": "function",
    "function": {
        "name": "write_file",
        "description": (
            "Write content to a file in the sandbox filesystem. "
            "Creates the file if it doesn't exist, overwrites if it does."
        ),
        "parameters": {
            "type": "object",
            "properties": {
                "path": {
                    "type": "string",
                    "description": (
                        "The path to write the file to, relative to /home/user or "
                        "absolute. Examples: 'output.txt', '/home/user/data.json'"
                    ),
                },
                "content": {
                    "type": "string",
                    "description": "The content to write to the file.",
                },
            },
            "required": ["path", "content"],
        },
    },
}

# List files in the sandbox
LIST_FILES_TOOL: dict[str, Any] = {
    "type": "function",
    "function": {
        "name": "list_files",
        "description": "List files in a directory in the sandbox filesystem.",
        "parameters": {
            "type": "object",
            "properties": {
                "path": {
                    "type": "string",
                    "description": (
                        "The directory path to list, relative to /home/user or "
                        "absolute. Defaults to /home/user if not specified."
                    ),
                },
            },
        },
    },
}
