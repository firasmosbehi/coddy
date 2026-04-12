"""Abstract base class for sandbox implementations."""

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Literal


@dataclass
class ExecutionResult:
    """Result of code execution in a sandbox.
    
    Attributes:
        stdout: Standard output from the execution.
        stderr: Standard error from the execution.
        exit_code: Exit code of the process (0 for success).
        output_files: List of paths to files created during execution.
        execution_time_ms: Time taken to execute in milliseconds.
        timed_out: Whether the execution timed out.
    """

    stdout: str
    stderr: str
    exit_code: int
    output_files: list[str] = field(default_factory=list)
    execution_time_ms: int = 0
    timed_out: bool = False

    def __post_init__(self) -> None:
        """Validate the execution result."""
        if self.exit_code != 0 and not self.stderr and not self.timed_out:
            # Some executions fail without stderr, ensure we have some indication
            self.stderr = f"Process exited with code {self.exit_code}"

    @property
    def success(self) -> bool:
        """Check if execution was successful."""
        return self.exit_code == 0 and not self.timed_out

    def to_dict(self) -> dict:
        """Convert the result to a dictionary."""
        return {
            "stdout": self.stdout,
            "stderr": self.stderr,
            "exit_code": self.exit_code,
            "output_files": self.output_files,
            "execution_time_ms": self.execution_time_ms,
            "timed_out": self.timed_out,
            "success": self.success,
        }


class Sandbox(ABC):
    """Abstract base class for sandbox implementations.
    
    A sandbox provides isolated execution of code with configurable
    resource limits and security constraints.
    """

    def __init__(
        self,
        timeout: int = 30,
        memory_limit: str = "512m",
        cpu_limit: float = 1.0,
        network: Literal["none", "bridge"] = "none",
    ) -> None:
        """Initialize the sandbox.
        
        Args:
            timeout: Maximum execution time in seconds.
            memory_limit: Memory limit (e.g., "512m", "1g").
            cpu_limit: CPU limit (number of cores).
            network: Network mode ("none" for no network, "bridge" for full access).
        """
        self.timeout = timeout
        self.memory_limit = memory_limit
        self.cpu_limit = cpu_limit
        self.network = network

    @abstractmethod
    async def execute(
        self,
        code: str,
        language: Literal["python", "nodejs"],
        timeout: int | None = None,
    ) -> ExecutionResult:
        """Execute code in the sandbox.
        
        Args:
            code: The code to execute.
            language: The programming language ("python" or "nodejs").
            timeout: Override the default timeout (seconds).
            
        Returns:
            ExecutionResult containing output and metadata.
        """
        raise NotImplementedError

    @abstractmethod
    async def upload_file(self, local_path: str, sandbox_path: str) -> None:
        """Copy a file into the sandbox filesystem.
        
        Args:
            local_path: Path to the local file.
            sandbox_path: Destination path in the sandbox.
        """
        raise NotImplementedError

    @abstractmethod
    async def download_file(self, sandbox_path: str) -> bytes:
        """Retrieve a file from the sandbox filesystem.
        
        Args:
            sandbox_path: Path to the file in the sandbox.
            
        Returns:
            File contents as bytes.
        """
        raise NotImplementedError

    @abstractmethod
    async def list_files(self, path: str = "/home/user") -> list[str]:
        """List files in the sandbox filesystem.
        
        Args:
            path: Directory path to list.
            
        Returns:
            List of file paths.
        """
        raise NotImplementedError

    @abstractmethod
    async def reset(self) -> None:
        """Reset the sandbox to a clean state."""
        raise NotImplementedError

    @abstractmethod
    async def cleanup(self) -> None:
        """Clean up sandbox resources (containers, processes, etc.)."""
        raise NotImplementedError
