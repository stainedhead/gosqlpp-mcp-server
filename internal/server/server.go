package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/config"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/sqlpp"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/tools"
)

// Server represents the MCP server
type Server struct {
	config      *config.Config
	logger      *logrus.Logger
	executor    *sqlpp.Executor
	toolHandler *tools.ToolHandler
	mcpServer   *mcp.Server
}

// New creates a new MCP server instance
func New(cfg *config.Config, logger *logrus.Logger) (*Server, error) {
	// Create sqlpp executor
	executor := sqlpp.NewExecutor(cfg.Sqlpp.ExecutablePath, cfg.Sqlpp.Timeout, logger)

	// Validate sqlpp executable
	if err := executor.ValidateExecutable(); err != nil {
		return nil, fmt.Errorf("sqlpp validation failed: %w", err)
	}

	// Create tool handler
	toolHandler := tools.NewToolHandler(executor, logger)

	// Create MCP server
	mcpServer := mcp.NewServer("mcp_sqlpp", "1.0.0", &mcp.ServerOptions{})

	// Register tools
	for _, tool := range toolHandler.GetTools() {
		toolName := tool.Name // Capture for closure
		handler := mcp.ToolHandler(func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResult, error) {
			result, err := toolHandler.ExecuteTool(toolName, params.Arguments)
			if err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: err.Error()},
					},
					IsError: true,
				}, nil
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: result},
				},
			}, nil
		})

		serverTool := &mcp.ServerTool{
			Tool: &mcp.Tool{
				Name:        toolName,
				Description: tool.Description,
				InputSchema: tool.InputSchema,
			},
			Handler: handler,
		}
		mcpServer.AddTools(serverTool)
	}

	server := &Server{
		config:      cfg,
		logger:      logger,
		executor:    executor,
		toolHandler: toolHandler,
		mcpServer:   mcpServer,
	}

	return server, nil
}

// Run starts the MCP server
func (s *Server) Run(ctx context.Context) error {
	s.logger.WithFields(logrus.Fields{
		"transport": s.config.Server.Transport,
		"host":      s.config.Server.Host,
		"port":      s.config.Server.Port,
	}).Info("Starting MCP server")

	// Create context that cancels on interrupt signals
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	switch s.config.Server.Transport {
	case "stdio":
		return s.runStdio(ctx)
	case "http":
		return s.runHTTP(ctx)
	default:
		return fmt.Errorf("unsupported transport: %s", s.config.Server.Transport)
	}
}

// runStdio runs the server in STDIO mode
func (s *Server) runStdio(ctx context.Context) error {
	s.logger.Info("Running MCP server in STDIO mode")

	// Create STDIO transport
	transport := mcp.NewStdioTransport()

	// Run server with transport
	return s.mcpServer.Run(ctx, transport)
}

// runHTTP runs the server in HTTP+SSE mode
func (s *Server) runHTTP(ctx context.Context) error {
	s.logger.WithFields(logrus.Fields{
		"host": s.config.Server.Host,
		"port": s.config.Server.Port,
	}).Info("Running MCP server in HTTP+SSE mode")

	// Create HTTP server
	mux := http.NewServeMux()

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add MCP SSE endpoint
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		transport := mcp.NewSSEServerTransport("/mcp", w)
		if err := s.mcpServer.Run(r.Context(), transport); err != nil {
			s.logger.WithError(err).Error("MCP server error")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// Configure HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: mux,
	}

	// Start HTTP server
	errChan := make(chan error, 1)
	go func() {
		s.logger.WithField("addr", httpServer.Addr).Info("HTTP server listening")
		errChan <- httpServer.ListenAndServe()
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down HTTP server")

		// Create shutdown context with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Graceful shutdown
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			s.logger.WithError(err).Error("Error during HTTP server shutdown")
			return err
		}

		return ctx.Err()
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			s.logger.WithError(err).Error("HTTP server error")
			return err
		}
		return nil
	}
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	s.logger.Info("Stopping MCP server")
	// The server will be stopped by context cancellation in Run()
	return nil
}
