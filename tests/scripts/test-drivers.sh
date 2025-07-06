#!/bin/bash

# Test script to get drivers from gosqlpp MCP server
SERVER_PATH="/Users/iggybdda/Code/stainedhead/Golang/gosqlpp-mcp-server/gosqlpp-mcp-server"

echo "Getting supported database drivers..."
echo "===================================="

# Create a temporary file for the MCP session
TEMP_FILE=$(mktemp)

# Initialize the MCP session
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' > $TEMP_FILE

# Send initialized notification
echo '{"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}}' >> $TEMP_FILE

# Call the drivers tool
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "drivers", "arguments": {}}}' >> $TEMP_FILE

# Send the commands to the server
cat $TEMP_FILE | $SERVER_PATH --transport stdio 2>/dev/null | jq -r 'select(.id == 3) | .result.content[0].text'

# Clean up
rm $TEMP_FILE
