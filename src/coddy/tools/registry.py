"""Tool registry for mapping tool names to handlers."""

from typing import Any

from coddy.tools.definitions import get_all_tools


class ToolRegistry:
    """Registry for managing available tools.
    
    The registry maintains a mapping of tool names to their definitions
    and provides methods for tool lookup and execution.
    """

    def __init__(self) -> None:
        """Initialize the tool registry with default tools."""
        self._tools: dict[str, dict[str, Any]] = {}
        self._register_default_tools()

    def _register_default_tools(self) -> None:
        """Register the default set of tools."""
        for tool in get_all_tools():
            name = tool.get("function", {}).get("name")
            if name:
                self._tools[name] = tool

    def register(self, name: str, definition: dict[str, Any]) -> None:
        """Register a new tool.
        
        Args:
            name: The tool name.
            definition: The tool definition in OpenAI format.
        """
        self._tools[name] = definition

    def get(self, name: str) -> dict[str, Any] | None:
        """Get a tool definition by name.
        
        Args:
            name: The tool name.
            
        Returns:
            The tool definition or None if not found.
        """
        return self._tools.get(name)

    def get_tool_definitions(self) -> list[dict[str, Any]]:
        """Get all tool definitions in OpenAI format.
        
        Returns:
            List of tool definitions.
        """
        return list(self._tools.values())

    def list_tools(self) -> list[str]:
        """List all registered tool names.
        
        Returns:
            List of tool names.
        """
        return list(self._tools.keys())

    def has_tool(self, name: str) -> bool:
        """Check if a tool is registered.
        
        Args:
            name: The tool name.
            
        Returns:
            True if the tool is registered.
        """
        return name in self._tools
