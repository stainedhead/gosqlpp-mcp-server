package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()
	
	// Change to temp directory to avoid loading any existing config
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	config, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Check defaults
	assert.Equal(t, "stdio", config.Server.Transport)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, "sqlpp", config.Sqlpp.ExecutablePath)
	assert.Equal(t, 300, config.Sqlpp.Timeout)
	assert.Equal(t, "info", config.Log.Level)
	assert.Equal(t, "text", config.Log.Format)
	assert.Equal(t, "us-east-1", config.AWS.Region)
	assert.Equal(t, "development", config.AWS.Environment)
}

func TestLoad_FromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test-config.yaml")
	
	configContent := `
server:
  transport: "http"
  port: 9090
  host: "0.0.0.0"
sqlpp:
  executable_path: "/usr/local/bin/sqlpp"
  timeout: 600
log:
  level: "debug"
  format: "json"
aws:
  region: "us-west-2"
  environment: "test"
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	config, err := Load(configFile)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Check loaded values
	assert.Equal(t, "http", config.Server.Transport)
	assert.Equal(t, 9090, config.Server.Port)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, "/usr/local/bin/sqlpp", config.Sqlpp.ExecutablePath)
	assert.Equal(t, 600, config.Sqlpp.Timeout)
	assert.Equal(t, "debug", config.Log.Level)
	assert.Equal(t, "json", config.Log.Format)
	assert.Equal(t, "us-west-2", config.AWS.Region)
	assert.Equal(t, "test", config.AWS.Environment)
}

func TestLoad_FromEnvironment(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Clear any existing environment variables first
	envVars := []string{
		"GOSQLPP_MCP_SERVER_TRANSPORT",
		"GOSQLPP_MCP_SERVER_PORT", 
		"GOSQLPP_MCP_SQLPP_EXECUTABLE_PATH",
		"GOSQLPP_MCP_LOG_LEVEL",
	}
	
	// Save original values
	originalVals := make(map[string]string)
	for _, envVar := range envVars {
		originalVals[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}
	
	// Set test environment variables
	os.Setenv("GOSQLPP_MCP_SERVER_TRANSPORT", "http")
	os.Setenv("GOSQLPP_MCP_SERVER_PORT", "3000")
	os.Setenv("GOSQLPP_MCP_SQLPP_EXECUTABLE_PATH", "/custom/sqlpp")
	os.Setenv("GOSQLPP_MCP_LOG_LEVEL", "error")
	
	defer func() {
		// Restore original environment variables
		for _, envVar := range envVars {
			if originalVal, exists := originalVals[envVar]; exists && originalVal != "" {
				os.Setenv(envVar, originalVal)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()

	config, err := Load("")
	require.NoError(t, err)
	require.NotNil(t, config)

	// Check environment variable values
	assert.Equal(t, "http", config.Server.Transport)
	assert.Equal(t, 3000, config.Server.Port)
	assert.Equal(t, "/custom/sqlpp", config.Sqlpp.ExecutablePath)
	assert.Equal(t, "error", config.Log.Level)
}

func TestValidate_InvalidTransport(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Transport: "invalid",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
		Sqlpp: SqlppConfig{
			Timeout: 300,
		},
	}

	err := validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transport")
}

func TestValidate_InvalidPort(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Transport: "http",
			Port:      -1,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
		Sqlpp: SqlppConfig{
			Timeout: 300,
		},
	}

	err := validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid port")
}

func TestValidate_InvalidLogLevel(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Transport: "stdio",
		},
		Log: LogConfig{
			Level:  "invalid",
			Format: "text",
		},
		Sqlpp: SqlppConfig{
			Timeout: 300,
		},
	}

	err := validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid log level")
}

func TestValidate_InvalidLogFormat(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Transport: "stdio",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "invalid",
		},
		Sqlpp: SqlppConfig{
			Timeout: 300,
		},
	}

	err := validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid log format")
}

func TestValidate_InvalidTimeout(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Transport: "stdio",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
		Sqlpp: SqlppConfig{
			Timeout: 0,
		},
	}

	err := validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sqlpp timeout")
}

func TestValidate_Valid(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Transport: "stdio",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
		Sqlpp: SqlppConfig{
			Timeout: 300,
		},
	}

	err := validate(config)
	assert.NoError(t, err)
}
