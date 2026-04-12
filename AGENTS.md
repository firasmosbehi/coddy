# AGENTS.md — LLM Code Execution Environment

## Project Goal

Build an orchestration layer that gives any LLM (Qwen, Llama, Mistral, etc.) the ability to write and execute code in a sandboxed environment, similar to how Claude uses its coding environment. The LLM can run Python or Node.js, see the output, iterate, and deliver results to the user.

---

## Architecture Overview

```
User Message
     │
     ▼
┌──────────────────────┐
│   Orchestrator        │  ← Python server (FastAPI or CLI)
│                       │
│  1. Build prompt      │
│  2. Call LLM          │
│  3. Parse response    │
│  4. Route tool calls  │
│  5. Loop until done   │
└──────┬───────────────┘
       │
       ├──► LLM API (OpenAI-compatible)
       │      - Qwen via Ollama / vLLM / cloud
       │      - Any model with tool-calling support
       │
       └──► Sandbox (code execution)
              - Docker container (primary)
              - E2B (alternative, managed)
              - Subprocess (dev only, not secure)
```

---

## Core Loop (the brain of the system)

This is the main execution cycle. Every user interaction follows this flow:

```
1. User sends message
2. Orchestrator appends message to conversation history
3. Orchestrator calls LLM with conversation + tool definitions
4. LLM responds with either:
   a. TOOL CALL → extract code → execute in sandbox → append result → GOTO 3
   b. TEXT RESPONSE → display to user → DONE
5. Repeat step 4 until the LLM produces a final text response
```

Important: the loop must support **multiple sequential tool calls** in a single turn. The LLM might run code, see an error, fix it, and run again before responding.

---

## Components to Implement

### 1. Sandbox Manager (`sandbox/`)

Responsible for executing arbitrary code safely and returning results.

**Interface:**

```python
class Sandbox:
    def execute(self, code: str, language: str, timeout: int = 30) -> ExecutionResult:
        """Run code and return stdout, stderr, exit_code, and output files."""

    def upload_file(self, local_path: str, sandbox_path: str) -> None:
        """Copy a file into the sandbox filesystem."""

    def download_file(self, sandbox_path: str) -> bytes:
        """Retrieve a file from the sandbox filesystem."""

    def reset(self) -> None:
        """Reset the sandbox to a clean state."""
```

**ExecutionResult schema:**

```python
@dataclass
class ExecutionResult:
    stdout: str
    stderr: str
    exit_code: int
    output_files: list[str]    # paths to files created during execution
    execution_time_ms: int
    timed_out: bool
```

**Implementation options (pick one):**

| Option | Security | Complexity | Best For |
|---|---|---|---|
| `subprocess` + `ulimit` | Low | Low | Local dev/testing only |
| Docker containers | Medium | Medium | Self-hosted production |
| gVisor / Firecracker | High | High | Multi-tenant production |
| E2B SDK | High | Low | Fast prototyping, managed |

**Docker approach (recommended starting point):**

- Base image with Python 3.12 + Node.js 20 + common packages pre-installed
- `--network=none` to disable network access (enable selectively if needed)
- Memory limit: 512MB, CPU limit: 1 core
- Execution timeout: 30 seconds default, 120 seconds max
- Read-only root filesystem, writable `/tmp` and `/home/user`
- One container per session (stateful) OR per execution (stateless)

**Pre-installed packages in the sandbox image:**

Python: `numpy`, `pandas`, `matplotlib`, `scipy`, `scikit-learn`, `requests`, `beautifulsoup4`, `Pillow`, `sympy`, `openpyxl`

Node.js: `lodash`, `axios`, `cheerio`, `csv-parse`, `mathjs`, `sharp`

### 2. LLM Client (`llm/`)

Handles communication with the LLM. Must support the OpenAI-compatible chat completions API since most local/cloud providers expose this format.

**Interface:**

```python
class LLMClient:
    def __init__(self, base_url: str, model: str, api_key: str = ""):
        ...

    def chat(self, messages: list[dict], tools: list[dict]) -> LLMResponse:
        """Send a chat completion request with tool definitions."""
```

**Supported backends (all use the same OpenAI-compatible API):**

- Ollama: `base_url="http://localhost:11434/v1"`
- vLLM: `base_url="http://localhost:8000/v1"`
- Together AI: `base_url="https://api.together.xyz/v1"`
- OpenRouter: `base_url="https://openrouter.ai/api/v1"`
- Any OpenAI-compatible provider

**Fallback for models without native tool calling:**

If the model doesn't support the `tools` parameter, use prompt-based tool calling instead. Instruct the model to emit code blocks in a structured format:

```xml
<tool_call>
<name>execute_code</name>
<arguments>
{"language": "python", "code": "print('hello world')"}
</arguments>
</tool_call>
```

The orchestrator must parse this from the raw text response.

### 3. Tool Definitions (`tools/`)

Define the tools the LLM can use. Start with one, expand later.

**Primary tool — `execute_code`:**

```json
{
  "type": "function",
  "function": {
    "name": "execute_code",
    "description": "Execute code in a sandboxed environment. Use this for calculations, data processing, file creation, plotting, or any task that benefits from running code. The environment is stateful: variables and files persist between calls within the same session.",
    "parameters": {
      "type": "object",
      "properties": {
        "language": {
          "type": "string",
          "enum": ["python", "nodejs"],
          "description": "The programming language to use"
        },
        "code": {
          "type": "string",
          "description": "The code to execute. For Python, use print() to show output. For Node.js, use console.log()."
        }
      },
      "required": ["language", "code"]
    }
  }
}
```

**Future tools to add:**

- `read_file` — read a user-uploaded file
- `write_file` — create a file for the user to download
- `web_fetch` — fetch a URL (requires enabling network in sandbox)
- `list_files` — list files in the working directory

### 4. Orchestrator (`orchestrator/`)

The main loop that ties everything together.

```python
class Orchestrator:
    def __init__(self, llm: LLMClient, sandbox: Sandbox, system_prompt: str):
        self.llm = llm
        self.sandbox = sandbox
        self.messages = [{"role": "system", "content": system_prompt}]
        self.tools = [EXECUTE_CODE_TOOL]
        self.max_tool_iterations = 10  # prevent infinite loops

    def handle_user_message(self, user_input: str) -> str:
        self.messages.append({"role": "user", "content": user_input})

        for _ in range(self.max_tool_iterations):
            response = self.llm.chat(self.messages, self.tools)
            self.messages.append(response.message)

            if response.has_tool_calls():
                for call in response.tool_calls:
                    result = self._execute_tool(call)
                    self.messages.append({
                        "role": "tool",
                        "tool_call_id": call.id,
                        "content": self._format_result(result)
                    })
            else:
                return response.text

        return "Error: exceeded maximum tool call iterations."

    def _execute_tool(self, call) -> ExecutionResult:
        args = call.arguments
        return self.sandbox.execute(
            code=args["code"],
            language=args["language"]
        )

    def _format_result(self, result: ExecutionResult) -> str:
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
```

### 5. System Prompt (`prompts/`)

The system prompt is critical. It teaches the LLM when and how to use the code execution tool.

**Key sections the system prompt must cover:**

```
1. ROLE
   You are an AI assistant with access to a code execution environment.
   You can run Python and Node.js code to help users.

2. WHEN TO USE CODE
   - Calculations, math, data analysis
   - File processing (CSV, JSON, images, etc.)
   - Generating charts, plots, visualizations
   - Any task that benefits from programmatic execution
   - DO NOT use code for simple factual questions or conversation

3. HOW TO USE CODE
   - Use the execute_code tool
   - The environment is stateful (variables persist between calls)
   - Use print() / console.log() to show intermediate results
   - If code errors, read the traceback, fix the issue, and retry
   - Install packages with pip/npm if needed (network may be disabled)

4. HANDLING RESULTS
   - After execution, explain the results to the user in natural language
   - If generating a file (chart, document, etc.), tell the user the filename
   - If an error occurs, fix the code and retry before giving up

5. CONSTRAINTS
   - Maximum execution time: 30 seconds per call
   - Maximum 10 consecutive tool calls per turn
   - Available pre-installed packages: [list them]
```

Store the system prompt in a separate file (`prompts/system.txt`) so it can be iterated on without code changes.

---

## File Structure

```
llm-code-executor/
├── AGENTS.md                  ← this file
├── README.md
├── requirements.txt
├── docker/
│   ├── Dockerfile             ← sandbox base image
│   └── entrypoint.sh
├── src/
│   ├── main.py                ← entry point (CLI or API server)
│   ├── orchestrator.py        ← core loop
│   ├── llm_client.py          ← LLM API client
│   ├── sandbox/
│   │   ├── base.py            ← abstract Sandbox interface
│   │   ├── docker_sandbox.py  ← Docker implementation
│   │   ├── subprocess_sandbox.py  ← dev-only fallback
│   │   └── e2b_sandbox.py     ← E2B managed implementation
│   ├── tools/
│   │   ├── definitions.py     ← tool JSON schemas
│   │   └── registry.py        ← tool name → handler mapping
│   └── prompts/
│       └── system.txt         ← system prompt template
├── tests/
│   ├── test_sandbox.py
│   ├── test_orchestrator.py
│   └── test_tool_parsing.py
└── examples/
    ├── cli_chat.py            ← simple CLI demo
    └── api_server.py          ← FastAPI web server
```

---

## Implementation Order

Follow this sequence. Each phase builds on the previous one.

### Phase 1 — Minimal Working Prototype

**Goal:** User can chat with Qwen in a terminal and it can run Python code.

1. Implement `subprocess_sandbox.py` (unsafe but fast to build)
2. Implement `llm_client.py` targeting Ollama
3. Implement `orchestrator.py` with the core loop
4. Write a basic system prompt
5. Create `cli_chat.py` that ties it all together
6. Test: "What is 2^100?" → LLM should write Python and return the answer

### Phase 2 — Docker Sandbox

**Goal:** Code execution is isolated and safe.

1. Create `Dockerfile` with Python + Node.js + packages
2. Implement `docker_sandbox.py` using the Docker SDK
3. Add timeout handling, memory limits, network isolation
4. Add file upload/download support
5. Switch the orchestrator from subprocess to Docker
6. Test: malicious code (fork bombs, infinite loops) should be contained

### Phase 3 — Stateful Sessions

**Goal:** Variables and files persist across multiple code executions in one conversation.

1. Keep containers alive for the duration of a session
2. Implement session ID tracking
3. Add cleanup logic (destroy container after N minutes of inactivity)
4. Test: define a variable in one call, use it in the next

### Phase 4 — Web API

**Goal:** Expose the system as an HTTP API.

1. Build FastAPI server with endpoints:
   - `POST /chat` — send message, get response (streaming optional)
   - `POST /sessions` — create a new session
   - `DELETE /sessions/{id}` — destroy a session
   - `POST /sessions/{id}/upload` — upload a file
   - `GET /sessions/{id}/files/{path}` — download a file
2. Add WebSocket support for streaming responses
3. Build a simple web UI (optional)

### Phase 5 — Hardening

**Goal:** Production-ready.

1. Add structured logging
2. Add rate limiting
3. Add conversation history persistence (SQLite or Redis)
4. Add support for multiple concurrent sessions
5. Add metrics (execution time, error rates, tool call frequency)
6. Write comprehensive tests
7. Add support for more LLM providers

---

## Configuration

Use environment variables or a config file:

```env
# LLM
LLM_BASE_URL=http://localhost:11434/v1
LLM_MODEL=qwen3-coder
LLM_API_KEY=                       # empty for Ollama

# Sandbox
SANDBOX_TYPE=docker                # docker | subprocess | e2b
SANDBOX_IMAGE=llm-sandbox:latest
SANDBOX_TIMEOUT=30                 # seconds
SANDBOX_MEMORY_LIMIT=512m
SANDBOX_CPU_LIMIT=1
SANDBOX_NETWORK=none               # none | bridge

# Orchestrator
MAX_TOOL_ITERATIONS=10
SESSION_TIMEOUT=3600               # seconds before session cleanup

# E2B (if using managed sandbox)
E2B_API_KEY=
```

---

## Key Design Decisions

| Decision | Recommendation | Rationale |
|---|---|---|
| Stateful vs stateless sandbox | **Stateful** (per session) | Much more useful — users expect variables to persist like a notebook |
| Container per session vs per execution | **Per session** | Avoids cold start overhead on every tool call |
| Tool calling: native vs prompt-based | **Native first**, prompt-based fallback | Native is more reliable but not all models support it |
| Streaming | **Yes, for text responses** | Better UX; tool call results can be non-streaming |
| Max iterations | **10** | Prevents runaway loops while allowing complex multi-step tasks |

---

## Testing Checklist

Use these prompts to verify the system works end-to-end:

- [ ] "What is the 50th Fibonacci number?" → should write and run code
- [ ] "Create a bar chart of the top 5 most populous countries" → should generate a matplotlib plot
- [ ] "Read this CSV and tell me the average of column X" → should handle file I/O
- [ ] "Write a function to sort a list, then test it" → should use stateful execution
- [ ] "What's the capital of France?" → should NOT use code, just answer
- [ ] Send code that runs an infinite loop → should timeout gracefully
- [ ] Send code that tries to access the network → should be blocked (if network=none)
- [ ] Send code that tries to read /etc/passwd → should be blocked or return nothing useful