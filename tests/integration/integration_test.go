package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/config"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/server"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/sqlpp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerCreation(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create a mock sqlpp executable
	mockSqlppDir := filepath.Join(tmpDir, "mock-bin")
	err := os.MkdirAll(mockSqlppDir, 0755)
	require.NoError(t, err)

	mockSqlpp := filepath.Join(mockSqlppDir, "sqlpp")
	mockScript := `#!/bin/bash
echo "sqlpp help information"
exit 0
`

	err = os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	// Create test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Transport: "stdio",
			Host:      "localhost",
			Port:      8080,
		},
		Sqlpp: config.SqlppConfig{
			ExecutablePath: mockSqlpp,
			Timeout:        30,
		},
		Log: config.LogConfig{
			Level:  "info",
			Format: "text",
		},
		AWS: config.AWSConfig{
			Region:      "us-east-1",
			Environment: "test",
		},
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create server
	srv, err := server.New(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, srv)
}

func TestSqlppExecutor(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create a mock sqlpp executable
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")
	mockScript := `#!/bin/bash
case "$1" in
    "--help")
        echo "sqlpp help information"
        ;;
    "--list-connections")
        echo '["conn1", "conn2"]'
        ;;
    "--stdin")
        # Read from stdin and respond based on content
        input=$(cat)
        case "$input" in
            "@drivers")
                echo '["mysql", "postgresql"]'
                ;;
            "@schema-tables"*)
                echo '{"tables": ["table1", "table2"]}'
                ;;
            *)
                echo "Unknown stdin command: $input"
                exit 1
                ;;
        esac
        ;;
    *)
        echo "Unknown command: $@"
        exit 1
        ;;
esac
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	logger := logrus.New()
	executor := sqlpp.NewExecutor(mockSqlpp, 30, logger)

	// Test validation
	err = executor.ValidateExecutable()
	require.NoError(t, err)

	// Test list connections
	result, err := executor.ListConnections()
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "conn1")

	// Test list drivers
	result, err = executor.ListDrivers()
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "mysql")

	// Test schema command
	result, err = executor.ExecuteSchemaCommand("tables", "test-conn", "", "json")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "table1")
}

func TestServerStartupShutdown(t *testing.T) {
	// Skip this test in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create a mock sqlpp executable
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")
	mockScript := `#!/bin/bash
echo "sqlpp help information"
exit 0
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	// Create test configuration for HTTP mode
	cfg := &config.Config{
		Server: config.ServerConfig{
			Transport: "http",
			Host:      "localhost",
			Port:      8081, // Use different port to avoid conflicts
		},
		Sqlpp: config.SqlppConfig{
			ExecutablePath: mockSqlpp,
			Timeout:        30,
		},
		Log: config.LogConfig{
			Level:  "info",
			Format: "text",
		},
		AWS: config.AWSConfig{
			Region:      "us-east-1",
			Environment: "test",
		},
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create server
	srv, err := server.New(cfg, logger)
	require.NoError(t, err)

	// Start server in background
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Run(ctx)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context to stop server
	cancel()

	// Wait for server to stop
	select {
	case err := <-errChan:
		// Server should stop with context cancellation
		assert.Equal(t, context.Canceled, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Server did not stop within timeout")
	}
}

func TestListConnectionsIntegration(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create a mock sqlpp executable that returns realistic connection data
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")
	mockScript := `#!/bin/bash
if [[ "$1" == "--help" ]]; then
	echo "sqlpp help information"
	exit 0
elif [[ "$1" == "--list-connections" ]]; then
	echo '[
		{"name": "main", "driver": "sqlite3", "notes": "Local SQLite database for testing", "is_default": true},
		{"name": "mysql", "driver": "mysql", "notes": "MySQL production database", "is_default": false},
		{"name": "postgres", "driver": "postgres", "notes": "PostgreSQL development database", "is_default": false}
	]'
	exit 0
fi
echo "Unknown command: $@"
exit 1
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	// Create executor and test list connections
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise during tests

	executor := sqlpp.NewExecutor(mockSqlpp, 30, logger)

	// Validate executable first
	err = executor.ValidateExecutable()
	require.NoError(t, err)

	// Test list connections
	result, err := executor.ListConnections()
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)

	// Verify we have connection data
	assert.Contains(t, result.Output, "main")
	assert.Contains(t, result.Output, "mysql")
	assert.Contains(t, result.Output, "postgres")

	// Verify at least one connection is marked as default
	assert.Contains(t, result.Output, `"is_default": true`)

	// Verify the output is valid JSON with expected structure
	assert.Contains(t, result.Output, `"name": "main"`)
	assert.Contains(t, result.Output, `"driver": "sqlite3"`)
	assert.Contains(t, result.Output, `"notes": "Local SQLite database for testing"`)
}
