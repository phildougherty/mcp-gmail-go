package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/phildougherty/mcp-google-calendar-go/internal/calendar"
	"github.com/sirupsen/logrus"
)

type Server struct {
	calendarClient *calendar.Client
	router         *mux.Router
	httpServer     *http.Server
	tools          *ToolRegistry
}

func NewServer(calendarClient *calendar.Client, port int) *Server {
	s := &Server{
		calendarClient: calendarClient,
		router:         mux.NewRouter(),
	}

	s.tools = NewToolRegistry(calendarClient)
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
	// MCP JSON-RPC 2.0 endpoint
	s.router.HandleFunc("/", s.handleMCPRequest).Methods("POST")
	
	// Health check
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
	
	// CORS middleware
	s.router.Use(corsMiddleware)
}

func (s *Server) Start(ctx context.Context) error {
	// Check authentication
	if !s.calendarClient.IsAuthenticated() {
		return fmt.Errorf("Calendar client not authenticated. Run with -auth flag first")
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

type JSONRPCRequest struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      map[string]interface{} `json:"clientInfo"`
}

type ToolsListParams struct{}

type ToolsCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func (s *Server) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON-RPC request: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var result interface{}

	switch req.Method {
	case "initialize":
		result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
				"resources": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "google-calendar-mcp-server",
				"version": "1.0.0",
			},
		}
	case "notifications/initialized":
		// No response needed for notifications
		w.WriteHeader(http.StatusOK)
		return
	case "tools/list":
		tools := s.tools.ListTools()
		result = map[string]interface{}{
			"tools": tools,
		}
	case "tools/call":
		var params ToolsCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			s.sendJSONRPCError(w, req.ID, -32602, fmt.Sprintf("Invalid params: %v", err))
			return
		}
		// Convert arguments to json.RawMessage
		argsBytes, err := json.Marshal(params.Arguments)
		if err != nil {
			s.sendJSONRPCError(w, req.ID, -32602, fmt.Sprintf("Invalid arguments: %v", err))
			return
		}
		var toolErr error
		result, toolErr = s.tools.CallTool(params.Name, json.RawMessage(argsBytes))
		if toolErr != nil {
			s.sendJSONRPCError(w, req.ID, -32603, fmt.Sprintf("Tool call failed: %v", toolErr))
			return
		}
	default:
		s.sendJSONRPCError(w, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
		return
	}

	response := JSONRPCResponse{
		JsonRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) sendJSONRPCError(w http.ResponseWriter, id interface{}, code int, message string) {
	response := JSONRPCResponse{
		JsonRPC: "2.0",
		ID:      id,
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
	w.WriteHeader(http.StatusOK) // JSON-RPC errors still return 200
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":        "healthy",
		"authenticated": s.calendarClient.IsAuthenticated(),
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