package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration for the MCP server
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Sqlpp  SqlppConfig  `mapstructure:"sqlpp"`
	Log    LogConfig    `mapstructure:"log"`
	AWS    AWSConfig    `mapstructure:"aws"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Transport string `mapstructure:"transport"` // "stdio" or "http"
	Port      int    `mapstructure:"port"`
	Host      string `mapstructure:"host"`
}

// SqlppConfig holds sqlpp executable configuration
type SqlppConfig struct {
	ExecutablePath string `mapstructure:"executable_path"` // Directory path containing sqlpp executable (defaults to .bin)
	Timeout        int    `mapstructure:"timeout"`         // timeout in seconds
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`       // "json" or "text"
	FileLogging bool   `mapstructure:"file_logging"` // Enable file logging
}

// AWSConfig holds AWS deployment configuration
type AWSConfig struct {
	Region      string `mapstructure:"region"`
	Environment string `mapstructure:"environment"`
}

// Load loads configuration from file, environment variables, and defaults
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Look for config in current directory and home directory
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.gosqlpp-mcp-server")

		// Also check for config in executable directory
		if execPath, err := os.Executable(); err == nil {
			v.AddConfigPath(filepath.Dir(execPath))
		}
	}

	// Environment variables
	v.SetEnvPrefix("GOSQLPP_MCP")
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults and env vars
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.transport", "stdio")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "localhost")

	// Sqlpp defaults
	v.SetDefault("sqlpp.executable_path", ".bin") // Default to .bin directory
	v.SetDefault("sqlpp.timeout", 300)            // 5 minutes

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("log.file_logging", false)

	// AWS defaults
	v.SetDefault("aws.region", "us-east-1")
	v.SetDefault("aws.environment", "development")
}

// validate validates the configuration
func validate(config *Config) error {
	// Validate transport
	if config.Server.Transport != "stdio" && config.Server.Transport != "http" {
		return fmt.Errorf("invalid transport: %s (must be 'stdio' or 'http')", config.Server.Transport)
	}

	// Validate port for HTTP transport
	if config.Server.Transport == "http" {
		if config.Server.Port < 1 || config.Server.Port > 65535 {
			return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", config.Server.Port)
		}
	}

	// Validate log level
	validLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	validLevel := false
	for _, level := range validLevels {
		if config.Log.Level == level {
			validLevel = true
			break
		}
	}
	if !validLevel {
		return fmt.Errorf("invalid log level: %s", config.Log.Level)
	}

	// Validate log format
	if config.Log.Format != "json" && config.Log.Format != "text" {
		return fmt.Errorf("invalid log format: %s (must be 'json' or 'text')", config.Log.Format)
	}

	// Validate sqlpp timeout
	if config.Sqlpp.Timeout < 1 {
		return fmt.Errorf("invalid sqlpp timeout: %d (must be greater than 0)", config.Sqlpp.Timeout)
	}

	return nil
}

// GetSqlppExecutablePath returns the full path to the sqlpp executable
func (c *SqlppConfig) GetSqlppExecutablePath() string {
	if c.ExecutablePath == "" {
		return filepath.Join(".bin", "sqlpp")
	}

	// If the path already includes the executable name, return as-is
	if filepath.Base(c.ExecutablePath) == "sqlpp" {
		return c.ExecutablePath
	}

	// Otherwise, join the directory path with the executable name
	return filepath.Join(c.ExecutablePath, "sqlpp")
}
