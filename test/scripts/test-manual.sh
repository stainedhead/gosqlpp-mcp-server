#!/bin/bash

# Manual test script for gosqlpp MCP server
# This script should be run from the project root directory

set -e

# Change to project root if not already there
if [ ! -f "go.mod" ]; then
    echo "Error: This script must be run from the project root directory"
    echo "Usage: ./test/scripts/test-manual.sh"
    exit 1
fi

echo "=== Manual Test Script for gosqlpp MCP Server ==="

# Build the application
echo "Building application..."
go build -o gosqlpp-mcp-server ./cmd/server

# Test help command
echo -e "\n1. Testing help command:"
./gosqlpp-mcp-server --help

# Test with invalid sqlpp path (should fail gracefully)
echo -e "\n2. Testing with invalid sqlpp path:"
if ./gosqlpp-mcp-server --config /dev/null 2>&1 | grep -q "sqlpp validation failed"; then
    echo "✓ Correctly failed with invalid sqlpp path"
else
    echo "✗ Did not fail as expected with invalid sqlpp path"
fi

# Create a mock sqlpp for testing
echo -e "\n3. Creating mock sqlpp executable:"
mkdir -p mock-bin
cat > mock-bin/sqlpp << 'EOF'
#!/bin/bash
case "$1" in
    "--help")
        echo "Mock sqlpp v1.0.0"
        echo "Usage: sqlpp [options] [command]"
        ;;
    "--list-connections")
        echo '["test-connection", "prod-connection"]'
        ;;
    "@drivers")
        echo '["mysql", "postgresql", "sqlite"]'
        ;;
    "@schema-tables")
        echo '{"tables": [{"name": "users", "columns": ["id", "name", "email"]}, {"name": "orders", "columns": ["id", "user_id", "total"]}]}'
        ;;
    *)
        echo "Mock sqlpp called with: $@"
        echo '{"result": "success", "message": "Mock response"}'
        ;;
esac
EOF

chmod +x mock-bin/sqlpp

# Create test config
echo -e "\n4. Creating test configuration:"
# Get current directory to create absolute path for sqlpp executable
CURRENT_DIR=$(pwd)
cat > test-config.yaml << EOF
server:
  transport: "http"
  host: "localhost"
  port: 8082

sqlpp:
  executable_path: "$CURRENT_DIR/mock-bin"  # Absolute path to directory containing mock sqlpp executable
  timeout: 30

log:
  level: "debug"
  format: "text"

aws:
  region: "us-east-1"
  environment: "test"
EOF

# Test server startup (HTTP mode)
echo -e "\n5. Testing server startup in HTTP mode:"
echo "Starting server in background..."
./gosqlpp-mcp-server --config test-config.yaml &
SERVER_PID=$!

# Give server time to start
sleep 2

# Test health endpoint
echo "Testing health endpoint..."
if curl -s http://localhost:8082/health | grep -q "OK"; then
    echo "✓ Health endpoint working"
else
    echo "✗ Health endpoint not working"
fi

# Stop the server
echo "Stopping server..."
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

# Test STDIO mode (just validate it starts)
echo -e "\n6. Testing STDIO mode startup:"
echo "quit" | timeout 5s ./gosqlpp-mcp-server --config test-config.yaml --transport stdio || echo "✓ STDIO mode started and stopped"

# Cleanup
echo -e "\n7. Cleaning up:"
rm -f gosqlpp-mcp-server mock-sqlpp test-config.yaml

echo -e "\n=== Manual tests completed ==="
