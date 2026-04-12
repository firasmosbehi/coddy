"""Tests for tool definitions and registry."""

import pytest

from coddy.tools.definitions import (
    EXECUTE_CODE_TOOL,
    LIST_FILES_TOOL,
    READ_FILE_TOOL,
    WRITE_FILE_TOOL,
    get_all_tools,
)
from coddy.tools.registry import ToolRegistry


class TestToolDefinitions:
    """Test cases for tool definitions."""

    def test_execute_code_structure(self) -> None:
        """Test that execute_code tool has correct structure."""
        assert EXECUTE_CODE_TOOL["type"] == "function"
        
        func = EXECUTE_CODE_TOOL["function"]
        assert func["name"] == "execute_code"
        assert "description" in func
        assert "parameters" in func
        
        params = func["parameters"]
        assert "language" in params["properties"]
        assert "code" in params["properties"]
        assert params["required"] == ["language", "code"]

    def test_execute_code_language_enum(self) -> None:
        """Test that execute_code has correct language enum."""
        language_prop = EXECUTE_CODE_TOOL["function"]["parameters"]["properties"]["language"]
        assert language_prop["enum"] == ["python", "nodejs"]

    def test_get_all_tools(self) -> None:
        """Test that get_all_tools returns all tools."""
        tools = get_all_tools()
        tool_names = [t["function"]["name"] for t in tools]
        
        assert "execute_code" in tool_names
        assert "read_file" in tool_names
        assert "write_file" in tool_names
        assert "list_files" in tool_names


class TestToolRegistry:
    """Test cases for ToolRegistry."""

    def test_default_tools_registered(self) -> None:
        """Test that default tools are registered on init."""
        registry = ToolRegistry()
        
        assert registry.has_tool("execute_code")
        assert registry.has_tool("read_file")
        assert registry.has_tool("write_file")
        assert registry.has_tool("list_files")

    def test_get_tool(self) -> None:
        """Test getting a tool definition."""
        registry = ToolRegistry()
        
        tool = registry.get("execute_code")
        assert tool is not None
        assert tool["function"]["name"] == "execute_code"

    def test_get_missing_tool(self) -> None:
        """Test getting a non-existent tool."""
        registry = ToolRegistry()
        assert registry.get("nonexistent") is None

    def test_list_tools(self) -> None:
        """Test listing all tool names."""
        registry = ToolRegistry()
        tools = registry.list_tools()
        
        assert "execute_code" in tools
        assert "read_file" in tools

    def test_register_new_tool(self) -> None:
        """Test registering a new tool."""
        registry = ToolRegistry()
        
        new_tool = {
            "type": "function",
            "function": {
                "name": "custom_tool",
                "description": "A custom tool",
                "parameters": {"type": "object", "properties": {}},
            },
        }
        
        registry.register("custom_tool", new_tool)
        assert registry.has_tool("custom_tool")
        assert registry.get("custom_tool") == new_tool

    def test_get_tool_definitions(self) -> None:
        """Test getting all tool definitions."""
        registry = ToolRegistry()
        definitions = registry.get_tool_definitions()
        
        assert len(definitions) >= 4
        tool_names = [d["function"]["name"] for d in definitions]
        assert "execute_code" in tool_names
