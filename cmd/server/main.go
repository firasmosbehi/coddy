// Coddy API Server - HTTP API for code execution
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/firasmosbehi/coddy/internal/api"
	"github.com/firasmosbehi/coddy/internal/config"
	"github.com/firasmosbehi/coddy/internal/llm"
	"github.com/firasmosbehi/coddy/internal/sandbox"
	"github.com/firasmosbehi/coddy/internal/session"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("🚀 Coddy API Server starting on %s", cfg.Address())

	// Create LLM client
	llmClient := llm.NewClient(cfg.LLMBaseURL, cfg.LLMModel, cfg.LLMAPIKey)

	// Create session manager
	sessionConfig := &session.StoreConfig{
		SessionTimeout: cfg.SessionTimeout,
		SandboxConfig: &sandbox.Config{
			Type:        cfg.SandboxType,
			Image:       cfg.SandboxImage,
			Timeout:     int(cfg.SandboxTimeout.Seconds()),
			MemoryLimit: cfg.SandboxMemory,
			CPULimit:    cfg.SandboxCPU,
			Network:     cfg.SandboxNetwork,
		},
	}

	sessionManager := session.NewSessionManager(sessionConfig)
	sessionManager.Start()
	defer sessionManager.Stop()

	// Create handlers
	handlers := api.NewHandlers(cfg, llmClient, sessionManager)

	// Setup routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", handlers.Health)

	// Stats
	mux.HandleFunc("/stats", handlers.GetStats)

	// Session management
	mux.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListSessions(w, r)
		case http.MethodPost:
			handlers.CreateSession(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/sessions/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetSession(w, r)
		case http.MethodDelete:
			handlers.DeleteSession(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Create server
	server := &http.Server{
		Addr:         cfg.Address(),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on http://%s", cfg.Address())
		log.Printf("Health check: http://%s/health", cfg.Address())
		log.Printf("API endpoints:")
		log.Printf("  POST   /sessions     - Create a new session")
		log.Printf("  GET    /sessions     - List all sessions")
		log.Printf("  GET    /sessions/:id - Get session details")
		log.Printf("  DELETE /sessions/:id - Delete a session")
		log.Printf("  GET    /stats        - Get server statistics")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server stopped")
	return nil
}
