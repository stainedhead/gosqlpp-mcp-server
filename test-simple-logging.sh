#!/bin/bash

# Test file logging with a simple tool call
echo "Testing file logging with MCP tool calls..."

# Remove any existing logs
rm -rf logs/

# Create a test MCP session with tool calls (corrected sequence)
cat << 'EOF' | ./mcp_sqlpp --file-logging --log-level trace --transport stdio
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"roots": {"listChanged": true}, "sampling": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}
{"jsonrpc": "2.0", "id": 3, "method": "tools/list"}
{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "list_connections", "arguments": {}}}
EOF

echo ""
echo "File logging test completed. Checking logs..."

# Check log file content
if [ -d "logs" ]; then
    log_file=$(ls logs/ | head -1)
    if [ -n "$log_file" ]; then
        echo "✓ Log file created: $log_file"
        echo "✓ Log file size: $(wc -c < logs/$log_file) bytes"
        echo ""
        echo "Full log file content:"
        cat logs/$log_file
        echo ""
        
        # Check for specific trace logs
        if grep -q '"level":"trace"' logs/$log_file; then
            echo "✓ TRACE level logs found"
        else
            echo "ℹ No TRACE level logs found (this is expected if tools weren't called)"
        fi
        
        # Check for tool related logs
        if grep -q 'tool\|list_connections' logs/$log_file; then
            echo "✓ Tool-related logs found"
        else
            echo "ℹ No tool-related logs found"
        fi
        
    else
        echo "✗ No log file found"
    fi
else
    echo "✗ Logs directory not created"
fi

echo ""
echo "✓ Test completed"
