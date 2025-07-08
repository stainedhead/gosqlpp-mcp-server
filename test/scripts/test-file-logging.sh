#!/bin/bash

# Test file logging functionality
# This script should be run from the project root directory

# Change to project root if not already there
if [ ! -f "mcp_sqlpp" ] && [ ! -f "go.mod" ]; then
    echo "Error: This script must be run from the project root directory"
    echo "Usage: ./test/scripts/test-file-logging.sh"
    exit 1
fi

echo "Testing file logging functionality..."

# Remove any existing logs
rm -rf logs/

# Start the server with file logging enabled in the background
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"roots": {"listChanged": true}, "sampling": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | ./mcp_sqlpp --file-logging --log-level trace &
SERVER_PID=$!

# Give the server time to start and process the message
sleep 3

# Kill the background process
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

# Check if log file was created
if [ -d "logs" ]; then
    echo "✓ Logs directory created"
    ls -la logs/
    
    # Check if log file exists and has content
    log_file=$(ls logs/ | head -1)
    if [ -n "$log_file" ]; then
        echo "✓ Log file created: $log_file"
        echo "✓ Log file size: $(wc -c < logs/$log_file) bytes"
        echo "First few lines of log file:"
        head -5 logs/$log_file
        echo "..."
        echo "Last few lines of log file:"
        tail -5 logs/$log_file
    else
        echo "✗ No log file found"
        exit 1
    fi
else
    echo "✗ Logs directory not created"
    exit 1
fi

echo "✓ File logging test completed successfully"
