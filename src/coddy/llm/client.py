"""LLM client for interacting with OpenAI-compatible APIs."""

from dataclasses import dataclass, field
from typing import Any


@dataclass
class ToolCall:
    """Represents a tool call from the LLM.
    
    Attributes:
        id: Unique identifier for the tool call.
        name: Name of the tool to call.
        arguments: Arguments for the tool call (parsed JSON).
    """

    id: str
    name: str
    arguments: dict[str, Any]


@dataclass
class LLMResponse:
    """Response from the LLM.
    
    Attributes:
        text: The text content of the response.
        tool_calls: List of tool calls to execute.
        message: The raw message dict for conversation history.
        usage: Token usage information.
    """

    text: str
    tool_calls: list[ToolCall] = field(default_factory=list)
    message: dict[str, Any] = field(default_factory=dict)
    usage: dict[str, int] = field(default_factory=dict)

    @property
    def has_tool_calls(self) -> bool:
        """Check if the response contains tool calls."""
        return len(self.tool_calls) > 0


class LLMClient:
    """Client for interacting with OpenAI-compatible LLM APIs."""

    def __init__(
        self,
        base_url: str,
        model: str,
        api_key: str = "",
        timeout: float = 60.0,
    ) -> None:
        """Initialize the LLM client.
        
        Args:
            base_url: Base URL for the API (e.g., http://localhost:11434/v1).
            model: Name of the model to use.
            api_key: API key for authentication (optional for local models).
            timeout: Request timeout in seconds.
        """
        self.base_url = base_url.rstrip("/")
        self.model = model
        self.api_key = api_key or "not-needed"
        self.timeout = timeout

    async def chat(
        self,
        messages: list[dict[str, Any]],
        tools: list[dict[str, Any]] | None = None,
        temperature: float = 0.7,
        max_tokens: int | None = None,
    ) -> LLMResponse:
        """Send a chat completion request.
        
        Args:
            messages: List of message dictionaries.
            tools: Optional list of tool definitions.
            temperature: Sampling temperature.
            max_tokens: Maximum tokens to generate.
            
        Returns:
            LLMResponse containing the response text and any tool calls.
            
        Raises:
            NotImplementedError: This is a placeholder implementation.
        """
        raise NotImplementedError("LLM client implementation pending Phase 1")

    def _parse_tool_calls(self, response_data: dict[str, Any]) -> list[ToolCall]:
        """Parse tool calls from the API response.
        
        Args:
            response_data: Raw response from the API.
            
        Returns:
            List of ToolCall objects.
        """
        tool_calls = []
        message = response_data.get("choices", [{}])[0].get("message", {})
        
        for tc in message.get("tool_calls", []):
            if tc.get("type") == "function":
                function = tc.get("function", {})
                import json
                try:
                    arguments = json.loads(function.get("arguments", "{}"))
                except json.JSONDecodeError:
                    arguments = {}
                
                tool_calls.append(
                    ToolCall(
                        id=tc.get("id", ""),
                        name=function.get("name", ""),
                        arguments=arguments,
                    )
                )
        
        return tool_calls
