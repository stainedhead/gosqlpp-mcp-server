package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/sqlpp"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/tools"
)

func main() {
	logger := logrus.New()
	executor := sqlpp.NewExecutor("../gosqlpp/sqlpp", 30, logger)
	handler := tools.NewToolHandler(executor, logger)

	toolList := handler.GetTools()
	fmt.Printf("Available MCP tools (%d total):\n", len(toolList))
	for _, tool := range toolList {
		fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
	}
}
