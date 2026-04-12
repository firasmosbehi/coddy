# 🚀 Coddy

[![Python 3.12+](https://img.shields.io/badge/python-3.12+-blue.svg)](https://www.python.org/downloads/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![FastAPI](https://img.shields.io/badge/FastAPI-009688.svg?logo=fastapi&logoColor=white)](https://fastapi.tiangolo.com/)
[![Docker](https://img.shields.io/badge/Docker-2496ED.svg?logo=docker&logoColor=white)](https://www.docker.com/)

> **Give every LLM the power to write, execute, and iterate on code safely.**

Coddy is an open, model-agnostic orchestration layer that enables any LLM (Qwen, Llama, Mistral, GPT-4, etc.) to run Python and Node.js code in isolated sandboxes — just like Claude's coding environment, but self-hostable and extensible.

---

## ✨ Features

- 🧠 **Model Agnostic** — Works with any OpenAI-compatible LLM (local or cloud)
- 🔒 **Secure Sandboxing** — Docker-based isolation with resource limits
- 📝 **Stateful Sessions** — Variables and files persist across interactions
- 🌐 **Dual Interface** — CLI for quick tasks, HTTP API for integrations
- ⚡ **Streaming Support** — Real-time response streaming via WebSocket
- 📁 **File Operations** — Upload, process, and download files
- 🐍 **Python & Node.js** — Built-in support with pre-installed packages

---

## 🚀 Quick Start

### Prerequisites

- Python 3.12+
- Docker
- An OpenAI-compatible LLM (Ollama recommended for local use)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/coddy.git
cd coddy

# Create virtual environment
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Build the sandbox Docker image
docker build -t coddy-sandbox -f docker/Dockerfile .
```

### Configuration

Create a `.env` file:

```env
# LLM Configuration
LLM_BASE_URL=http://localhost:11434/v1
LLM_MODEL=qwen3-coder
LLM_API_KEY=  # Leave empty for Ollama

# Sandbox Configuration
SANDBOX_TYPE=docker
SANDBOX_TIMEOUT=30
SANDBOX_MEMORY_LIMIT=512m
```

### Usage

#### CLI Mode

```bash
# Start interactive chat
python -m coddy

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
python -m coddy.api

# Create a session
curl -X POST http://localhost:8000/sessions
# {"session_id": "abc-123", ...}

# Chat with the session
curl -X POST http://localhost:8000/sessions/abc-123/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Create a plot of sin(x)"}'
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
├── src/coddy/           # Core library
│   ├── api/             # FastAPI application
│   ├── orchestrator/    # Main execution loop
│   ├── llm/             # LLM client
│   ├── sandbox/         # Sandbox implementations
│   ├── tools/           # Tool definitions
│   └── prompts/         # System prompts
├── docker/              # Sandbox Dockerfile
├── tests/               # Test suite
├── examples/            # Usage examples
└── docs/                # Documentation
```

---

## 🛠️ Development

```bash
# Run tests
pytest tests/

# Run with coverage
pytest tests/ --cov=src/coddy

# Linting
ruff check src/
ruff format src/

# Type checking
mypy src/
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
- Built with [FastAPI](https://fastapi.tiangolo.com/), [Docker](https://www.docker.com/), and ❤️

---

<p align="center">
  <b>⭐ Star this repo if you find it useful!</b>
</p>
