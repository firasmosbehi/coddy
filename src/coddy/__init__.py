"""Coddy — LLM Code Execution Environment.

Give every LLM the power to write, execute, and iterate on code safely.
"""

__version__ = "0.1.0"
__author__ = "Coddy Contributors"
__license__ = "MIT"

from coddy.config import Settings
from coddy.sandbox.base import Sandbox, ExecutionResult
from coddy.llm.client import LLMClient, LLMResponse
from coddy.orchestrator.core import Orchestrator

__all__ = [
    "Settings",
    "Sandbox",
    "ExecutionResult",
    "LLMClient",
    "LLMResponse",
    "Orchestrator",
]
