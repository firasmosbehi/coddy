package session

import (
	"context"
	"testing"
	"time"

	"github.com/firasmosbehi/coddy/internal/sandbox"
)

func TestMemoryStore_Create(t *testing.T) {
	config := &StoreConfig{
		SessionTimeout: 1 * time.Hour,
		SandboxConfig: &sandbox.Config{
			Type:    "subprocess",
			Timeout: 30,
		},
	}

	store := NewMemoryStore(config)
	defer store.Close()

	ctx := context.Background()
	session, err := store.Create(ctx, "subprocess")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	if session.ID == "" {
		t.Error("session ID should not be empty")
	}

	if session.SandboxType != "subprocess" {
		t.Errorf("expected sandbox type 'subprocess', got %s", session.SandboxType)
	}
}

func TestMemoryStore_Get(t *testing.T) {
	config := &StoreConfig{
		SessionTimeout: 1 * time.Hour,
		SandboxConfig: &sandbox.Config{
			Type:    "subprocess",
			Timeout: 30,
		},
	}

	store := NewMemoryStore(config)
	defer store.Close()

	ctx := context.Background()

	session, _ := store.Create(ctx, "subprocess")

	retrieved, err := store.Get(ctx, session.ID)
	if err != nil {
		t.Errorf("failed to get session: %v", err)
	}

	if retrieved.ID != session.ID {
		t.Errorf("expected session ID %s, got %s", session.ID, retrieved.ID)
	}

	_, err = store.Get(ctx, "non-existent")
	if err != ErrSessionNotFound {
		t.Errorf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestMemoryStore_GetExpired(t *testing.T) {
	config := &StoreConfig{
		SessionTimeout: 100 * time.Millisecond,
		SandboxConfig: &sandbox.Config{
			Type:    "subprocess",
			Timeout: 30,
		},
	}

	store := NewMemoryStore(config)
	defer store.Close()

	ctx := context.Background()

	session, _ := store.Create(ctx, "subprocess")

	time.Sleep(200 * time.Millisecond)

	_, err := store.Get(ctx, session.ID)
	if err != ErrSessionExpired {
		t.Errorf("expected ErrSessionExpired, got %v", err)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	config := &StoreConfig{
		SessionTimeout: 1 * time.Hour,
		SandboxConfig: &sandbox.Config{
			Type:    "subprocess",
			Timeout: 30,
		},
	}

	store := NewMemoryStore(config)
	defer store.Close()

	ctx := context.Background()

	session, _ := store.Create(ctx, "subprocess")

	err := store.Delete(ctx, session.ID)
	if err != nil {
		t.Errorf("failed to delete session: %v", err)
	}

	_, err = store.Get(ctx, session.ID)
	if err != ErrSessionNotFound {
		t.Errorf("expected ErrSessionNotFound, got %v", err)
	}
}

func TestMemoryStore_List(t *testing.T) {
	config := &StoreConfig{
		SessionTimeout: 1 * time.Hour,
		SandboxConfig: &sandbox.Config{
			Type:    "subprocess",
			Timeout: 30,
		},
	}

	store := NewMemoryStore(config)
	defer store.Close()

	ctx := context.Background()

	session1, _ := store.Create(ctx, "subprocess")
	session2, _ := store.Create(ctx, "subprocess")

	ids, err := store.List(ctx)
	if err != nil {
		t.Errorf("failed to list sessions: %v", err)
	}

	if len(ids) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(ids))
	}

	found := make(map[string]bool)
	for _, id := range ids {
		found[id] = true
	}

	if !found[session1.ID] {
		t.Error("session1 not found in list")
	}

	if !found[session2.ID] {
		t.Error("session2 not found in list")
	}
}

func TestMemoryStore_Cleanup(t *testing.T) {
	config := &StoreConfig{
		SessionTimeout: 100 * time.Millisecond,
		SandboxConfig: &sandbox.Config{
			Type:    "subprocess",
			Timeout: 30,
		},
	}

	store := NewMemoryStore(config)
	defer store.Close()

	ctx := context.Background()

	session1, _ := store.Create(ctx, "subprocess")
	store.Create(ctx, "subprocess")

	time.Sleep(200 * time.Millisecond)

	session3, _ := store.Create(ctx, "subprocess")

	err := store.Cleanup(ctx)
	if err != nil {
		t.Errorf("cleanup failed: %v", err)
	}

	_, err = store.Get(ctx, session1.ID)
	if err != ErrSessionNotFound {
		t.Error("expired session should be removed")
	}

	_, err = store.Get(ctx, session3.ID)
	if err != nil {
		t.Errorf("non-expired session should exist: %v", err)
	}
}

func TestMemoryStore_Stats(t *testing.T) {
	config := &StoreConfig{
		SessionTimeout: 1 * time.Hour,
		SandboxConfig: &sandbox.Config{
			Type:    "subprocess",
			Timeout: 30,
		},
	}

	store := NewMemoryStore(config)
	defer store.Close()

	ctx := context.Background()

	stats := store.Stats()
	if stats.TotalSessions != 0 {
		t.Errorf("expected 0 sessions initially, got %d", stats.TotalSessions)
	}

	store.Create(ctx, "subprocess")
	store.Create(ctx, "subprocess")

	stats = store.Stats()
	if stats.TotalSessions != 2 {
		t.Errorf("expected 2 sessions, got %d", stats.TotalSessions)
	}
}
