# sqlpp Integration Guide

This document provides detailed information about how the gosqlpp MCP server integrates with the sqlpp CLI tool.

## Overview

The gosqlpp MCP server acts as a bridge between MCP clients and the sqlpp CLI tool, translating MCP tool calls into appropriate sqlpp commands executed via stdin.

## sqlpp CLI Interface

### Command Execution Method

**Important**: The sqlpp CLI only accepts commands via stdin or from files. It does NOT accept SQL commands as direct command-line arguments.

Correct usage:
```bash
# SQL via stdin
echo "SELECT * FROM users;" | sqlpp --stdin --connection main --output table

# Schema commands via stdin
echo "@schema-tables" | sqlpp --stdin --connection main --output json

# From file
sqlpp script.sql --connection main --output table
```

Incorrect usage (will fail):
```bash
# This does NOT work - no --command flag exists
sqlpp --command "SELECT * FROM users;" --connection main
```

### Supported Commands

#### SQL Statements
Any valid SQL statement supported by your database backend:
- SELECT, INSERT, UPDATE, DELETE
- DDL statements (CREATE, ALTER, DROP)
- Stored procedure calls
- Complex queries with joins, subqueries, etc.

#### Schema Commands
Special commands for database introspection:
- `@drivers` - List available database drivers
- `@schema-tables [filter]` - List tables (optional filter pattern)
- `@schema-views [filter]` - List views
- `@schema-procedures [filter]` - List stored procedures
- `@schema-functions [filter]` - List functions
- `@schema-all [filter]` - Show all schema information

#### Command-Line Flags
- `--stdin` - Read commands from standard input (required for MCP server)
- `--connection, -c` - Specify database connection name
- `--output, -o` - Output format (table, json, yaml, csv)
- `--list-connections, -l` - List available connections

## MCP Server Implementation

### Architecture

```
MCP Tool Call → Tool Handler → sqlpp Executor → sqlpp CLI (stdin) → Database
```

### Execution Flow

1. **MCP Tool Call**: Client calls MCP tool (e.g., `execute-sql-command`)
2. **Parameter Validation**: Server validates required parameters
3. **Command Construction**: Server builds appropriate sqlpp command
4. **stdin Execution**: Server executes sqlpp with `--stdin` flag
5. **Input Streaming**: SQL/schema command sent via stdin pipe
6. **Output Capture**: Server captures stdout/stderr from sqlpp
7. **Response Formatting**: Server formats response for MCP client

### Error Handling

Common errors and their meanings:

#### "Error: unknown flag: --command"
- **Cause**: Attempting to use non-existent `--command` flag
- **Fix**: Use stdin interface instead
- **Code**: Send command via stdin pipe, not command-line argument

#### "failed to open file @schema-tables"
- **Cause**: Schema command treated as filename
- **Fix**: Send schema commands via stdin, not as file arguments
- **Code**: Use `executeStdinCommand("@schema-tables")`

#### "sqlpp: command not found"
- **Cause**: sqlpp not in PATH or incorrect executable_path
- **Fix**: Verify sqlpp installation and configuration
- **Code**: Check `executable_path` in config.yaml

## Configuration

### sqlpp Configuration
The MCP server relies on sqlpp's existing configuration:

```yaml
# .sqlppconfig or sqlpp-config.yaml
connections:
  main:
    driver: sqlite3
    dsn: "./test.db"
  postgres:
    driver: postgres
    dsn: "postgres://user:pass@localhost/db"
```

### MCP Server Configuration
```yaml
# config.yaml
sqlpp:
  executable_path: "sqlpp"  # or absolute path
  timeout: 300
```

## Testing and Validation

### Pre-deployment Testing

1. **Verify sqlpp Installation**:
```bash
sqlpp --version
sqlpp --help | grep -E "(stdin|connection|output)"
```

2. **Test Database Connectivity**:
```bash
sqlpp --list-connections
echo "SELECT 1;" | sqlpp --stdin --connection main
```

3. **Test Schema Commands**:
```bash
echo "@drivers" | sqlpp --stdin
echo "@schema-tables" | sqlpp --stdin --connection main --output json
```

### MCP Server Testing

1. **Start Server**:
```bash
./gosqlpp-mcp-server --transport stdio
```

2. **Test Tools** (via MCP client):
- `list-connections`
- `execute-sql-command` with simple SELECT
- `schema-tables` with main connection
- `drivers`

## Troubleshooting

### Debug Mode
Enable debug logging to see sqlpp command execution:
```bash
./gosqlpp-mcp-server --log-level debug
```

### Manual sqlpp Testing
Always test sqlpp directly before troubleshooting MCP server:
```bash
# Test the exact command the MCP server would run
echo "SELECT * FROM users;" | sqlpp --stdin --connection main --output table
```

### Common Issues

1. **Timeout Errors**: Increase timeout in config.yaml
2. **Permission Errors**: Check sqlpp executable permissions
3. **Connection Errors**: Verify sqlpp database configuration
4. **Schema Command Errors**: Ensure commands are sent via stdin

## Version Compatibility

### sqlpp Requirements
- Must support `--stdin` flag
- Must support schema commands (`@schema-*`)
- Must support `--connection` and `--output` flags

### Verification
```bash
# Check for required features
sqlpp --help | grep -E "(stdin|schema|connection|output)"

# Test schema command support
echo "@drivers" | sqlpp --stdin
```

## Performance Considerations

### Process Overhead
Each MCP tool call spawns a new sqlpp process. For high-frequency usage:
- Consider connection pooling at the database level
- Monitor process creation overhead
- Use appropriate timeout values

### Memory Usage
- sqlpp processes are short-lived
- Memory usage depends on query result size
- Configure appropriate output limits for large datasets

## Security Considerations

### Input Validation
- All SQL commands are passed through to sqlpp
- Rely on database-level security and permissions
- Consider SQL injection risks in user-provided queries

### Process Isolation
- Each sqlpp execution is isolated
- No persistent state between calls
- Temporary files cleaned up automatically

### Credential Management
- Database credentials managed by sqlpp configuration
- No credentials stored in MCP server
- Use sqlpp's secure credential management features
