"""Configuration management for Coddy.

Uses pydantic-settings for environment variable loading.
"""

from typing import Literal

from pydantic import Field, field_validator
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings loaded from environment variables.
    
    Attributes:
        llm_base_url: Base URL for the LLM API (OpenAI-compatible).
        llm_model: Name of the LLM model to use.
        llm_api_key: API key for the LLM provider (optional for local models).
        sandbox_type: Type of sandbox to use (subprocess, docker, or e2b).
        sandbox_image: Docker image for sandbox containers.
        sandbox_timeout: Default execution timeout in seconds.
        sandbox_memory_limit: Memory limit for sandbox containers.
        sandbox_cpu_limit: CPU limit for sandbox containers.
        sandbox_network: Network mode for sandbox containers.
        max_tool_iterations: Maximum consecutive tool calls per message.
        session_timeout: Session inactivity timeout in seconds.
        host: Host to bind the API server to.
        port: Port to run the API server on.
        log_level: Logging level.
        e2b_api_key: E2B API key (optional).
    """

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
    )

    # LLM Configuration
    llm_base_url: str = Field(
        default="http://localhost:11434/v1",
        description="Base URL for the LLM API",
    )
    llm_model: str = Field(
        default="qwen3-coder",
        description="Name of the LLM model to use",
    )
    llm_api_key: str = Field(
        default="",
        description="API key for the LLM provider",
    )

    # Sandbox Configuration
    sandbox_type: Literal["subprocess", "docker", "e2b"] = Field(
        default="docker",
        description="Type of sandbox to use",
    )
    sandbox_image: str = Field(
        default="coddy-sandbox:latest",
        description="Docker image for sandbox containers",
    )
    sandbox_timeout: int = Field(
        default=30,
        ge=1,
        le=120,
        description="Default execution timeout in seconds",
    )
    sandbox_memory_limit: str = Field(
        default="512m",
        description="Memory limit for sandbox containers",
    )
    sandbox_cpu_limit: float = Field(
        default=1.0,
        gt=0,
        description="CPU limit for sandbox containers",
    )
    sandbox_network: Literal["none", "bridge"] = Field(
        default="none",
        description="Network mode for sandbox containers",
    )

    # Orchestrator Configuration
    max_tool_iterations: int = Field(
        default=10,
        ge=1,
        le=50,
        description="Maximum consecutive tool calls per message",
    )
    session_timeout: int = Field(
        default=3600,
        ge=60,
        description="Session inactivity timeout in seconds",
    )

    # Server Configuration
    host: str = Field(
        default="0.0.0.0",
        description="Host to bind the API server to",
    )
    port: int = Field(
        default=8000,
        ge=1,
        le=65535,
        description="Port to run the API server on",
    )
    log_level: Literal["debug", "info", "warning", "error"] = Field(
        default="info",
        description="Logging level",
    )

    # E2B Configuration
    e2b_api_key: str = Field(
        default="",
        description="E2B API key (optional)",
    )

    @field_validator("sandbox_timeout")
    @classmethod
    def validate_timeout(cls, v: int) -> int:
        """Ensure timeout is within reasonable bounds."""
        if v > 120:
            raise ValueError("Timeout cannot exceed 120 seconds")
        return v

    @property
    def is_local_llm(self) -> bool:
        """Check if using a local LLM (no API key required)."""
        return not self.llm_api_key or self.llm_base_url.startswith("http://localhost")
