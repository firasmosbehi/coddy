"""Tests for sandbox base classes."""

import pytest

from coddy.sandbox.base import ExecutionResult, Sandbox


class TestExecutionResult:
    """Test cases for ExecutionResult dataclass."""

    def test_successful_execution(self) -> None:
        """Test successful execution result."""
        result = ExecutionResult(
            stdout="hello\n",
            stderr="",
            exit_code=0,
            output_files=[],
            execution_time_ms=100,
            timed_out=False,
        )
        
        assert result.success is True
        assert result.stdout == "hello\n"
        assert result.exit_code == 0

    def test_failed_execution(self) -> None:
        """Test failed execution result."""
        result = ExecutionResult(
            stdout="",
            stderr="Error: division by zero",
            exit_code=1,
        )
        
        assert result.success is False
        assert result.stderr == "Error: division by zero"

    def test_timed_out_execution(self) -> None:
        """Test timed out execution result."""
        result = ExecutionResult(
            stdout="",
            stderr="",
            exit_code=-1,
            timed_out=True,
        )
        
        assert result.success is False
        assert result.timed_out is True

    def test_to_dict(self) -> None:
        """Test conversion to dictionary."""
        result = ExecutionResult(
            stdout="output",
            stderr="error",
            exit_code=0,
            output_files=["file.txt"],
            execution_time_ms=50,
            timed_out=False,
        )
        
        d = result.to_dict()
        assert d["stdout"] == "output"
        assert d["stderr"] == "error"
        assert d["exit_code"] == 0
        assert d["output_files"] == ["file.txt"]
        assert d["execution_time_ms"] == 50
        assert d["timed_out"] is False
        assert d["success"] is True


class TestSandboxAbstract:
    """Test cases for abstract Sandbox class."""

    def test_cannot_instantiate_abstract(self) -> None:
        """Test that Sandbox cannot be instantiated directly."""
        with pytest.raises(TypeError):
            Sandbox()  # type: ignore

    def test_sandbox_init_parameters(self) -> None:
        """Test that sandbox initialization stores parameters."""
        
        class ConcreteSandbox(Sandbox):
            async def execute(self, code, language, timeout=None):
                pass
            
            async def upload_file(self, local_path, sandbox_path):
                pass
            
            async def download_file(self, sandbox_path):
                return b""
            
            async def list_files(self, path="/home/user"):
                return []
            
            async def reset(self):
                pass
            
            async def cleanup(self):
                pass
        
        sandbox = ConcreteSandbox(
            timeout=60,
            memory_limit="1g",
            cpu_limit=2.0,
            network="bridge",
        )
        
        assert sandbox.timeout == 60
        assert sandbox.memory_limit == "1g"
        assert sandbox.cpu_limit == 2.0
        assert sandbox.network == "bridge"
