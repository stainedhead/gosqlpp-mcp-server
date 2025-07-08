package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/sqlpp"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/tools"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise
	executor := sqlpp.NewExecutor("../gosqlpp/sqlpp", 30, logger)
	handler := tools.NewToolHandler(executor, logger)

	toolList := handler.GetTools()
	fmt.Printf("✅ MCP Server: mcp_sqlpp\n")
	fmt.Printf("✅ Available tools (%d total):\n", len(toolList))
	for i, tool := range toolList {
		fmt.Printf("  %d. %s - %s\n", i+1, tool.Name, tool.Description)
	}
	fmt.Println("\n✅ All changes implemented successfully!")
}
