package tools

import (
	"testing"

	"github.com/sirupsen/logrus"
	"git	expectedTools := []string{
		"list_schema_all",
		"list_schema_tables", 
		"list_schema_views",
		"list_schema_procedures",
		"list_schema_functions",
		"list_connections",
		"execute_sql_command",
		"list_drivers",
	}ainedhead/gosqlpp-mcp-server/internal/sqlpp"
	"github.com/stainedhead/gosqlpp-mcp-server/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockExecutor is a mock implementation of the sqlpp executor
type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) ExecuteSchemaCommand(schemaType, connection, filter, output string) (*types.SqlppResult, error) {
	args := m.Called(schemaType, connection, filter, output)
	return args.Get(0).(*types.SqlppResult), args.Error(1)
}

func (m *MockExecutor) ExecuteSQLCommand(connection, command, output string) (*types.SqlppResult, error) {
	args := m.Called(connection, command, output)
	return args.Get(0).(*types.SqlppResult), args.Error(1)
}

func (m *MockExecutor) ListConnections() (*types.SqlppResult, error) {
	args := m.Called()
	return args.Get(0).(*types.SqlppResult), args.Error(1)
}

func (m *MockExecutor) ListDrivers() (*types.SqlppResult, error) {
	args := m.Called()
	return args.Get(0).(*types.SqlppResult), args.Error(1)
}

func (m *MockExecutor) ValidateExecutable() error {
	args := m.Called()
	return args.Error(0)
}

// Ensure MockExecutor implements the interface
var _ sqlpp.ExecutorInterface = (*MockExecutor)(nil)

func TestNewToolHandler(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()

	handler := NewToolHandler(mockExecutor, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, mockExecutor, handler.executor)
	assert.Equal(t, logger, handler.logger)
}

func TestGetTools(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	tools := handler.GetTools()

	assert.Len(t, tools, 8)

	toolNames := make([]string, len(tools))
	for i, tool := range tools {
		toolNames[i] = tool.Name
	}

	expectedTools := []string{
		"list_schema_all",
		"list_schema_tables",
		"list_schema_views",
		"list_schema_procedures",
		"list_schema_functions",
		"list_connections",
		"execute_sql_command",
		"list_drivers",
	}

	for _, expected := range expectedTools {
		assert.Contains(t, toolNames, expected)
	}
}

func TestExecuteTool_SchemaCommand_Success(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	expectedResult := &types.SqlppResult{
		Success: true,
		Output:  `{"tables": ["table1", "table2"]}`,
	}

	mockExecutor.On("ExecuteSchemaCommand", "tables", "test-conn", "test*", "json").Return(expectedResult, nil)

	arguments := map[string]interface{}{
		"connection": "test-conn",
		"filter":     "test*",
		"output":     "json",
	}

	result, err := handler.ExecuteTool("list_schema_tables", arguments)
	require.NoError(t, err)
	assert.Contains(t, result, "table1")

	mockExecutor.AssertExpectations(t)
}

func TestExecuteTool_SchemaCommand_MissingConnection(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	arguments := map[string]interface{}{
		"filter": "test*",
		"output": "json",
	}

	result, err := handler.ExecuteTool("list_schema_tables", arguments)
	require.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "connection parameter is required")
}

func TestExecuteTool_ExecuteSQL_Success(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	expectedResult := &types.SqlppResult{
		Success: true,
		Output:  `{"rows": [{"id": 1, "name": "test"}]}`,
	}

	mockExecutor.On("ExecuteSQLCommand", "test-conn", "SELECT * FROM users", "json").Return(expectedResult, nil)

	arguments := map[string]interface{}{
		"connection": "test-conn",
		"command":    "SELECT * FROM users",
		"output":     "json",
	}

	result, err := handler.ExecuteTool("execute_sql_command", arguments)
	require.NoError(t, err)
	assert.Contains(t, result, "test")

	mockExecutor.AssertExpectations(t)
}

func TestExecuteTool_ExecuteSQL_MissingParameters(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	// Missing connection
	arguments := map[string]interface{}{
		"command": "SELECT * FROM users",
		"output":  "json",
	}

	result, err := handler.ExecuteTool("execute_sql_command", arguments)
	require.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "connection parameter is required")

	// Missing command
	arguments = map[string]interface{}{
		"connection": "test-conn",
		"output":     "json",
	}

	result, err = handler.ExecuteTool("execute_sql_command", arguments)
	require.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "command parameter is required")
}

func TestExecuteTool_ListConnections_Success(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	expectedResult := &types.SqlppResult{
		Success: true,
		Output:  `["conn1", "conn2", "conn3"]`,
	}

	mockExecutor.On("ListConnections").Return(expectedResult, nil)

	result, err := handler.ExecuteTool("list_connections", map[string]interface{}{})
	require.NoError(t, err)
	assert.Contains(t, result, "conn1")

	mockExecutor.AssertExpectations(t)
}

func TestExecuteTool_ListDrivers_Success(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	expectedResult := &types.SqlppResult{
		Success: true,
		Output:  `["mysql", "postgresql", "sqlite"]`,
	}

	mockExecutor.On("ListDrivers").Return(expectedResult, nil)

	result, err := handler.ExecuteTool("list_drivers", map[string]interface{}{})
	require.NoError(t, err)
	assert.Contains(t, result, "mysql")

	mockExecutor.AssertExpectations(t)
}

func TestExecuteTool_UnknownTool(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	result, err := handler.ExecuteTool("unknown_tool", map[string]interface{}{})
	require.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "unknown tool")
}

func TestExecuteTool_SqlppFailure(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	expectedResult := &types.SqlppResult{
		Success: false,
		Error:   "Connection failed",
	}

	mockExecutor.On("ListConnections").Return(expectedResult, nil)

	result, err := handler.ExecuteTool("list_connections", map[string]interface{}{})
	require.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "Connection failed")

	mockExecutor.AssertExpectations(t)
}

func TestGetStringArg(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	arguments := map[string]interface{}{
		"string_arg": "test_value",
		"int_arg":    123,
		"nil_arg":    nil,
	}

	// Test existing string argument
	result := handler.getStringArg(arguments, "string_arg", "default")
	assert.Equal(t, "test_value", result)

	// Test non-string argument
	result = handler.getStringArg(arguments, "int_arg", "default")
	assert.Equal(t, "default", result)

	// Test missing argument
	result = handler.getStringArg(arguments, "missing_arg", "default")
	assert.Equal(t, "default", result)

	// Test nil argument
	result = handler.getStringArg(arguments, "nil_arg", "default")
	assert.Equal(t, "default", result)
}

func TestFormatResult(t *testing.T) {
	mockExecutor := &MockExecutor{}
	logger := logrus.New()
	handler := NewToolHandler(mockExecutor, logger)

	// Test with JSON output
	result := handler.formatResult(`{"key": "value"}`)
	assert.Contains(t, result, "\"key\": \"value\"")

	// Test with plain text output
	result = handler.formatResult("Plain text output")
	assert.Equal(t, "Plain text output", result)
}
