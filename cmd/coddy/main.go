// Coddy CLI - Interactive chat with code execution
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/firasmosbehi/coddy/internal/config"
	"github.com/firasmosbehi/coddy/internal/llm"
	"github.com/firasmosbehi/coddy/internal/orchestrator"
	"github.com/firasmosbehi/coddy/internal/sandbox"
	"github.com/firasmosbehi/coddy/internal/session"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("🚀 Coddy - LLM Code Execution Environment")
	fmt.Println("=========================================")
	fmt.Printf("Model: %s\n", cfg.LLMModel)
	fmt.Printf("Sandbox: %s\n", cfg.SandboxType)
	fmt.Printf("Timeout: %v\n", cfg.SandboxTimeout)
	fmt.Println()

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

	// Create a new session
	ctx := context.Background()
	sess, err := sessionManager.Store().Create(ctx, cfg.SandboxType)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	fmt.Printf("Session created: %s\n", sess.ID)
	fmt.Println()

	// Create LLM client
	llmClient := llm.NewClient(cfg.LLMBaseURL, cfg.LLMModel, cfg.LLMAPIKey)

	// Interactive loop
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Type your message (or 'quit' to exit, 'clear' to reset history)")
	fmt.Println()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle special commands
		switch strings.ToLower(input) {
		case "quit", "exit", "q":
			fmt.Println("Goodbye!")
			return nil
		case "clear", "reset":
			// Clear history by creating new session orchestrator
			fmt.Println("History cleared.")
			continue
		case "help", "?":
			printHelp()
			continue
		case "session":
			fmt.Printf("Current session: %s\n", sess.ID)
			continue
		case "stats":
			stats := sessionManager.Stats()
			fmt.Printf("Active sessions: %d\n", stats.TotalSessions)
			continue
		}

		// Create session orchestrator for this message
		orch, err := orchestrator.NewSessionOrchestrator(cfg, llmClient, sessionManager.Store(), sess.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating orchestrator: %v\n", err)
			continue
		}

		// Process message
		fmt.Println()
		response, err := orch.HandleMessage(ctx, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		fmt.Println(response)
		fmt.Println()
	}

	return scanner.Err()
}

func printHelp() {
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help     - Show this help")
	fmt.Println("  clear    - Clear conversation history")
	fmt.Println("  session  - Show current session ID")
	fmt.Println("  stats    - Show session statistics")
	fmt.Println("  quit     - Exit the program")
	fmt.Println()
	fmt.Println("Tips:")
	fmt.Println("- Ask the LLM to write and run code")
	fmt.Println("- Variables persist between code executions")
	fmt.Println("- Files created in the sandbox are accessible")
	fmt.Println()
}
