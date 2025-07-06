package types

// SqlppResult represents the result of a sqlpp command execution
type SqlppResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// ToolParameter represents a parameter for MCP tools
type ToolParameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// Connection represents a database connection
type Connection struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Status string `json:"status"`
}

// Driver represents a database driver
type Driver struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version,omitempty"`
}
