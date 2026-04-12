// Coddy API Server - HTTP API for code execution
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/firasmosbehi/coddy/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("🚀 Coddy API Server starting on %s\n", cfg.Address())
	fmt.Println("Note: API server is not fully implemented yet. Use the CLI for now.")

	// TODO: Implement full API server in Phase 4
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message": "Coddy API - Coming soon"}`))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	log.Fatal(http.ListenAndServe(cfg.Address(), nil))
}
