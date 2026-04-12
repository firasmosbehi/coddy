// Package session provides session management for Coddy.
package session

import (
	"context"
	"log"
	"time"
)

// CleanupWorker periodically cleans up expired sessions.
type CleanupWorker struct {
	store       Store
	interval    time.Duration
	stopChan    chan struct{}
	stoppedChan chan struct{}
}

// NewCleanupWorker creates a new cleanup worker.
func NewCleanupWorker(store Store, interval time.Duration) *CleanupWorker {
	if interval == 0 {
		interval = 5 * time.Minute
	}

	return &CleanupWorker{
		store:       store,
		interval:    interval,
		stopChan:    make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}
}

// Start starts the cleanup worker.
func (w *CleanupWorker) Start() {
	go w.run()
}

// Stop stops the cleanup worker.
func (w *CleanupWorker) Stop() {
	close(w.stopChan)
	<-w.stoppedChan
}

func (w *CleanupWorker) run() {
	defer close(w.stoppedChan)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run cleanup immediately on start
	w.cleanup()

	for {
		select {
		case <-ticker.C:
			w.cleanup()
		case <-w.stopChan:
			return
		}
	}
}

func (w *CleanupWorker) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := w.store.Cleanup(ctx); err != nil {
		log.Printf("Cleanup error: %v", err)
	}
}

// SessionManager combines the store and cleanup worker.
type SessionManager struct {
	store  *MemoryStore
	worker *CleanupWorker
}

// NewSessionManager creates a new session manager.
func NewSessionManager(config *StoreConfig) *SessionManager {
	store := NewMemoryStore(config)
	worker := NewCleanupWorker(store, config.SessionTimeout/2)

	return &SessionManager{
		store:  store,
		worker: worker,
	}
}

// Start starts the session manager.
func (m *SessionManager) Start() {
	m.worker.Start()
}

// Stop stops the session manager.
func (m *SessionManager) Stop() {
	m.worker.Stop()
	m.store.Close()
}

// Store returns the underlying store.
func (m *SessionManager) Store() *MemoryStore {
	return m.store
}

// Stats returns statistics.
func (m *SessionManager) Stats() StoreStats {
	return m.store.Stats()
}
