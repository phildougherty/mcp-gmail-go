package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/phildougherty/mcp-gmail-go/internal/gmail"
	"github.com/sirupsen/logrus"
)

type Server struct {
	gmailClient *gmail.Client
	router      *mux.Router
	httpServer  *http.Server
	tools       *ToolRegistry
}

func NewServer(gmailClient *gmail.Client, port int) *Server {
	s := &Server{
		gmailClient: gmailClient,
		router:      mux.NewRouter(),
	}

	s.tools = NewToolRegistry(gmailClient)
	s.setupRoutes()
	
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) setupRoutes() {
	// MCP protocol endpoints
	s.router.HandleFunc("/mcp/tools", s.handleListTools).Methods("GET")
	s.router.HandleFunc("/mcp/tools/{name}", s.handleCallTool).Methods("POST")
	
	// SSE endpoint for streaming
	s.router.HandleFunc("/mcp/sse", s.handleSSE).Methods("GET")
	
	// Health check
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
	
	// CORS middleware
	s.router.Use(corsMiddleware)
}

func (s *Server) Start(ctx context.Context) error {
	// Check authentication
	if !s.gmailClient.IsAuthenticated() {
		return fmt.Errorf("Gmail client not authenticated. Run with -auth flag first")
	}

	// Start server in goroutine
	go func() {
		logrus.Infof("Server listening on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Server error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(shutdownCtx)
}

func (s *Server) handleListTools(w http.ResponseWriter, r *http.Request) {
	tools := s.tools.ListTools()
	
	response := map[string]interface{}{
		"tools": tools,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCallTool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	toolName := vars["name"]
	
	var request ToolCallRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	
	result, err := s.tools.CallTool(toolName, request.Arguments)
	if err != nil {
		logrus.Errorf("Tool call failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Error: %v", err),
				},
			},
			"isError": true,
		})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Send initial connection message
	fmt.Fprintf(w, "data: {\"type\":\"connection\",\"status\":\"connected\"}\n\n")
	
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Fprintf(w, "data: {\"type\":\"ping\"}\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":        "healthy",
		"authenticated": s.gmailClient.IsAuthenticated(),
		"timestamp":     time.Now().UTC(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}