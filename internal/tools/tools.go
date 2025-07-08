package tools

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/sirupsen/logrus"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/sqlpp"
)

const (
	// MaxLogOutputLength is the maximum length of output to include in tool logs
	MaxLogOutputLength = 500
)

// truncateForLogging truncates output for logging purposes to avoid overwhelming logs
func truncateForLogging(output string) string {
	if len(output) <= MaxLogOutputLength {
		return output
	}
	return output[:MaxLogOutputLength] + "... (truncated)"
}

// ToolHandler handles MCP tool execution
type ToolHandler struct {
	executor sqlpp.ExecutorInterface
	logger   *logrus.Logger
}

// NewToolHandler creates a new tool handler
func NewToolHandler(executor sqlpp.ExecutorInterface, logger *logrus.Logger) *ToolHandler {
	return &ToolHandler{
		executor: executor,
		logger:   logger,
	}
}

// Tool represents a simplified tool definition
type Tool struct {
	Name        string
	Description string
	InputSchema *jsonschema.Schema
}

// GetTools returns all available MCP tools
func (h *ToolHandler) GetTools() []Tool {
	return []Tool{
		h.createSchemaAllTool(),
		h.createSchemaTablesTool(),
		h.createSchemaViewsTool(),
		h.createSchemaProceduresTool(),
		h.createSchemaFunctionsTool(),
		h.createListConnectionsTool(),
		h.createExecuteSQLTool(),
		h.createDriversTool(),
	}
}

// ExecuteTool executes a tool with the given name and arguments
func (h *ToolHandler) ExecuteTool(name string, arguments map[string]interface{}) (string, error) {
	h.logger.WithFields(logrus.Fields{
		"tool":      name,
		"arguments": arguments,
	}).Debug("Executing tool")

	var result string
	var err error

	switch name {
	case "list_schema_all":
		result, err = h.executeSchemaCommand("all", arguments)
	case "list_schema_tables":
		result, err = h.executeSchemaCommand("tables", arguments)
	case "list_schema_views":
		result, err = h.executeSchemaCommand("views", arguments)
	case "list_schema_procedures":
		result, err = h.executeSchemaCommand("procedures", arguments)
	case "list_schema_functions":
		result, err = h.executeSchemaCommand("functions", arguments)
	case "list_connections":
		result, err = h.executeListConnections(arguments)
	case "execute_sql_command":
		result, err = h.executeSQL(arguments)
	case "list_drivers":
		result, err = h.executeDrivers(arguments)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}

	// Log tool execution result
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"tool":  name,
			"error": err.Error(),
		}).Error("Tool execution failed")
	} else {
		h.logger.WithFields(logrus.Fields{
			"tool":        name,
			"result_size": len(result),
		}).Debug("Tool execution succeeded")

		// Log truncated result at TRACE level for detailed debugging
		if h.logger.Level <= logrus.TraceLevel {
			h.logger.WithFields(logrus.Fields{
				"tool":           name,
				"result_preview": truncateForLogging(result),
			}).Trace("Tool execution result preview")
		}
	}

	return result, err
}

// Schema command tools
func (h *ToolHandler) createSchemaAllTool() Tool {
	schema := h.createSchemaToolSchema()
	return Tool{
		Name:        "list_schema_all",
		Description: "Retrieve all schema information (tables, views, procedures, functions) from the database",
		InputSchema: &schema,
	}
}

func (h *ToolHandler) createSchemaTablesTool() Tool {
	schema := h.createSchemaToolSchema()
	return Tool{
		Name:        "list_schema_tables",
		Description: "Retrieve table schema information from the database",
		InputSchema: &schema,
	}
}

func (h *ToolHandler) createSchemaViewsTool() Tool {
	schema := h.createSchemaToolSchema()
	return Tool{
		Name:        "list_schema_views",
		Description: "Retrieve view schema information from the database",
		InputSchema: &schema,
	}
}

func (h *ToolHandler) createSchemaProceduresTool() Tool {
	schema := h.createSchemaToolSchema()
	return Tool{
		Name:        "list_schema_procedures",
		Description: "Retrieve stored procedure schema information from the database",
		InputSchema: &schema,
	}
}

func (h *ToolHandler) createSchemaFunctionsTool() Tool {
	schema := h.createSchemaToolSchema()
	return Tool{
		Name:        "list_schema_functions",
		Description: "Retrieve function schema information from the database",
		InputSchema: &schema,
	}
}

func (h *ToolHandler) createSchemaToolSchema() jsonschema.Schema {
	return jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"connection": {
				Type:        "string",
				Description: "Database connection name to use",
			},
			"filter": {
				Type:        "string",
				Description: "Filter pattern to apply to results (optional)",
			},
			"output": {
				Type:        "string",
				Description: "Output format (json, table, csv, etc.)",
			},
		},
		Required: []string{"connection"},
	}
}

// List connections tool
func (h *ToolHandler) createListConnectionsTool() Tool {
	schema := jsonschema.Schema{
		Type:       "object",
		Properties: map[string]*jsonschema.Schema{},
	}
	return Tool{
		Name:        "list_connections",
		Description: "List all available database connections",
		InputSchema: &schema,
	}
}

// Execute SQL tool
func (h *ToolHandler) createExecuteSQLTool() Tool {
	schema := jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"connection": {
				Type:        "string",
				Description: "Database connection name to use",
			},
			"command": {
				Type:        "string",
				Description: "SQL command(s) to execute. Multiple commands can be separated by GO statements",
			},
			"output": {
				Type:        "string",
				Description: "Output format (json, table, csv, etc.)",
			},
		},
		Required: []string{"connection", "command"},
	}
	return Tool{
		Name:        "execute_sql_command",
		Description: "Execute SQL commands against the database",
		InputSchema: &schema,
	}
}

// Drivers tool
func (h *ToolHandler) createDriversTool() Tool {
	schema := jsonschema.Schema{
		Type:       "object",
		Properties: map[string]*jsonschema.Schema{},
	}
	return Tool{
		Name:        "list_drivers",
		Description: "List all available database drivers",
		InputSchema: &schema,
	}
}

// Tool execution methods
func (h *ToolHandler) executeSchemaCommand(schemaType string, arguments map[string]interface{}) (string, error) {
	connection := h.getStringArg(arguments, "connection", "")
	filter := h.getStringArg(arguments, "filter", "")
	output := h.getStringArg(arguments, "output", "")

	if connection == "" {
		return "", fmt.Errorf("connection parameter is required")
	}

	result, err := h.executor.ExecuteSchemaCommand(schemaType, connection, filter, output)
	if err != nil {
		return "", fmt.Errorf("error executing schema command: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("sqlpp command failed: %s", result.Error)
	}

	return h.formatResult(result.Output), nil
}

func (h *ToolHandler) executeListConnections(arguments map[string]interface{}) (string, error) {
	result, err := h.executor.ListConnections()
	if err != nil {
		return "", fmt.Errorf("error listing connections: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("sqlpp command failed: %s", result.Error)
	}

	return h.formatResult(result.Output), nil
}

func (h *ToolHandler) executeSQL(arguments map[string]interface{}) (string, error) {
	connection := h.getStringArg(arguments, "connection", "")
	command := h.getStringArg(arguments, "command", "")
	output := h.getStringArg(arguments, "output", "")

	if connection == "" {
		return "", fmt.Errorf("connection parameter is required")
	}

	if command == "" {
		return "", fmt.Errorf("command parameter is required")
	}

	result, err := h.executor.ExecuteSQLCommand(connection, command, output)
	if err != nil {
		return "", fmt.Errorf("error executing SQL command: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("sqlpp command failed: %s", result.Error)
	}

	return h.formatResult(result.Output), nil
}

func (h *ToolHandler) executeDrivers(arguments map[string]interface{}) (string, error) {
	result, err := h.executor.ListDrivers()
	if err != nil {
		return "", fmt.Errorf("error listing drivers: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("sqlpp command failed: %s", result.Error)
	}

	return h.formatResult(result.Output), nil
}

// Helper methods
func (h *ToolHandler) getStringArg(arguments map[string]interface{}, key, defaultValue string) string {
	if val, ok := arguments[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (h *ToolHandler) formatResult(output string) string {
	// Try to parse as JSON for better formatting
	var jsonData interface{}
	if err := json.Unmarshal([]byte(output), &jsonData); err == nil {
		// It's valid JSON, format it nicely
		if formatted, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
			return string(formatted)
		}
	}

	// Not JSON or formatting failed, return as plain text
	return output
}
