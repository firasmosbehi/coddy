"""Tests for configuration module."""

import pytest
from pydantic import ValidationError

from coddy.config import Settings


class TestSettings:
    """Test cases for Settings class."""

    def test_default_values(self) -> None:
        """Test that default values are set correctly."""
        settings = Settings()
        
        assert settings.llm_base_url == "http://localhost:11434/v1"
        assert settings.llm_model == "qwen3-coder"
        assert settings.llm_api_key == ""
        assert settings.sandbox_type == "docker"
        assert settings.sandbox_timeout == 30
        assert settings.sandbox_memory_limit == "512m"
        assert settings.sandbox_cpu_limit == 1.0
        assert settings.sandbox_network == "none"
        assert settings.max_tool_iterations == 10
        assert settings.session_timeout == 3600
        assert settings.host == "0.0.0.0"
        assert settings.port == 8000
        assert settings.log_level == "info"

    def test_is_local_llm_with_local_url(self) -> None:
        """Test is_local_llm with localhost URL."""
        settings = Settings(llm_base_url="http://localhost:11434/v1")
        assert settings.is_local_llm is True

    def test_is_local_llm_with_remote_url(self) -> None:
        """Test is_local_llm with remote URL."""
        settings = Settings(
            llm_base_url="https://api.openai.com/v1",
            llm_api_key="sk-test",
        )
        assert settings.is_local_llm is False

    def test_timeout_validation(self) -> None:
        """Test that timeout validation works."""
        # Valid timeout
        settings = Settings(sandbox_timeout=60)
        assert settings.sandbox_timeout == 60

        # Too high timeout should fail
        with pytest.raises(ValidationError):
            Settings(sandbox_timeout=200)

    def test_sandbox_type_validation(self) -> None:
        """Test sandbox type validation."""
        # Valid types
        Settings(sandbox_type="subprocess")
        Settings(sandbox_type="docker")
        Settings(sandbox_type="e2b")

    def test_port_validation(self) -> None:
        """Test port validation."""
        # Valid ports
        Settings(port=1)
        Settings(port=65535)
        
        # Invalid ports
        with pytest.raises(ValidationError):
            Settings(port=0)
        with pytest.raises(ValidationError):
            Settings(port=70000)
