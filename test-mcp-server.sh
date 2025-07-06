#!/bin/bash

echo "Testing gosqlpp MCP Server..."
echo "=============================="

# Test if the server starts and responds to basic MCP protocol
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | /Users/iggybdda/Code/stainedhead/Golang/gosqlpp-mcp-server/gosqlpp-mcp-server --transport stdio

echo ""
echo "If you see an 'initialized' response above, the MCP server is working correctly!"
echo "You can now use it in Q CLI sessions with tools like:"
echo "- mcp_docker___list_connections"
echo "- mcp_docker___schema_tables" 
echo "- mcp_docker___execute_sql_command"
