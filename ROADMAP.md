# Coddy — LLM Code Execution Environment

> **Vision**: Give every LLM the power to write, execute, and iterate on code safely — just like Claude's coding environment, but model-agnostic and self-hostable.

---

## 1. Problem Statement

Current LLM interactions are limited to text-in, text-out. When users ask for:
- Complex calculations
- Data analysis and visualization  
- File processing
- Code verification

...the LLM can only *describe* what code would do, not *actually run it* and show results. This creates a gap between intention and execution.

### Why Existing Solutions Fall Short

| Solution | Limitation |
|----------|------------|
| Claude's built-in environment | Proprietary, only works with Claude |
| ChatGPT Code Interpreter | Closed, API not available for custom use |
| Jupyter notebooks | Manual execution, no LLM integration |
| E2B | Managed service, requires external API key |

### The Opportunity

Build an **open, model-agnostic orchestration layer** that:
- Works with any OpenAI-compatible LLM (local or cloud)
- Runs code in isolated, configurable sandboxes
- Can be self-hosted or embedded in other applications
- Provides both CLI and API interfaces

---

## 2. Core Concepts

### 2.1 The "Coddy" Loop

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   User      │────▶│ Orchestrator │────▶│     LLM     │
│  Request    │     │   (Coddy)    │     │             │
└─────────────┘     └──────────────┘     └──────┬──────┘
       ▲                                          │
       │                                          │
       │    ┌──────────────┐    ┌─────────────┐  │
       └────│   Sandbox    │◀───│  Tool Call  │◀─┘
            │ (execution)  │    │  execute()  │
            └──────────────┘    └─────────────┘
```

### 2.2 Key Differentiators

| Feature | Coddy |
|---------|-------|
| **Model Agnostic** | Works with Qwen, Llama, Mistral, GPT-4, etc. |
| **Multiple Backends** | Subprocess (dev), Docker (self-hosted), E2B (managed) |
| **Stateful Sessions** | Variables and files persist across tool calls |
| **Dual Interface** | CLI for quick tasks, API for integrations |
| **Streaming Support** | Real-time response streaming |
| **File Operations** | Upload, process, download files |

### 2.3 Supported Languages (MVP)

- **Python 3.12** — Primary language with data science stack
- **Node.js 20** — Secondary language for JavaScript/TypeScript tasks

### 2.4 Session Model

A **Session** is a stateful conversation context with:
- A dedicated sandbox container
- Conversation history
- Persistent filesystem (`/home/user`)
- Configurable timeout and resource limits

Sessions are created on first use and destroyed after inactivity or explicit cleanup.

---

## 3. Architecture Overview

```
coddy/
├── API Layer (FastAPI)
│   ├── REST endpoints for chat, sessions, files
│   └── WebSocket for streaming
│
├── Orchestrator
│   ├── Message history management
│   ├── Tool call routing
│   ├── Session lifecycle
│   └── Error handling & retries
│
├── LLM Client
│   ├── OpenAI-compatible API client
│   ├── Native tool calling support
│   └── Prompt-based fallback
│
├── Sandbox Manager
│   ├── Abstract Sandbox interface
│   ├── Docker implementation (primary)
│   ├── Subprocess implementation (dev)
│   └── E2B implementation (optional)
│
└── Tool Registry
    ├── execute_code
    ├── read_file
    ├── write_file
    └── list_files
```

---

## 4. Roadmap & Milestones

### Phase 0: Foundation (Week 1)
**Goal**: Project structure, tooling, and planning

| Task | Deliverable |
|------|-------------|
| Set up Python project structure | `src/`, `tests/`, `docker/` directories |
| Create virtual environment | `requirements.txt` with dependencies |
| Set up linting/formatting | `ruff`, `mypy` configuration |
| Write initial README | Basic setup instructions |
| Design data models | Pydantic models for all schemas |

**Dependencies to install:**
```
fastapi, uvicorn, pydantic, httpx, docker, pytest, pytest-asyncio
```

---

### Phase 1: Core Engine (Weeks 2-3)
**Goal**: Minimal working system with subprocess sandbox

#### 1.1 Sandbox Interface
- [ ] Abstract `Sandbox` base class
- [ ] `SubprocessSandbox` implementation
- [ ] `ExecutionResult` dataclass
- [ ] Timeout handling

#### 1.2 LLM Client
- [ ] `LLMClient` class
- [ ] OpenAI-compatible API support
- [ ] Response parsing (native tool calls)
- [ ] Prompt-based tool call fallback

#### 1.3 Tool System
- [ ] Tool definition schemas
- [ ] Tool registry
- [ ] `execute_code` tool implementation

#### 1.4 Orchestrator
- [ ] Core execution loop
- [ ] Message history management
- [ ] Tool call dispatch
- [ ] Result formatting

#### 1.5 CLI Interface
- [ ] Basic chat CLI
- [ ] Configuration via environment variables
- [ ] Simple REPL

**Milestone 1 Demo:**
```
$ python -m coddy.cli
> What is 2^100?
[Tool Call: execute_code]
STDOUT: 1267650600228229401496703205376
Exit code: 0

The result of 2^100 is 1,267,650,600,228,229,401,496,703,205,376.
```

---

### Phase 2: Secure Sandbox (Weeks 4-5)
**Goal**: Production-ready Docker sandbox

#### 2.1 Docker Infrastructure
- [ ] Create `Dockerfile` with Python + Node.js
- [ ] Pre-install packages (numpy, pandas, matplotlib, etc.)
- [ ] Security hardening (non-root user, read-only fs)
- [ ] Entrypoint script

#### 2.2 Docker Sandbox Implementation
- [ ] `DockerSandbox` class
- [ ] Container lifecycle management
- [ ] Volume mounting for persistence
- [ ] Network isolation (`--network=none`)
- [ ] Resource limits (memory, CPU)

#### 2.3 File Operations
- [ ] `upload_file()` method
- [ ] `download_file()` method
- [ ] `list_files()` tool
- [ ] Output file detection

#### 2.4 Safety Features
- [ ] Timeout enforcement (30s default, 120s max)
- [ ] Memory limit enforcement (512MB default)
- [ ] CPU limit (1 core)
- [ ] Infinite loop detection

**Milestone 2 Demo:**
```
$ python -m coddy.cli
> Create a plot of sin(x) from 0 to 2π
[Tool Call: execute_code]
Files created: /home/user/plot.png
Exit code: 0

I've created a plot of sin(x). The file is available at plot.png.
```

---

### Phase 3: Stateful Sessions (Week 6)
**Goal**: Persistent sessions across multiple interactions

#### 3.1 Session Management
- [ ] `Session` class
- [ ] Session ID generation (UUID)
- [ ] In-memory session store
- [ ] Session timeout handling

#### 3.2 Stateful Execution
- [ ] Container reuse per session
- [ ] Variable persistence between calls
- [ ] Filesystem persistence

#### 3.3 Session Cleanup
- [ ] Background cleanup task
- [ ] Destroy after inactivity (1 hour default)
- [ ] Explicit session deletion

#### 3.4 Configuration
- [ ] `.env` file support
- [ ] Session timeout configuration
- [ ] Sandbox type selection

**Milestone 3 Demo:**
```
$ python -m coddy.cli
> x = 42
[Tool Call: execute_code]
Exit code: 0

> print(f"The value is {x}")
[Tool Call: execute_code]
STDOUT: The value is 42
Exit code: 0
```

---

### Phase 4: Web API (Weeks 7-8)
**Goal**: HTTP API and WebSocket support

#### 4.1 FastAPI Server
- [ ] Project structure for API
- [ ] Health check endpoint
- [ ] Error handling middleware

#### 4.2 REST Endpoints
- [ ] `POST /sessions` — Create session
- [ ] `DELETE /sessions/{id}` — Destroy session
- [ ] `POST /sessions/{id}/chat` — Send message
- [ ] `POST /sessions/{id}/upload` — Upload file
- [ ] `GET /sessions/{id}/files/{path}` — Download file
- [ ] `GET /sessions/{id}/files` — List files

#### 4.3 WebSocket Support
- [ ] `ws://` endpoint for streaming
- [ ] Streaming LLM responses
- [ ] Real-time tool call notifications

#### 4.4 API Documentation
- [ ] Auto-generated OpenAPI docs
- [ ] Example requests/responses
- [ ] Postman collection (optional)

**Milestone 4 Demo:**
```bash
# Create a session
curl -X POST http://localhost:8000/sessions
# { "session_id": "abc-123", "created_at": "..." }

# Chat with the session
curl -X POST http://localhost:8000/sessions/abc-123/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "What is 2+2?"}'
```

---

### Phase 5: Enhanced Tools (Week 9)
**Goal**: Expanded tool capabilities

#### 5.1 File Tools
- [ ] `read_file` tool — read files from sandbox
- [ ] `write_file` tool — write text files
- [ ] `list_files` tool — directory listing

#### 5.2 Web Tool (Optional)
- [ ] `web_fetch` tool — fetch URLs
- [ ] Requires network-enabled sandbox
- [ ] HTML/text extraction

#### 5.3 Tool Chaining
- [ ] Multi-tool call support
- [ ] Tool call dependencies
- [ ] Parallel tool execution

---

### Phase 6: Hardening & Polish (Week 10)
**Goal**: Production-ready system

#### 6.1 Testing
- [ ] Unit tests for all components
- [ ] Integration tests for API
- [ ] Sandbox security tests
- [ ] Test coverage > 80%

#### 6.2 Observability
- [ ] Structured logging (JSON)
- [ ] Execution metrics
- [ ] Error tracking
- [ ] Performance monitoring

#### 6.3 Documentation
- [ ] Comprehensive README
- [ ] API documentation
- [ ] Deployment guide
- [ ] Contributing guide

#### 6.4 Examples
- [ ] CLI chat example
- [ ] API client example (Python)
- [ ] Jupyter notebook integration
- [ ] Docker Compose setup

---

## 5. Technical Specifications

### 5.1 Data Models

```python
# Execution Result
class ExecutionResult(BaseModel):
    stdout: str
    stderr: str
    exit_code: int
    output_files: list[str]
    execution_time_ms: int
    timed_out: bool

# Session
class Session(BaseModel):
    id: str
    created_at: datetime
    last_activity: datetime
    sandbox_type: str
    messages: list[dict]

# Chat Request/Response
class ChatRequest(BaseModel):
    message: str
    session_id: str | None = None

class ChatResponse(BaseModel):
    response: str
    session_id: str
    tool_calls: list[ToolCall] | None = None
```

### 5.2 Configuration Schema

```yaml
# config.yaml
llm:
  base_url: "http://localhost:11434/v1"
  model: "qwen3-coder"
  api_key: null

sandbox:
  type: "docker"  # subprocess | docker | e2b
  image: "coddy-sandbox:latest"
  timeout: 30
  memory_limit: "512m"
  cpu_limit: 1.0
  network: "none"

orchestrator:
  max_tool_iterations: 10
  session_timeout: 3600

server:
  host: "0.0.0.0"
  port: 8000
```

### 5.3 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/sessions` | Create new session |
| DELETE | `/sessions/{id}` | Destroy session |
| POST | `/sessions/{id}/chat` | Send message (REST) |
| WS | `/ws/sessions/{id}` | WebSocket chat |
| POST | `/sessions/{id}/upload` | Upload file |
| GET | `/sessions/{id}/files` | List files |
| GET | `/sessions/{id}/files/{path}` | Download file |

---

## 6. Future Enhancements (Post-MVP)

### 6.1 Additional Features
- [ ] Multi-language support (Go, Rust, Ruby)
- [ ] Package installation on demand
- [ ] Git integration (clone repos)
- [ ] Database connections (SQLite, PostgreSQL)
- [ ] Persistent conversation history (database)

### 6.2 Advanced Sandboxing
- [ ] gVisor integration (stronger isolation)
- [ ] Firecracker microVMs
- [ ] GPU support for ML workloads

### 6.3 Enterprise Features
- [ ] Authentication (API keys, OAuth)
- [ ] Rate limiting
- [ ] Usage quotas
- [ ] Audit logging
- [ ] Multi-tenancy

### 6.4 Integrations
- [ ] VS Code extension
- [ ] Slack/Discord bot
- [ ] LangChain/LlamaIndex integration
- [ ] MCP (Model Context Protocol) server

---

## 7. Development Workflow

### 7.1 Directory Structure

```
coddy/
├── src/
│   ├── coddy/
│   │   ├── __init__.py
│   │   ├── __main__.py          # CLI entry point
│   │   ├── config.py            # Configuration management
│   │   ├── api/                 # FastAPI application
│   │   │   ├── __init__.py
│   │   │   ├── server.py
│   │   │   ├── routes.py
│   │   │   └── websocket.py
│   │   ├── orchestrator/
│   │   │   ├── __init__.py
│   │   │   ├── core.py
│   │   │   └── session.py
│   │   ├── llm/
│   │   │   ├── __init__.py
│   │   │   ├── client.py
│   │   │   └── parser.py
│   │   ├── sandbox/
│   │   │   ├── __init__.py
│   │   │   ├── base.py
│   │   │   ├── docker.py
│   │   │   ├── subprocess.py
│   │   │   └── e2b.py
│   │   ├── tools/
│   │   │   ├── __init__.py
│   │   │   ├── definitions.py
│   │   │   ├── registry.py
│   │   │   └── handlers.py
│   │   └── prompts/
│   │       └── system.txt
│   └── coddy_cli/               # CLI implementation
│       ├── __init__.py
│       └── main.py
├── docker/
│   ├── Dockerfile
│   └── entrypoint.sh
├── tests/
│   ├── __init__.py
│   ├── test_sandbox.py
│   ├── test_orchestrator.py
│   └── test_api.py
├── examples/
│   ├── cli_chat.py
│   └── api_client.py
├── pyproject.toml
├── requirements.txt
├── README.md
├── ROADMAP.md
└── AGENTS.md
```

### 7.2 Development Commands

```bash
# Setup
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# Development
ruff check src/              # Linting
ruff format src/             # Formatting
mypy src/                    # Type checking
pytest tests/                # Run tests
pytest tests/ -v --cov       # With coverage

# Running
python -m coddy              # CLI mode
python -m coddy.api          # API server mode

# Docker
docker build -t coddy-sandbox -f docker/Dockerfile .
docker run -d --name coddy-test coddy-sandbox
```

---

## 8. Success Metrics

| Metric | Target |
|--------|--------|
| End-to-end latency | < 5s for simple code execution |
| Sandbox startup time | < 2s (Docker) |
| Test coverage | > 80% |
| Concurrent sessions | > 10 (on 8GB RAM) |
| API uptime | 99.9% |

---

## 9. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Sandbox escape | Low | High | Use gVisor/Firecracker in production |
| Resource exhaustion | Medium | Medium | Memory/CPU limits, rate limiting |
| LLM API failures | Medium | Low | Retry logic, fallback models |
| Complex state bugs | Medium | Medium | Comprehensive tests, session isolation |

---

## 10. Next Steps

1. **Review this roadmap** — Provide feedback on scope and priorities
2. **Set up development environment** — Python 3.12+, Docker, Ollama
3. **Begin Phase 0** — Project structure and foundation
4. **Start with Phase 1** — Build the core engine

---

*Last updated: 2026-04-12*
*Status: Planning Phase*
