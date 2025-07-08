package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/config"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/server"
)

var (
	configPath string
	logLevel   string
	transport  string
	port       int
	host       string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mcp_sqlpp",
	Short: "MCP server for sqlpp database CLI tool",
	Long: `mcp_sqlpp is a Model Context Protocol (MCP) server that provides
access to the sqlpp database CLI tool through standardized tool interfaces.

It supports both STDIO and HTTP+SSE transports for flexible integration
with MCP clients and AI development tools.`,
	RunE: runServer,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "", "Log level (trace, debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVarP(&transport, "transport", "t", "", "Transport mode (stdio, http)")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 0, "HTTP server port (only for HTTP transport)")
	rootCmd.PersistentFlags().StringVar(&host, "host", "", "HTTP server host (only for HTTP transport)")
}

func runServer(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override config with command line flags
	if logLevel != "" {
		cfg.Log.Level = logLevel
	}
	if transport != "" {
		cfg.Server.Transport = transport
	}
	if port != 0 {
		cfg.Server.Port = port
	}
	if host != "" {
		cfg.Server.Host = host
	}

	// Setup logger
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	logger.SetLevel(level)

	// Set log format
	if cfg.Log.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Log startup information
	logger.WithFields(logrus.Fields{
		"version":    "1.0.0",
		"transport":  cfg.Server.Transport,
		"log_level":  cfg.Log.Level,
		"sqlpp_path": cfg.Sqlpp.ExecutablePath,
	}).Info("Starting mcp_sqlpp MCP server")

	// Create and run server
	srv, err := server.New(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		if err == context.Canceled {
			logger.Info("Server stopped gracefully")
			return nil
		}
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
