# MCP Protocol Testing Guide

## Understanding the Error

The error you encountered:
```
{"jsonrpc":"2.0","id":1751990095,"error":{"code":0,"message":"method \"tools/call\" is invalid during session initialization"}}
```

This happens because MCP requires a proper handshake sequence before tools can be called. You cannot directly call tools without first initializing the connection.

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
