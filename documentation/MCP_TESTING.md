# MCP Protocol Testing Guide

## Understanding the Error

The error you encountered:
```
{"jsonrpc":"2.0","id":1751990095,"error":{"code":0,"message":"method \"tools/call\" is invalid during session initialization"}}
```

This happens because MCP requires a proper handshake sequence before tools can be called. You cannot directly call tools without first initializing the connection.

## Testing Philosophy

### Unit Tests
The project includes comprehensive unit tests that verify:

#### Tool Functionality Tests
- **`list_connections` validation**: Ensures the tool returns at least one connection and properly identifies default connections
- **Connection format verification**: Validates that connections include required fields (`name`, `driver`, `notes`, `is_default`)
- **Empty result handling**: Tests behavior when no connections are configured
- **Output truncation**: Tests the logging truncation function to ensure large outputs are properly handled

#### Integration Tests  
- **Real sqlpp integration**: Tests with mock sqlpp executables that return realistic data
- **End-to-end validation**: Verifies the complete flow from tool call to sqlpp execution to result formatting

### Logging Features

The MCP server includes comprehensive logging with different levels of detail:

#### Log Levels
- **ERROR**: Full error details, stderr output, and failure context
- **DEBUG**: Execution metadata including output sizes and parameters (default)
- **TRACE**: Truncated output previews for detailed debugging (max 500 characters)

#### What Gets Logged

**sqlpp Executor Level:**
- Command arguments and input parameters
- Execution success/failure status
- Output size (always logged)
- Truncated output preview (TRACE level only)
- Error details and stderr (on failure)

**Tools Layer:**
- Tool name and arguments
- Execution result size
- Truncated result preview (TRACE level only)
- Error details (on failure)

#### Security and Performance
- **Output truncation**: Large outputs are truncated to 500 characters to prevent log flooding
- **Sensitive data protection**: Full output only shown at TRACE level, which should be disabled in production
- **Performance optimization**: Default DEBUG level shows only metadata, not content
- **Configurable**: Log level can be set via command line or configuration

#### Testing Logging
Use the test script to see different logging levels:
```bash
# Test with detailed TRACE logging (shows output previews)
./test-logging.sh

# Or manually test different levels
./mcp_sqlpp -t stdio --log-level trace  # Shows truncated output
./mcp_sqlpp -t stdio --log-level debug  # Shows only sizes (default)
./mcp_sqlpp -t stdio --log-level error  # Shows only errors
```

### New Connection Tests

The following tests have been added to ensure `list_connections` works correctly:

#### `TestExecuteTool_ListConnections_WithDefaultConnection`
- Verifies that `list_connections` returns multiple connections
- Ensures at least one connection is marked with `"is_default": true`
- Validates the JSON structure includes all required fields
- Tests realistic connection data with different drivers

#### `TestExecuteTool_ListConnections_EmptyResult`
- Tests behavior when no connections are configured
- Ensures graceful handling of empty connection lists

#### `TestListConnectionsIntegration`
- Integration test with mock sqlpp executable
- Validates real-world output format
- Ensures the executable validation and connection listing work together

## Correct MCP Protocol Sequence

MCP follows a specific initialization sequence:

### 1. Initialize Request
The client must first send an `initialize` request:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": {}
    },
    "clientInfo": {
      "name": "test-client",
      "version": "1.0.0"
    }
  }
}
```

### 2. Initialized Notification
After receiving the initialize response, send an `initialized` notification:

```json
{
  "jsonrpc": "2.0",
  "method": "notifications/initialized"
}
```

### 3. List Tools (Optional)
You can list available tools:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

### 4. Call Tools
Now you can call tools:

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "list_drivers",
    "arguments": {}
  }
}
```

## Available Tools

The mcp_sqlpp server provides the following tools:

- `list_drivers` - List available database drivers
- `list_connections` - List configured database connections
- `list_schema_all` - List all database schema objects
- `list_schema_tables` - List database tables
- `list_schema_views` - List database views
- `list_schema_procedures` - List database procedures
- `list_schema_functions` - List database functions
- `execute_sql_command` - Execute SQL commands

## Testing Scripts

Use the provided test scripts to verify the protocol:

```bash
# Simple test that shows the complete sequence
./test-mcp-simple.sh

# More detailed test with step-by-step output
./test-mcp-protocol.sh
```

## Manual Testing

For manual testing, you can use the following command sequence:

```bash
# Build the server
make build

# Test with proper sequence (pipe multiple JSON messages)
(
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{"tools":{}},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
  echo '{"jsonrpc":"2.0","method":"notifications/initialized"}'
  echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
  echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_drivers","arguments":{}}}'
) | ./mcp_sqlpp -t stdio
```

## Common Issues

1. **Tool names**: Make sure to use the correct tool names (e.g., `list_drivers`, not `mcp_sqlpp___list_drivers`)
2. **Initialization**: Always initialize the connection before calling tools
3. **JSON format**: Ensure proper JSON formatting with correct escaping
4. **Protocol version**: Use the correct protocol version `2024-11-05`

## Integration with MCP Clients

When integrating with MCP clients like Claude Desktop, VS Code extensions, or other tools, they will handle the initialization sequence automatically. The error you saw typically occurs during manual testing or when using incomplete test scripts.
