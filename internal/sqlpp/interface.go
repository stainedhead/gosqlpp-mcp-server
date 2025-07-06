package sqlpp

import "github.com/stainedhead/gosqlpp-mcp-server/pkg/types"

// ExecutorInterface defines the interface for sqlpp command execution
type ExecutorInterface interface {
	ExecuteSchemaCommand(schemaType, connection, filter, output string) (*types.SqlppResult, error)
	ExecuteSQLCommand(connection, command, output string) (*types.SqlppResult, error)
	ListConnections() (*types.SqlppResult, error)
	ListDrivers() (*types.SqlppResult, error)
	ValidateExecutable() error
}
