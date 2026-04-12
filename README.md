# 🚀 Coddy

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg?logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Give every LLM the power to write, execute, and iterate on code safely.**

Coddy is an open, model-agnostic orchestration layer that enables any LLM (Qwen, Llama, Mistral, GPT-4, etc.) to run Python and Node.js code in isolated sandboxes — just like Claude's coding environment, but self-hostable and written in Go.

---

## ✨ Features

- 🧠 **Model Agnostic** — Works with any OpenAI-compatible LLM
- 🔒 **Secure Sandboxing** — Docker-based isolation with resource limits
- 📝 **Stateful Sessions** — Variables and files persist across interactions
- 🌐 **Dual Interface** — CLI for quick tasks, HTTP API for integrations
- ⚡ **Concurrent** — Built with Go goroutines for high performance
- 📦 **Single Binary** — Easy deployment with static binary
- 🐍 **Python & Node.js** — Built-in support with pre-installed packages

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Python 3.12+ and/or Node.js 20 (for subprocess sandbox)
- Docker (for Docker sandbox)
- An OpenAI-compatible LLM (Ollama recommended for local use)

### Installation

```bash
# Clone the repository
git clone https://github.com/firasmosbehi/coddy.git
cd coddy

# Download dependencies
make deps

# Build binaries
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
SANDBOX_TYPE=subprocess  # Use 'docker' for production
SANDBOX_TIMEOUT=30
```

### Usage

#### CLI Mode

```bash
# Run the CLI
./build/coddy

# Example session
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
│       │  │ (Docker)│
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
│   ├── config/            # Configuration
│   ├── llm/               # LLM client
│   ├── sandbox/           # Sandbox implementations
│   ├── orchestrator/      # Core execution loop
│   ├── tools/             # Tool definitions
│   └── api/               # HTTP handlers
├── pkg/models/             # Shared types
├── docker/                 # Sandbox Dockerfile
└── Makefile               # Build automation
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

# Build Docker sandbox
make docker
```

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
