// Package session provides session management for Coddy.
package session

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/firasmosbehi/coddy/internal/sandbox"
	"github.com/firasmosbehi/coddy/pkg/models"
)

var (
	// ErrSessionNotFound is returned when a session doesn't exist.
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired is returned when a session has expired.
	ErrSessionExpired = errors.New("session expired")
)

// Store defines the interface for session storage.
type Store interface {
	// Create creates a new session.
	Create(ctx context.Context, sandboxType string) (*models.Session, error)

	// Get retrieves a session by ID.
	Get(ctx context.Context, id string) (*models.Session, error)

	// GetWithSandbox retrieves a session with its sandbox.
	GetWithSandbox(ctx context.Context, id string) (*SessionWithSandbox, error)

	// Update updates a session.
	Update(ctx context.Context, session *models.Session) error

	// Delete deletes a session.
	Delete(ctx context.Context, id string) error

	// List returns all active session IDs.
	List(ctx context.Context) ([]string, error)

	// Cleanup removes expired sessions.
	Cleanup(ctx context.Context) error
}

// SessionWithSandbox combines a session with its sandbox instance.
type SessionWithSandbox struct {
	Session *models.Session
	Sandbox sandbox.Sandbox
}

// MemoryStore is an in-memory implementation of Store.
type MemoryStore struct {
	mu        sync.RWMutex
	sessions  map[string]*sessionEntry
	sandboxes map[string]sandbox.Sandbox
	config    *StoreConfig
}

// sessionEntry wraps a session with internal metadata.
type sessionEntry struct {
	session   *models.Session
	createdAt time.Time
}

// StoreConfig holds configuration for the session store.
type StoreConfig struct {
	// SessionTimeout is the duration after which a session expires.
	SessionTimeout time.Duration

	// SandboxConfig is the configuration for creating sandboxes.
	SandboxConfig *sandbox.Config
}

// NewMemoryStore creates a new in-memory session store.
func NewMemoryStore(config *StoreConfig) *MemoryStore {
	if config.SessionTimeout == 0 {
		config.SessionTimeout = 1 * time.Hour
	}

	return &MemoryStore{
		sessions:  make(map[string]*sessionEntry),
		sandboxes: make(map[string]sandbox.Sandbox),
		config:    config,
	}
}

// Create creates a new session with a sandbox.
func (s *MemoryStore) Create(ctx context.Context, sandboxType string) (*models.Session, error) {
	// Generate unique ID
	id := generateSessionID()

	// Create sandbox
	sbConfig := s.config.SandboxConfig
	if sbConfig == nil {
		sbConfig = &sandbox.Config{
			Type: sandboxType,
		}
	}
	sbConfig.Type = sandboxType

	sb, err := sandbox.Factory(sbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create sandbox: %w", err)
	}

	// Create session
	session := models.NewSession(id, sandboxType)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[id] = &sessionEntry{
		session:   session,
		createdAt: time.Now().UTC(),
	}
	s.sandboxes[id] = sb

	return session, nil
}

// Get retrieves a session by ID.
func (s *MemoryStore) Get(ctx context.Context, id string) (*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.sessions[id]
	if !exists {
		return nil, ErrSessionNotFound
	}

	// Check expiration
	if s.isExpired(entry) {
		return nil, ErrSessionExpired
	}

	// Update last activity
	entry.session.Touch()

	return entry.session, nil
}

// GetWithSandbox retrieves a session with its sandbox.
func (s *MemoryStore) GetWithSandbox(ctx context.Context, id string) (*SessionWithSandbox, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.sessions[id]
	if !exists {
		return nil, ErrSessionNotFound
	}

	// Check expiration
	if s.isExpired(entry) {
		// Clean up expired session
		if sb, ok := s.sandboxes[id]; ok {
			sb.Close()
			delete(s.sandboxes, id)
		}
		delete(s.sessions, id)
		return nil, ErrSessionExpired
	}

	// Update last activity
	entry.session.Touch()

	sb, exists := s.sandboxes[id]
	if !exists {
		return nil, fmt.Errorf("sandbox not found for session %s", id)
	}

	return &SessionWithSandbox{
		Session: entry.session,
		Sandbox: sb,
	}, nil
}

// Update updates a session.
func (s *MemoryStore) Update(ctx context.Context, session *models.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.sessions[session.ID]
	if !exists {
		return ErrSessionNotFound
	}

	entry.session = session
	entry.session.Touch()

	return nil
}

// Delete deletes a session and its sandbox.
func (s *MemoryStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close and remove sandbox
	if sb, exists := s.sandboxes[id]; exists {
		sb.Close()
		delete(s.sandboxes, id)
	}

	delete(s.sessions, id)

	return nil
}

// List returns all active session IDs.
func (s *MemoryStore) List(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ids []string
	now := time.Now().UTC()

	for id, entry := range s.sessions {
		// Skip expired sessions
		if now.Sub(entry.session.LastActivity) > s.config.SessionTimeout {
			continue
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// Cleanup removes expired sessions.
func (s *MemoryStore) Cleanup(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	var expired []string

	for id, entry := range s.sessions {
		if now.Sub(entry.session.LastActivity) > s.config.SessionTimeout {
			expired = append(expired, id)
		}
	}

	// Clean up expired sessions
	for _, id := range expired {
		if sb, exists := s.sandboxes[id]; exists {
			sb.Close()
			delete(s.sandboxes, id)
		}
		delete(s.sessions, id)
	}

	return nil
}

// isExpired checks if a session entry has expired.
func (s *MemoryStore) isExpired(entry *sessionEntry) bool {
	return time.Since(entry.session.LastActivity) > s.config.SessionTimeout
}

// generateSessionID generates a unique session ID.
func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

// Close closes all sandboxes and cleans up.
func (s *MemoryStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, sb := range s.sandboxes {
		sb.Close()
		delete(s.sandboxes, id)
	}

	s.sessions = make(map[string]*sessionEntry)

	return nil
}

// Stats returns statistics about the store.
func (s *MemoryStore) Stats() StoreStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return StoreStats{
		TotalSessions:  len(s.sessions),
		TotalSandboxes: len(s.sandboxes),
	}
}

// StoreStats holds statistics about the session store.
type StoreStats struct {
	TotalSessions  int
	TotalSandboxes int
}
