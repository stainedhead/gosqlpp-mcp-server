# mcp_sqlpp

A Model Context Protocol (MCP) server that provides access to the [sqlpp](https://github.com/stainedhead/gosqlpp) database CLI tool through standardized tool interfaces. This enables AI development tools and agents to interact with databases through a secure, controlled interface.

## Features

- **MCP Protocol Compliance**: Full support for the Model Context Protocol specification
- **Dual Transport Support**: Both STDIO and HTTP+SSE transports for flexible integration
- **Database Schema Tools**: Access table, view, procedure, and function schemas
- **SQL Execution**: Execute SQL commands with proper output formatting
- **Connection Management**: List and manage database connections
- **Driver Information**: Query available database drivers
- **Comprehensive Logging**: Multiple log levels with optional file logging and automatic rotation
- **Containerized Deployment**: Docker support with AWS App Runner deployment
- **Production Ready**: Comprehensive logging, health checks, monitoring, and extensive test coverage

## Architecture

The server acts as a bridge between MCP clients and the sqlpp CLI tool:

```
MCP Client → mcp_sqlpp → sqlpp CLI (via stdin) → Database
```

### Key Components

- **MCP Server**: Handles protocol communication and tool registration
- **Tool Handler**: Manages tool execution and parameter validation
- **sqlpp Executor**: Wraps sqlpp CLI execution using stdin interface with proper error handling
- **Configuration System**: Flexible configuration via files, environment variables, and CLI flags

**Important**: All SQL commands and schema queries are sent to sqlpp via stdin using the `--stdin` flag. The sqlpp CLI does not accept SQL commands as direct command-line arguments.

For detailed information about sqlpp integration, see [SQLPP_INTEGRATION.md](documentation/SQLPP_INTEGRATION.md).

## Installation

### Prerequisites

- Go 1.21 or later
- [sqlpp](https://github.com/stainedhead/gosqlpp) CLI tool installed and configured
  - **Required**: sqlpp version that supports `--stdin` flag and schema commands (`@schema-*`)
  - **Verify**: Run `sqlpp --help` to ensure `--stdin` flag is available
  - **Test**: Verify schema commands work: `echo "@drivers" | sqlpp --stdin`
  - **Integration Details**: See [SQLPP_INTEGRATION.md](documentation/SQLPP_INTEGRATION.md) for complete setup and configuration
- Docker (for containerized deployment)
- AWS CLI and CDK (for AWS deployment)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/stainedhead/gosqlpp-mcp-server.git
cd gosqlpp-mcp-server
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o mcp_sqlpp ./cmd/server
```

4. Run with default configuration:
```bash
./mcp_sqlpp
```

### Docker Deployment

1. Build the Docker image:
```bash
docker build -f deployment/docker/Dockerfile -t mcp_sqlpp .
```

2. Run the container:
```bash
docker run -p 8080:8080 mcp_sqlpp --transport http
```

## Configuration

The server supports multiple configuration methods with the following precedence:
1. Command line flags (highest priority)
2. Environment variables
3. Configuration file
4. Default values (lowest priority)

### Configuration File

Create a `config.yaml` file:

```yaml
server:
  transport: "stdio"  # or "http"
  host: "localhost"
  port: 8080

sqlpp:
  executable_path: "sqlpp"
  timeout: 300

log:
  level: "info"
  format: "text"  # or "json"
  file_logging: false  # Enable file logging with automatic rolling dates

aws:
  region: "us-east-1"
  environment: "development"
```

### Environment Variables

All configuration options can be set via environment variables with the `GOSQLPP_MCP_` prefix:

```bash
export GOSQLPP_MCP_SERVER_TRANSPORT=http
export GOSQLPP_MCP_SERVER_PORT=8080
export GOSQLPP_MCP_SQLPP_EXECUTABLE_PATH=/usr/local/bin/sqlpp
export GOSQLPP_MCP_LOG_LEVEL=debug
export GOSQLPP_MCP_LOG_FILE_LOGGING=true
```

### Command Line Flags

```bash
./mcp_sqlpp --help
```

Available options:
- `--transport, -t`: Transport mode (stdio, http)
- `--port, -p`: HTTP server port (for HTTP transport)
- `--host`: HTTP server host (for HTTP transport)
- `--config, -c`: Configuration file path
- `--log-level, -l`: Log level (trace, debug, info, warn, error, fatal, panic)
- `--file-logging, -f`: Enable file logging with automatic rolling dates

## Logging

The MCP server provides comprehensive logging capabilities with multiple levels and optional file logging.

### Log Levels

- **TRACE**: Most verbose, includes truncated output previews of tool results and sqlpp responses
- **DEBUG**: Tool execution details, sqlpp commands, and metadata (default for development)
- **INFO**: General application flow and important events
- **WARN/ERROR/FATAL**: Warnings, errors, and fatal conditions

### Console Logging

By default, logs are written to the console in text format:

```bash
# Set log level via command line
./mcp_sqlpp --log-level debug

# Set log level via environment variable
export GOSQLPP_MCP_LOG_LEVEL=trace
./mcp_sqlpp
```

### File Logging

Enable file logging for persistent debugging and monitoring:

```bash
# Enable via command line flag
./mcp_sqlpp --file-logging --log-level trace

# Enable via configuration file
# Set file_logging: true in config.yaml
```

**File Logging Features:**
- **Automatic File Naming**: `logs/mcp_sqlpp_YYYY-MM-DD.log`
- **Log Rotation**: 100MB max file size, 10 backup files, 30-day retention
- **JSON Format**: Structured logs for better parsing and analysis
- **Dual Output**: Logs appear in both console and file
- **Compression**: Old log files are automatically compressed

**What Gets Logged:**
- MCP protocol initialization and tool calls
- Tool execution with arguments and result sizes
- sqlpp command execution and output (truncated at TRACE level)
- Error details and debugging information
- Performance metrics and timing

## MCP Tools

The server provides the following MCP tools:

**Note**: All tools internally use sqlpp's `--stdin` interface to send commands. Schema commands like `@schema-tables` and SQL statements are sent as input via stdin to the sqlpp process.

For detailed information about MCP protocol testing and tool validation, see [MCP_TESTING.md](documentation/MCP_TESTING.md).

### Schema Commands Reference

The following schema commands are supported (sent via stdin to sqlpp):
- `@drivers` - List all available database drivers
- `@schema-tables [filter]` - List database tables
- `@schema-views [filter]` - List database views
- `@schema-procedures [filter]` - List stored procedures
- `@schema-functions [filter]` - List functions
- `@schema-all [filter]` - Show all schema information

### Schema Tools

#### `list_schema_all`
Retrieve all schema information (tables, views, procedures, functions).

**Parameters:**
- `connection` (required): Database connection name
- `filter` (optional): Filter pattern for results
- `output` (optional): Output format (json, table, csv)

#### `list_schema_tables`
Retrieve table schema information.

**Parameters:** Same as `list_schema_all`

#### `list_schema_views`
Retrieve view schema information.

**Parameters:** Same as `list_schema_all`

#### `list_schema_procedures`
Retrieve stored procedure schema information.

**Parameters:** Same as `list_schema_all`

#### `list_schema_functions`
Retrieve function schema information.

**Parameters:** Same as `list_schema_all`

### Connection Management

#### `list_connections`
List all available database connections.

**Parameters:** None

### SQL Execution

#### `execute_sql_command`
Execute SQL commands against the database.

**Parameters:**
- `connection` (required): Database connection name
- `command` (required): SQL command(s) to execute
- `output` (optional): Output format

### Driver Information

#### `list_drivers`
List all available database drivers.

**Parameters:** None

## Usage Examples

### STDIO Mode (for MCP clients)

```bash
./mcp_sqlpp --transport stdio
```

### HTTP Mode (for testing and web integration)

```bash
./mcp_sqlpp --transport http --port 8080
```

### File Logging

Enable detailed file logging with automatic rolling dates:

```bash
# Enable file logging via command line
./mcp_sqlpp --file-logging --log-level trace --transport stdio

# File logging via configuration
# Set file_logging: true in config.yaml
./mcp_sqlpp --config config.yaml
```

When file logging is enabled:
- Logs are written to `logs/mcp_sqlpp_YYYY-MM-DD.log`
- Files automatically rotate (100MB max, 10 backups, 30 day retention)
- TRACE level includes truncated output previews for debugging
- All MCP calls, tool executions, and sqlpp interactions are logged

Test the health endpoint:
```bash
curl http://localhost:8080/health
```

### Using with MCP Clients

Configure your MCP client to connect to the server:

**STDIO Transport:**
```json
{
  "command": "./mcp_sqlpp",
  "args": ["--transport", "stdio"]
}
```

**HTTP Transport:**
```json
{
  "url": "http://localhost:8080/mcp"
}
```

## AWS Deployment

### Prerequisites

1. Configure AWS credentials with appropriate permissions
2. Install AWS CDK: `npm install -g aws-cdk`
3. Set up OIDC for GitHub Actions (if using CI/CD)

### Manual Deployment

1. Navigate to the CDK directory:
```bash
cd deployment/cdk
```

2. Install Python dependencies:
```bash
pip install -r requirements.txt
```

3. Bootstrap CDK (first time only):
```bash
cdk bootstrap
```

4. Deploy the stack:
```bash
ENVIRONMENT=development cdk deploy
```

### CI/CD Deployment

The project includes GitHub Actions workflows for automated deployment:

1. **Development**: Deploys on push to `develop` branch
2. **Production**: Deploys on push to `main` branch

Required GitHub secrets:
- `AWS_ROLE_ARN`: ARN of the IAM role for OIDC authentication

### Deployment Commands

From your development environment:

```bash
# Deploy to development
ENVIRONMENT=development cdk deploy --profile your-aws-profile

# Deploy to production
ENVIRONMENT=production cdk deploy --profile your-aws-profile

# Destroy resources
ENVIRONMENT=development cdk destroy --profile your-aws-profile
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed output
go test -v ./...
```

### Test Categories

**Unit Tests:**
- **Configuration**: Loading and validation of config files and environment variables
- **Tool Execution**: Validation of all MCP tools with various input scenarios
- **Connection Management**: Testing connection listing and default connection handling
- **Logging Functions**: Validation of output truncation and logging levels

**Integration Tests:**
- **Real sqlpp Integration**: Tests with actual sqlpp executable to verify end-to-end functionality
- **MCP Protocol Compliance**: Validation of proper MCP handshake and tool calling sequences
- **Error Handling**: Testing of various failure scenarios and error propagation

**Specific Test Highlights:**
- **`list_connections` Validation**: Ensures proper connection format and default connection detection
- **Empty Result Handling**: Tests graceful handling of empty or missing data
- **Output Truncation**: Validates that large outputs are properly truncated for logging
- **TRACE Level Logging**: Verification that detailed debugging logs are captured correctly

### Manual Testing

Additional manual testing scripts are available:

```bash
# Test MCP protocol compliance
bash scripts/test-mcp-protocol.sh

# Test logging functionality  
bash scripts/test-logging.sh

# Test all tool integrations
bash scripts/test-tools.sh
```

For detailed testing information, see [MCP_TESTING.md](documentation/MCP_TESTING.md).

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Project Structure

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── server/          # MCP server implementation
│   ├── sqlpp/           # sqlpp CLI wrapper
│   └── tools/           # MCP tool definitions
├── pkg/types/           # Shared types
├── deployment/
│   ├── cdk/             # AWS CDK infrastructure
│   └── docker/          # Docker configuration
├── .github/workflows/   # CI/CD pipelines
└── documentation/       # Additional documentation
```

## Monitoring and Logging

### Logging

The server provides structured logging with configurable levels:
- `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`

Logs can be output in text or JSON format for different environments.

### Health Checks

HTTP transport includes a health check endpoint at `/health` that returns:
- HTTP 200 OK when the service is healthy
- Validates sqlpp executable availability

### AWS CloudWatch

When deployed to AWS App Runner, logs are automatically sent to CloudWatch with:
- Structured JSON logging
- Request/response tracking
- Error monitoring
- Performance metrics

## Security Considerations

- **Input Validation**: All tool parameters are validated before execution
- **Process Isolation**: sqlpp runs as isolated child processes
- **Non-root Execution**: Container runs with non-root user
- **Network Security**: Configurable egress rules in AWS deployment
- **Secrets Management**: No hardcoded credentials or secrets

## Troubleshooting

For comprehensive troubleshooting and MCP protocol testing, see [MCP_TESTING.md](documentation/MCP_TESTING.md).

### Common Issues

1. **"Error: unknown flag: --command"**
   - **Cause**: This error indicates an older version of the MCP server that incorrectly tried to use a non-existent `--command` flag
   - **Solution**: Ensure you're using the latest version of mcp_sqlpp. The server should use stdin to send commands to sqlpp
   - **Verification**: Check that your server version includes the stdin-based implementation

2. **"failed to open file @schema-tables"**
   - **Cause**: Schema commands are being treated as file paths instead of special commands
   - **Solution**: Ensure schema commands are sent via stdin, not as command-line arguments
   - **Verification**: Schema commands should work when sent as input to `sqlpp --stdin`

3. **sqlpp not found**
   - Ensure sqlpp is installed and in PATH
   - Check `executable_path` configuration
   - Verify permissions

4. **Connection timeout**
   - Increase `timeout` configuration
   - Check database connectivity
   - Review sqlpp connection configuration

5. **Permission denied**
   - Check file permissions on sqlpp executable
   - Verify user permissions for database access
   - Review container user configuration

### Testing sqlpp Integration

Before using the MCP server, verify sqlpp works correctly:

```bash
# Test basic connectivity
sqlpp --list-connections

# Test SQL execution
echo "SELECT 1 as test;" | sqlpp --stdin --connection main

# Test schema commands
echo "@schema-tables" | sqlpp --stdin --connection main

# Test with output formatting
echo "SELECT 'Hello World' as message;" | sqlpp --stdin --connection main --output table
```

For detailed sqlpp integration and configuration information, see [SQLPP_INTEGRATION.md](documentation/SQLPP_INTEGRATION.md).

### Debug Mode

Enable debug logging for detailed troubleshooting:

```bash
./mcp_sqlpp --log-level debug
```

### File Logging for Debugging

For comprehensive debugging, enable file logging to capture all interactions:

```bash
# Enable file logging with TRACE level for maximum detail
./mcp_sqlpp --file-logging --log-level trace --transport stdio

# View recent logs
tail -f logs/mcp_sqlpp_$(date +%Y-%m-%d).log

# Search for specific tool or error logs
grep "list_connections" logs/mcp_sqlpp_*.log
grep "error" logs/mcp_sqlpp_*.log
```

File logging captures:
- All MCP protocol messages (initialize, tool calls, responses)
- Tool execution details with arguments and results
- sqlpp command execution with truncated output previews
- Error details and debugging information

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and linting
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and questions:
- Create an issue on GitHub
- Check the [documentation](./documentation/) directory for detailed guides:
  - [MCP_TESTING.md](documentation/MCP_TESTING.md) - MCP protocol testing and troubleshooting
  - [SQLPP_INTEGRATION.md](documentation/SQLPP_INTEGRATION.md) - sqlpp setup and integration details
- Review sqlpp documentation for database-specific issues
