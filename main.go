package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/phildougherty/mcp-gmail-go/internal/config"
	"github.com/phildougherty/mcp-gmail-go/internal/gmail"
	"github.com/phildougherty/mcp-gmail-go/internal/mcp"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		port    = flag.Int("port", 8080, "Server port")
		authCmd = flag.Bool("auth", false, "Run OAuth authentication flow")
		debug   = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	// Configure logging
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Gmail client
	gmailClient, err := gmail.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Gmail client: %v", err)
	}

	// Handle auth command
	if *authCmd {
		if err := gmailClient.Authenticate(); err != nil {
			log.Fatalf("Authentication failed: %v", err)
		}
		fmt.Println("Authentication successful!")
		return
	}

	// Create MCP server
	server := mcp.NewServer(gmailClient, *port)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logrus.Info("Received shutdown signal")
		cancel()
	}()

	// Start server
	logrus.Infof("Starting Gmail MCP server on port %d", *port)
	if err := server.Start(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	logrus.Info("Server shutdown complete")
}