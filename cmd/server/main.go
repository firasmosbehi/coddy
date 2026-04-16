// Coddy API Server - HTTP API for code execution
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	// Setup routes with middleware
	mux := http.NewServeMux()

	// Apply middleware chain
	var handler http.Handler = mux
	handler = api.LoggingMiddleware(handler)
	handler = api.RecoveryMiddleware(handler)
	handler = api.CORSMiddleware(handler)

	// Health check (no auth required)
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
		path := r.URL.Path

		// Handle file operations
		if strings.Contains(path, "/files/") {
			switch r.Method {
			case http.MethodGet:
				handlers.DownloadFile(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		if strings.HasSuffix(path, "/files") {
			switch r.Method {
			case http.MethodGet:
				handlers.ListFiles(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		if strings.HasSuffix(path, "/messages") {
			switch r.Method {
			case http.MethodGet:
				handlers.GetMessages(w, r)
			case http.MethodDelete:
				handlers.ClearMessages(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		if strings.HasSuffix(path, "/upload") {
			switch r.Method {
			case http.MethodPost:
				handlers.UploadFile(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		// Regular session operations
		switch r.Method {
		case http.MethodGet:
			handlers.GetSession(w, r)
		case http.MethodDelete:
			handlers.DeleteSession(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// WebSocket endpoint
	mux.HandleFunc("/ws/sessions/", handlers.HandleWebSocket)

	// Create server
	server := &http.Server{
		Addr:         cfg.Address(),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on http://%s", cfg.Address())
		log.Printf("")
		log.Printf("📚 API Documentation:")
		log.Printf("  Health:    GET  /health")
		log.Printf("  Stats:     GET  /stats")
		log.Printf("  Sessions:  POST /sessions              - Create session")
		log.Printf("             GET  /sessions              - List sessions")
		log.Printf("             GET  /sessions/:id          - Get session")
		log.Printf("             DEL  /sessions/:id          - Delete session")
		log.Printf("  Files:     POST /sessions/:id/upload   - Upload file")
		log.Printf("             GET  /sessions/:id/files    - List files")
		log.Printf("             GET  /sessions/:id/files/*  - Download file")
		log.Printf("  WebSocket: /ws/sessions/:id            - Real-time chat")
		log.Printf("")

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
