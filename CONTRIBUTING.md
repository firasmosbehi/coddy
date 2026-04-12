# Contributing to Coddy

Thank you for your interest in contributing to Coddy! We welcome contributions from everyone. This document provides guidelines to help you get started.

## 🚀 Quick Start

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/coddy.git`
3. Create a virtual environment: `python -m venv .venv`
4. Activate it: `source .venv/bin/activate` (Windows: `.venv\Scripts\activate`)
5. Install dependencies: `pip install -r requirements-dev.txt`
6. Create a branch: `git checkout -b feature/your-feature-name`

## 📋 Development Setup

### Prerequisites

- Python 3.12+
- Docker
- Git

### Install Development Dependencies

```bash
pip install -r requirements-dev.txt
```

### Pre-commit Hooks

We use pre-commit hooks to ensure code quality:

```bash
pre-commit install
```

## 🏗️ Project Structure

```
coddy/
├── src/coddy/          # Main source code
├── tests/              # Test files
├── docker/             # Docker configurations
├── examples/           # Usage examples
└── docs/               # Documentation
```

## 📝 Coding Standards

### Style Guide

We use:
- **Ruff** for linting and formatting
- **MyPy** for type checking
- **Pytest** for testing

### Running Checks

```bash
# Format code
ruff format src/ tests/

# Lint code
ruff check src/ tests/

# Type check
mypy src/

# Run tests
pytest tests/

# Run tests with coverage
pytest tests/ --cov=src/coddy --cov-report=html
```

### Code Style Guidelines

1. **Type Hints**: Use type hints for all function signatures
2. **Docstrings**: Use Google-style docstrings
3. **Comments**: Explain "why", not "what"
4. **Naming**: Use `snake_case` for functions/variables, `PascalCase` for classes

Example:

```python
from typing import Optional

def process_data(data: str, timeout: int = 30) -> Optional[dict]:
    """Process input data with optional timeout.
    
    Args:
        data: The input data to process.
        timeout: Maximum time to wait in seconds.
        
    Returns:
        Processed data as a dictionary, or None if processing fails.
        
    Raises:
        ValueError: If data is empty or invalid.
    """
    if not data:
        raise ValueError("Data cannot be empty")
    # ... implementation
```

## 🧪 Testing

### Writing Tests

- Place tests in `tests/` directory
- Name test files `test_*.py`
- Use descriptive test function names
- Follow the Arrange-Act-Assert pattern

Example:

```python
def test_sandbox_execute_returns_result():
    # Arrange
    sandbox = SubprocessSandbox()
    code = "print('hello')"
    
    # Act
    result = sandbox.execute(code, "python")
    
    # Assert
    assert result.stdout == "hello\n"
    assert result.exit_code == 0
```

### Test Coverage

Aim for at least 80% code coverage. Run:

```bash
pytest tests/ --cov=src/coddy --cov-report=term-missing
```

## 🐛 Reporting Bugs

When reporting bugs, please include:

1. **Description**: Clear description of the bug
2. **Steps to Reproduce**: Minimal steps to reproduce the issue
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**: Python version, OS, Coddy version
6. **Logs/Tracebacks**: Any relevant error messages

Use the [bug report template](https://github.com/yourusername/coddy/issues/new?template=bug_report.md).

## 💡 Feature Requests

We welcome feature requests! Please:

1. Check if the feature has already been requested
2. Provide a clear use case
3. Explain why the feature would be useful

Use the [feature request template](https://github.com/yourusername/coddy/issues/new?template=feature_request.md).

## 🔀 Pull Request Process

1. **Update Documentation**: Update README.md or other docs if needed
2. **Add Tests**: Include tests for new functionality
3. **Update CHANGELOG**: Add entry to CHANGELOG.md
4. **Ensure CI Passes**: All checks must pass
5. **Request Review**: Tag maintainers for review

### PR Checklist

- [ ] Code follows style guidelines
- [ ] Tests added and passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Commits are descriptive
- [ ] PR description explains the changes

## 🏷️ Commit Message Guidelines

Use conventional commits format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, semicolons, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Build process or auxiliary tool changes

Examples:

```
feat(sandbox): add Docker sandbox implementation

fix(orchestrator): handle timeout edge case

docs(readme): update installation instructions
```

## 🙋 Getting Help

- Join our [Discord server](https://discord.gg/coddy) (coming soon)
- Open a [discussion](https://github.com/yourusername/coddy/discussions)
- Email: maintainers@coddy.dev (coming soon)

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Coddy! 🎉
