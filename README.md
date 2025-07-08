# mcp_sqlpp

A Model Context Protocol (MCP) server that provides access to the [sqlpp](https://github.com/stainedhead/gosqlpp) database CLI tool through standardized tool interfaces. This enables AI development tools and agents to interact with databases through a secure, controlled interface.

## Features

- **MCP Protocol Compliance**: Full support for the Model Context Protocol specification
- **Dual Transport Support**: Both STDIO and HTTP+SSE transports for flexible integration
- **Database Schema Tools**: Access table, view, procedure, and function schemas
- **SQL Execution**: Execute SQL commands with proper output formatting
- **Connection Management**: List and manage database connections
- **Driver Information**: Query available database drivers
- **Containerized Deployment**: Docker support with AWS App Runner deployment
- **Production Ready**: Comprehensive logging, health checks, and monitoring

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

## Installation

### Prerequisites

- Go 1.21 or later
- [sqlpp](https://github.com/stainedhead/gosqlpp) CLI tool installed and configured
  - **Required**: sqlpp version that supports `--stdin` flag and schema commands (`@schema-*`)
  - **Verify**: Run `sqlpp --help` to ensure `--stdin` flag is available
  - **Test**: Verify schema commands work: `echo "@drivers" | sqlpp --stdin`
- Docker (for containerized deployment)
- AWS CLI and CDK (for AWS deployment)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/stainedhead/mcp_sqlpp.git
cd mcp_sqlpp
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o gosqlpp-mcp-server ./cmd/server
```

4. Run with default configuration:
```bash
./gosqlpp-mcp-server
```

### Docker Deployment

1. Build the Docker image:
```bash
docker build -f deployment/docker/Dockerfile -t gosqlpp-mcp-server .
```

2. Run the container:
```bash
docker run -p 8080:8080 gosqlpp-mcp-server --transport http
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
```

### Command Line Flags

```bash
./gosqlpp-mcp-server --help
```

Common options:
- `--transport, -t`: Transport mode (stdio, http)
- `--port, -p`: HTTP server port
- `--config, -c`: Configuration file path
- `--log-level, -l`: Log level

## MCP Tools

The server provides the following MCP tools:

**Note**: All tools internally use sqlpp's `--stdin` interface to send commands. Schema commands like `@schema-tables` and SQL statements are sent as input via stdin to the sqlpp process.

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
./gosqlpp-mcp-server --transport stdio
```

### HTTP Mode (for testing and web integration)

```bash
./gosqlpp-mcp-server --transport http --port 8080
```

Test the health endpoint:
```bash
curl http://localhost:8080/health
```

### Using with MCP Clients

Configure your MCP client to connect to the server:

**STDIO Transport:**
```json
{
  "command": "./gosqlpp-mcp-server",
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
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test package
go test ./internal/config
```

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

### Common Issues

1. **"Error: unknown flag: --command"**
   - **Cause**: This error indicates an older version of the MCP server that incorrectly tried to use a non-existent `--command` flag
   - **Solution**: Ensure you're using the latest version of gosqlpp-mcp-server. The server should use stdin to send commands to sqlpp
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

### Debug Mode

Enable debug logging for detailed troubleshooting:

```bash
./gosqlpp-mcp-server --log-level debug
```

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
- Check the [documentation](./documentation/) directory
- Review sqlpp documentation for database-specific issues
