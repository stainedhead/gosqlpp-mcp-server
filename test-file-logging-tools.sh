#!/bin/bash

# Test file logging with actual tool calls
echo "Testing file logging with MCP tool calls..."

# Remove any existing logs
rm -rf logs/

# Create a test MCP session with tool calls
cat << 'EOF' | ./mcp_sqlpp --file-logging --log-level trace --transport stdio
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"roots": {"listChanged": true}, "sampling": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "notifications/initialized"}
{"jsonrpc": "2.0", "id": 3, "method": "tools/list"}
{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "list_connections", "arguments": {}}}
{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "list_drivers", "arguments": {}}}
EOF

echo ""
echo "File logging test completed. Checking logs..."

# Check if log file was created
if [ -d "logs" ]; then
    echo "✓ Logs directory created"
    
    # Check if log file exists and has content
    log_file=$(ls logs/ | head -1)
    if [ -n "$log_file" ]; then
        echo "✓ Log file created: $log_file"
        echo "✓ Log file size: $(wc -c < logs/$log_file) bytes"
        echo ""
        echo "Checking for TRACE level logs..."
        
        # Check for specific trace logs
        if grep -q '"level":"trace"' logs/$log_file; then
            echo "✓ TRACE level logs found"
            echo "Sample TRACE logs:"
            grep '"level":"trace"' logs/$log_file | head -3
        else
            echo "✗ No TRACE level logs found"
        fi
        
        echo ""
        echo "Checking for tool call logs..."
        if grep -q '"tool"' logs/$log_file; then
            echo "✓ Tool call logs found"
        else
            echo "✗ No tool call logs found"
        fi
        
        echo ""
        echo "Checking for sqlpp execution logs..."
        if grep -q '"sqlpp"' logs/$log_file; then
            echo "✓ sqlpp execution logs found"
        else
            echo "✗ No sqlpp execution logs found"
        fi
        
    else
        echo "✗ No log file found"
        exit 1
    fi
else
    echo "✗ Logs directory not created"
    exit 1
fi

echo ""
echo "✓ File logging with tool calls test completed"
