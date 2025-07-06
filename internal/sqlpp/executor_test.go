package sqlpp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExecutor(t *testing.T) {
	logger := logrus.New()
	executor := NewExecutor("/usr/bin/sqlpp", 300, logger)

	assert.NotNil(t, executor)
	assert.Equal(t, "/usr/bin/sqlpp", executor.executablePath)
	assert.Equal(t, 300, int(executor.timeout.Seconds()))
}

func TestExecuteSchemaCommand(t *testing.T) {
	// Create a mock sqlpp executable
	tmpDir := t.TempDir()
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")

	// Create a simple shell script that echoes the arguments
	mockScript := `#!/bin/bash
echo "Mock sqlpp called with: $@"
echo '{"tables": ["table1", "table2"]}'
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	executor := NewExecutor(mockSqlpp, 30, logger)

	result, err := executor.ExecuteSchemaCommand("tables", "test-conn", "test*", "json")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "Mock sqlpp called with")
	assert.Contains(t, result.Output, "--stdin")
	assert.Contains(t, result.Output, "--connection test-conn")
	assert.Contains(t, result.Output, "--output json")
}

func TestExecuteSQLCommand(t *testing.T) {
	// Create a mock sqlpp executable
	tmpDir := t.TempDir()
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")

	mockScript := `#!/bin/bash
echo "Executing SQL: $@"
echo '{"rows": [{"id": 1, "name": "test"}]}'
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	logger := logrus.New()
	executor := NewExecutor(mockSqlpp, 30, logger)

	result, err := executor.ExecuteSQLCommand("test-conn", "SELECT * FROM users", "json")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "Executing SQL")
	assert.Contains(t, result.Output, "--stdin")
	assert.Contains(t, result.Output, "--connection test-conn")
	assert.Contains(t, result.Output, "--output json")
}

func TestListConnections(t *testing.T) {
	// Create a mock sqlpp executable
	tmpDir := t.TempDir()
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")

	mockScript := `#!/bin/bash
echo "Arguments: $@"
echo "Available connections:"
echo '["conn1", "conn2", "conn3"]'
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	logger := logrus.New()
	executor := NewExecutor(mockSqlpp, 30, logger)

	result, err := executor.ListConnections()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "Available connections")
	assert.Contains(t, result.Output, "--list-connections")
}

func TestListDrivers(t *testing.T) {
	// Create a mock sqlpp executable
	tmpDir := t.TempDir()
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")

	mockScript := `#!/bin/bash
echo "Arguments: $@"
echo "Available drivers:"
echo '["mysql", "postgresql", "sqlite"]'
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	logger := logrus.New()
	executor := NewExecutor(mockSqlpp, 30, logger)

	result, err := executor.ListDrivers()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "Available drivers")
	assert.Contains(t, result.Output, "--stdin")
}

func TestExecuteCommand_Failure(t *testing.T) {
	// Create a mock sqlpp executable that fails
	tmpDir := t.TempDir()
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")

	mockScript := `#!/bin/bash
echo "Error: Connection failed" >&2
exit 1
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	logger := logrus.New()
	executor := NewExecutor(mockSqlpp, 30, logger)

	result, err := executor.ListConnections()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "Connection failed")
}

func TestValidateExecutable_Success(t *testing.T) {
	// Create a mock sqlpp executable
	tmpDir := t.TempDir()
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")

	mockScript := `#!/bin/bash
echo "sqlpp help information"
exit 0
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	logger := logrus.New()
	executor := NewExecutor(mockSqlpp, 30, logger)

	err = executor.ValidateExecutable()
	assert.NoError(t, err)
}

func TestValidateExecutable_Failure(t *testing.T) {
	logger := logrus.New()
	executor := NewExecutor("/nonexistent/sqlpp", 30, logger)

	err := executor.ValidateExecutable()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sqlpp executable not found")
}

func TestExecuteCommand_EmptyArgs(t *testing.T) {
	// Create a mock sqlpp executable
	tmpDir := t.TempDir()
	mockSqlpp := filepath.Join(tmpDir, "mock-sqlpp")

	mockScript := `#!/bin/bash
echo "No arguments provided"
`

	err := os.WriteFile(mockSqlpp, []byte(mockScript), 0755)
	require.NoError(t, err)

	logger := logrus.New()
	executor := NewExecutor(mockSqlpp, 30, logger)

	result, err := executor.executeCommand([]string{})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "No arguments provided")
}
