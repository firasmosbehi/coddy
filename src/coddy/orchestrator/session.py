"""Session management for stateful conversations."""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Any
from uuid import uuid4

from coddy.sandbox.base import Sandbox


@dataclass
class Session:
    """A conversation session with persistent state.
    
    A session maintains:
    - Conversation history (messages)
    - A sandbox instance (for code execution)
    - Metadata (creation time, last activity)
    
    Attributes:
        id: Unique session identifier.
        sandbox: The sandbox instance for this session.
        messages: Conversation history.
        created_at: Session creation timestamp.
        last_activity: Last activity timestamp.
    """

    id: str = field(default_factory=lambda: str(uuid4()))
    sandbox: Sandbox | None = None
    messages: list[dict[str, Any]] = field(default_factory=list)
    created_at: datetime = field(default_factory=datetime.utcnow)
    last_activity: datetime = field(default_factory=datetime.utcnow)

    def __post_init__(self) -> None:
        """Ensure ID is a string."""
        if not isinstance(self.id, str):
            self.id = str(self.id)

    def add_message(self, role: str, content: str, **kwargs: Any) -> None:
        """Add a message to the conversation history.
        
        Args:
            role: Message role (system, user, assistant, tool).
            content: Message content.
            **kwargs: Additional message fields.
        """
        message: dict[str, Any] = {
            "role": role,
            "content": content,
            **kwargs,
        }
        self.messages.append(message)
        self.touch()

    def touch(self) -> None:
        """Update the last activity timestamp."""
        self.last_activity = datetime.utcnow()

    @property
    def is_expired(self, timeout_seconds: int = 3600) -> bool:
        """Check if the session has expired due to inactivity.
        
        Args:
            timeout_seconds: Inactivity timeout in seconds.
            
        Returns:
            True if the session has expired.
        """
        elapsed = (datetime.utcnow() - self.last_activity).total_seconds()
        return elapsed > timeout_seconds

    def to_dict(self) -> dict[str, Any]:
        """Convert session to a dictionary (for API responses).
        
        Returns:
            Dictionary representation of the session.
        """
        return {
            "id": self.id,
            "created_at": self.created_at.isoformat(),
            "last_activity": self.last_activity.isoformat(),
            "message_count": len(self.messages),
        }
