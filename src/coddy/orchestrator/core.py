"""Core orchestrator for managing LLM interactions and tool execution."""

from typing import Any

from coddy.config import Settings
from coddy.llm.client import LLMClient, LLMResponse, ToolCall
from coddy.sandbox.base import ExecutionResult, Sandbox
from coddy.tools.registry import ToolRegistry


class Orchestrator:
    """Orchestrates the conversation loop between user, LLM, and tools.
    
    The orchestrator manages:
    - Conversation history
    - LLM API calls
    - Tool execution
    - Session lifecycle
    
    Example:
        ```python
        settings = Settings()
        llm = LLMClient(...)
        sandbox = SubprocessSandbox()
        orchestrator = Orchestrator(llm, sandbox, settings)
        
        response = await orchestrator.handle_message("What is 2+2?")
        print(response)
        ```
    """

    def __init__(
        self,
        llm: LLMClient,
        sandbox: Sandbox,
        settings: Settings,
        system_prompt: str | None = None,
    ) -> None:
        """Initialize the orchestrator.
        
        Args:
            llm: The LLM client.
            sandbox: The sandbox for code execution.
            settings: Application settings.
            system_prompt: Optional custom system prompt.
        """
        self.llm = llm
        self.sandbox = sandbox
        self.settings = settings
        self.tool_registry = ToolRegistry()
        
        # Initialize conversation with system prompt
        self.messages: list[dict[str, Any]] = []
        if system_prompt:
            self.messages.append({"role": "system", "content": system_prompt})
        else:
            self.messages.append({
                "role": "system",
                "content": self._default_system_prompt(),
            })

    def _default_system_prompt(self) -> str:
        """Generate the default system prompt.
        
        Returns:
            The system prompt as a string.
        """
        return """You are an AI assistant with access to a code execution environment.

You can run Python and Node.js code to help users with calculations, data analysis, file processing, and more.

## When to Use Code

Use the execute_code tool for:
- Mathematical calculations
- Data analysis and visualization
- File processing (CSV, JSON, images, etc.)
- Text processing and transformation
- Any task that benefits from programmatic execution

DO NOT use code for:
- Simple factual questions (e.g., "What is the capital of France?")
- General conversation

## How to Use Code

- Use print() in Python or console.log() in Node.js to show results
- The environment is stateful: variables persist between calls
- If code fails, read the error and retry with fixes
- Maximum execution time: 30 seconds

## Handling Results

- After execution, explain the results to the user
- If generating files (charts, etc.), mention the filename
- If an error occurs, fix the code before giving up
"""

    async def handle_message(self, user_input: str) -> str:
        """Handle a user message through the full orchestration loop.
        
        This method:
        1. Adds the user message to history
        2. Calls the LLM
        3. Executes any tool calls
        4. Loops until the LLM provides a final response
        5. Returns the final response text
        
        Args:
            user_input: The user's message.
            
        Returns:
            The LLM's final response text.
        """
        # Add user message to history
        self.messages.append({"role": "user", "content": user_input})

        # Loop until we get a final response or hit the iteration limit
        for iteration in range(self.settings.max_tool_iterations):
            # Call the LLM
            response = await self.llm.chat(
                messages=self.messages,
                tools=self.tool_registry.get_tool_definitions(),
            )

            # Add the assistant's response to history
            self.messages.append(response.message)

            # Check for tool calls
            if response.has_tool_calls:
                for tool_call in response.tool_calls:
                    result = await self._execute_tool(tool_call)
                    self.messages.append({
                        "role": "tool",
                        "tool_call_id": tool_call.id,
                        "content": self._format_result(result),
                    })
            else:
                # Final response - no tool calls
                return response.text

        # Hit the iteration limit
        return "Error: Exceeded maximum tool call iterations. Please try a simpler request."

    async def _execute_tool(self, tool_call: ToolCall) -> ExecutionResult:
        """Execute a tool call.
        
        Args:
            tool_call: The tool call to execute.
            
        Returns:
            The execution result.
        """
        if tool_call.name == "execute_code":
            code = tool_call.arguments.get("code", "")
            language = tool_call.arguments.get("language", "python")
            return await self.sandbox.execute(code, language)
        else:
            return ExecutionResult(
                stdout="",
                stderr=f"Unknown tool: {tool_call.name}",
                exit_code=1,
            )

    def _format_result(self, result: ExecutionResult) -> str:
        """Format an execution result for the LLM.
        
        Args:
            result: The execution result.
            
        Returns:
            Formatted result string.
        """
        parts = []
        
        if result.stdout:
            parts.append(f"STDOUT:\n{result.stdout}")
        
        if result.stderr:
            parts.append(f"STDERR:\n{result.stderr}")
        
        if result.timed_out:
            parts.append("EXECUTION TIMED OUT")
        
        if result.output_files:
            parts.append(f"Files created: {', '.join(result.output_files)}")
        
        parts.append(f"Exit code: {result.exit_code} | Time: {result.execution_time_ms}ms")
        
        return "\n".join(parts) if parts else "Code executed successfully (no output)."

    def clear_history(self) -> None:
        """Clear the conversation history (except system prompt)."""
        if self.messages and self.messages[0].get("role") == "system":
            self.messages = [self.messages[0]]
        else:
            self.messages = []
