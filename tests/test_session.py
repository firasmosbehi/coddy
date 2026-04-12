"""Tests for session management."""

import time
from datetime import datetime, timedelta

from coddy.orchestrator.session import Session


class TestSession:
    """Test cases for Session class."""

    def test_session_creation(self) -> None:
        """Test that session is created with required fields."""
        session = Session()
        
        assert session.id is not None
        assert len(session.id) > 0
        assert isinstance(session.id, str)
        assert session.sandbox is None
        assert session.messages == []
        assert isinstance(session.created_at, datetime)
        assert isinstance(session.last_activity, datetime)

    def test_add_message(self) -> None:
        """Test adding messages to session."""
        session = Session()
        
        session.add_message("user", "Hello")
        assert len(session.messages) == 1
        assert session.messages[0]["role"] == "user"
        assert session.messages[0]["content"] == "Hello"

    def test_add_message_with_metadata(self) -> None:
        """Test adding message with additional metadata."""
        session = Session()
        
        session.add_message("assistant", "Hi", tool_calls=[{"id": "123"}])
        assert session.messages[0]["tool_calls"] == [{"id": "123"}]

    def test_touch_updates_activity(self) -> None:
        """Test that touch updates last_activity."""
        session = Session()
        old_activity = session.last_activity
        
        time.sleep(0.01)  # Small delay
        session.touch()
        
        assert session.last_activity > old_activity

    def test_is_expired(self) -> None:
        """Test session expiration detection."""
        session = Session()
        
        # Not expired (just created)
        assert session.is_expired(timeout_seconds=3600) is False
        
        # Manually set last activity to be old
        session.last_activity = datetime.utcnow() - timedelta(seconds=4000)
        assert session.is_expired(timeout_seconds=3600) is True

    def test_is_not_expired(self) -> None:
        """Test that active session is not expired."""
        session = Session()
        session.last_activity = datetime.utcnow() - timedelta(seconds=60)
        
        assert session.is_expired(timeout_seconds=3600) is False

    def test_to_dict(self) -> None:
        """Test session serialization."""
        session = Session()
        session.add_message("user", "Hello")
        session.add_message("assistant", "Hi")
        
        data = session.to_dict()
        
        assert "id" in data
        assert "created_at" in data
        assert "last_activity" in data
        assert data["message_count"] == 2
        assert isinstance(data["created_at"], str)
