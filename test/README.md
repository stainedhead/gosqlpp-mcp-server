# Test Directory Structure

This directory contains all test-related files, scripts, and data for the gosqlpp MCP server project.

## Directory Structure

```
test/
├── README.md                          # This file
├── config/                            # Test configuration files
│   └── test-config-file-logging.yaml  # Configuration for testing file logging
├── data/                              # Test data files
│   ├── final-test.json                # Final test MCP commands
│   ├── test-mcp-commands.json         # Standard MCP test commands
│   └── test.db                        # Test SQLite database (if present)
├── scripts/                           # Test scripts
│   ├── test-file-logging.sh           # Test file logging functionality
│   ├── test-file-logging-tools.sh     # Test file logging with tool calls
│   ├── test-manual.sh                 # Manual testing script for the server
│   └── test-simple-logging.sh         # Simple logging test
└── test-new-config.go                 # Go program to test new configuration logic
```

## Test Scripts

### File Logging Tests
- **test-file-logging.sh**: Basic file logging functionality test
- **test-file-logging-tools.sh**: File logging with actual MCP tool calls
- **test-simple-logging.sh**: Simple logging test with corrected MCP sequence

### Manual Testing
- **test-manual.sh**: Comprehensive manual testing script that:
  - Builds the application
  - Tests help command
  - Tests invalid sqlpp path handling
  - Creates mock sqlpp executable
  - Tests server startup and configuration

### Configuration Testing
- **test-new-config.go**: Go program to verify the new configuration logic works correctly

## Test Data

### MCP Command Files
- **test-mcp-commands.json**: Standard MCP protocol commands for testing
- **final-test.json**: Final test commands for driver listing

### Configuration Files
- **test-config-file-logging.yaml**: Configuration specifically for testing file logging with trace level

### Database Files
- **test.db**: SQLite test database (created as needed by tests)

## Usage

### Running Test Scripts
All test scripts should be run from the project root directory:

```bash
# Manual testing
./test/scripts/test-manual.sh

# File logging tests
./test/scripts/test-file-logging.sh
./test/scripts/test-file-logging-tools.sh
./test/scripts/test-simple-logging.sh

# Configuration testing
cd test && go run test-new-config.go
```

### Using Test Configuration
```bash
# Run server with test file logging configuration
./mcp_sqlpp --config test/config/test-config-file-logging.yaml

# Test with MCP commands
cat test/data/test-mcp-commands.json | ./mcp_sqlpp --transport stdio
```

## Important Notes

### Path Resolution Behavior
The MCP server resolves relative paths in `executable_path` relative to the binary's location, not the working directory. This ensures the server can find sqlpp regardless of where it's launched from.

- **test-new-config.go** demonstrates this behavior with comprehensive tests
- **Test configurations** use absolute paths to avoid confusion during testing
- **Test scripts** create mock executables in predictable locations relative to the project root

## Notes

- All test scripts expect to be run from the project root directory
- Test scripts may create temporary files in the project root (mock executables, log files, etc.)
- The test database (test.db) may be created by various tests and scripts
- Test configurations may reference files outside the test directory (e.g., ../gosqlpp for the sqlpp executable)
