# 🚀 Coddy

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg?logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Give every LLM the power to write, execute, and iterate on code safely.**

Coddy is an open, model-agnostic orchestration layer that enables any LLM (Qwen, Llama, Mistral, GPT-4, etc.) to run Python and Node.js code in isolated sandboxes — just like Claude's coding environment, but self-hostable and written in Go.

---

## ✨ Features

- 🧠 **Model Agnostic** — Works with any OpenAI-compatible LLM
- 🔒 **Sandboxed Execution** — Subprocess sandbox (dev), Docker support planned
- 📝 **Stateful Sessions** — Variables and files persist across interactions
- 🌐 **Full REST API** — Complete HTTP API with file operations
- 🔌 **WebSocket Support** — Real-time streaming chat
- ⚡ **Concurrent** — Built with Go goroutines for high performance
- 📦 **Single Binary** — Easy deployment with static binary
- 🐍 **Python & Node.js** — Built-in support for both languages
- 🔐 **Middleware** — Logging, CORS, rate limiting, recovery

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Python 3.12+ and/or Node.js 20 (for subprocess sandbox)
- An OpenAI-compatible LLM (Ollama recommended for local use)

### Installation

```bash
# Clone the repository
git clone https://github.com/firasmosbehi/coddy.git
cd coddy

# Download dependencies and build
make deps
make build

# Run CLI
./build/coddy

# Or run API server
./build/coddy-server
```

### Configuration

Create a `.env` file:

```env
# LLM Configuration
LLM_BASE_URL=http://localhost:11434/v1
LLM_MODEL=qwen3-coder

# Sandbox Configuration
SANDBOX_TYPE=subprocess
SANDBOX_TIMEOUT=30
```

---

## 📚 API Reference

### Health & Status

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/stats` | Server statistics |

### Sessions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/sessions` | Create new session |
| GET | `/sessions` | List all sessions |
| GET | `/sessions/:id` | Get session details |
| DELETE | `/sessions/:id` | Delete session |

### Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sessions/:id/messages` | Get chat history |
| DELETE | `/sessions/:id/messages` | Clear chat history |

### File Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/sessions/:id/upload` | Upload file |
| GET | `/sessions/:id/files` | List files |
| GET | `/sessions/:id/files/:path` | Download file |

### WebSocket

| Protocol | Endpoint | Description |
|----------|----------|-------------|
| WS | `/ws/sessions/:id` | Real-time streaming chat |

### Example Usage

```bash
# Create a session
SESSION_ID=$(curl -s -X POST http://localhost:8000/sessions | jq -r '.id')

# Upload a file
curl -X POST -F "file=@data.csv" http://localhost:8000/sessions/$SESSION_ID/upload

# List files
curl http://localhost:8000/sessions/$SESSION_ID/files

# Download a file
curl http://localhost:8000/sessions/$SESSION_ID/files/data.csv -o data.csv

# Get messages
curl http://localhost:8000/sessions/$SESSION_ID/messages

# Delete session
curl -X DELETE http://localhost:8000/sessions/$SESSION_ID
```

### WebSocket Example (JavaScript)

```javascript
const ws = new WebSocket('ws://localhost:8000/ws/sessions/:id');

ws.onopen = () => {
  ws.send(JSON.stringify({
    type: "chat",
    payload: JSON.stringify({message: "Calculate 2+2"})
  }));
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.type === "chunk") {
    console.log(msg.content);
  } else if (msg.type === "done") {
    console.log("Complete!");
  }
};
```

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│           API Server (HTTP/WebSocket)    │
├─────────────────────────────────────────┤
│  Middleware: Logging | CORS | Recovery   │
├─────────────────────────────────────────┤
│  Handlers: Sessions | Files | Messages   │
├─────────────────────────────────────────┤
│         Session Manager                  │
├─────────────────────────────────────────┤
│  Orchestrator  ←→  LLM Client            │
├─────────────────────────────────────────┤
│         Sandbox (Subprocess/Docker)      │
└─────────────────────────────────────────┘
```

---

## 📁 Project Structure

```
coddy/
├── cmd/
│   ├── coddy/             # CLI application
│   └── server/            # API server
├── internal/
│   ├── api/               # HTTP handlers, middleware, WebSocket
│   ├── config/            # Configuration
│   ├── llm/               # LLM client
│   ├── sandbox/           # Sandbox implementations
│   ├── orchestrator/      # Core execution loop
│   ├── session/           # Session management
│   └── tools/             # Tool definitions
├── pkg/models/            # Public/shared types
├── docker/                # Sandbox Dockerfile
├── Makefile               # Build automation
└── go.mod                 # Go module
```

---

## 🛠️ Development

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Build binaries
make build

# Run API server
make run-server
```

---

## 🗺️ Roadmap

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 0 | ✅ | Foundation, structure, config |
| Phase 1 | ✅ | Core engine, CLI, subprocess sandbox |
| Phase 2 | 🚧 | Docker sandbox (reference ready) |
| Phase 3 | ✅ | Stateful sessions |
| Phase 4 | ✅ | **Full Web API with WebSocket** |
| Phase 5 | ⏳ | Hardening, monitoring, production |

---

## ⚠️ Security Warning

**The subprocess sandbox runs code directly on your machine.** It provides NO security isolation and should:
- NEVER be used in production
- NEVER be used with untrusted code
- ONLY be used for local development

Use the Docker sandbox (when implemented) for any production deployment.

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## 📄 License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.

---

## 🙏 Acknowledgments

- Inspired by Claude's code execution environment
- Built with Go ❤️

---

<p align="center">
  <b>⭐ Star this repo if you find it useful!</b>
</p>
