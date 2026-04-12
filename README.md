# 🚀 Coddy

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg?logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Give every LLM the power to write, execute, and iterate on code safely.**

Coddy is an open, model-agnostic orchestration layer that enables any LLM (Qwen, Llama, Mistral, GPT-4, etc.) to run Python and Node.js code in isolated sandboxes — just like Claude's coding environment, but self-hostable and written in Go.

---

## ✨ Features

- 🧠 **Model Agnostic** — Works with any OpenAI-compatible LLM
- 🔒 **Sandboxed Execution** — Subprocess sandbox (dev), Docker coming soon
- 📝 **Stateful Sessions** — Variables and files persist across interactions
- 🌐 **Dual Interface** — CLI for quick tasks, HTTP API for integrations
- ⚡ **Concurrent** — Built with Go goroutines for high performance
- 📦 **Single Binary** — Easy deployment with static binary
- 🐍 **Python & Node.js** — Built-in support for both languages

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
```

### Configuration

Create a `.env` file:

```env
# LLM Configuration
LLM_BASE_URL=http://localhost:11434/v1
LLM_MODEL=qwen3-coder

# Sandbox Configuration
SANDBOX_TYPE=subprocess  # Use 'docker' when implemented
SANDBOX_TIMEOUT=30
```

### Usage

#### CLI Mode

```bash
# Run the CLI
./build/coddy

# Example session
🚀 Coddy - LLM Code Execution Environment
=========================================
Model: qwen3-coder
Sandbox: subprocess
Timeout: 30s

Type your message (or 'quit' to exit, 'clear' to reset history)

> What is the 50th Fibonacci number?
[Tool Call: execute_code]
STDOUT: 12586269025
Exit code: 0

The 50th Fibonacci number is 12,586,269,025.
```

#### API Mode

```bash
# Start the API server
./build/coddy-server

# Health check
curl http://localhost:8000/health
```

---

## 🏗️ Architecture

```
User Request
     │
     ▼
┌─────────────────┐
│   Orchestrator   │
│    (Coddy)       │
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌───────┐  ┌─────────┐
│  LLM  │  │ Sandbox │
│       │  │(subprocess)
└───────┘  └─────────┘
```

---

## 📁 Project Structure

```
coddy/
├── cmd/                    # Entry points
│   ├── coddy/             # CLI application
│   └── server/            # API server
├── internal/               # Private packages
│   ├── config/            # Environment-based config
│   ├── llm/               # LLM client
│   ├── sandbox/           # Sandbox implementations
│   │   ├── subprocess.go  # Working implementation
│   │   ├── docker.go      # Stub (see docker_impl.go.reference)
│   │   └── docker_impl.go.reference  # Full Docker implementation
│   ├── orchestrator/      # Core execution loop
│   ├── tools/             # Tool definitions
│   └── api/               # HTTP handlers
├── pkg/models/             # Public/shared types
├── docker/                 # Sandbox Dockerfile
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

# Lint code
make lint

# Build for release
make build
```

---

## 🗺️ Roadmap

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 0 | ✅ | Project structure, data models, configuration |
| Phase 1 | ✅ | Core engine: subprocess sandbox, LLM client, CLI |
| Phase 2 | 🚧 | Docker sandbox (stub implemented, full version pending) |
| Phase 3 | ⏳ | Stateful sessions with session manager |
| Phase 4 | ⏳ | Full HTTP API with WebSocket support |
| Phase 5 | ⏳ | Enhanced tools and hardening |

### Docker Sandbox

The Docker sandbox implementation is stubbed due to complex dependency requirements. A complete reference implementation is available in `internal/sandbox/docker_impl.go.reference`. To enable:

1. Set up Docker SDK dependencies
2. Rename `docker_impl.go.reference` to `docker_impl.go`
3. Update imports and build

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
