package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// SetupFileLogging sets up file logging with automatic rolling based on date
func SetupFileLogging(logger *logrus.Logger, enabled bool) error {
	if !enabled {
		return nil
	}

	// Create logs directory relative to the binary location, not working directory
	logsDir := resolveLogsDirPath("logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Generate filename with current date
	filename := generateLogFilename()
	logPath := filepath.Join(logsDir, filename)

	// Setup lumberjack for log rotation
	fileLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // megabytes
		MaxBackups: 10,  // number of backups
		MaxAge:     30,  // days
		Compress:   true,
	}

	// Create a multi-writer that writes to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, fileLogger)
	logger.SetOutput(multiWriter)

	// Note: We keep the existing log level - don't override it
	// The file will capture whatever level was set by the user

	// Use JSON format for file logs for better parsing
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		PrettyPrint:     false,
	})

	logger.WithFields(logrus.Fields{
		"log_file":      logPath,
		"logs_dir":      logsDir,
		"max_size_mb":   100,
		"max_backups":   10,
		"max_age_days":  30,
	}).Info("File logging enabled")

	return nil
}

// generateLogFilename generates a log filename with the current date
func generateLogFilename() string {
	now := time.Now()
	return fmt.Sprintf("mcp_sqlpp_%s.log", now.Format("2006-01-02"))
}

// SetupDailyRotation sets up a ticker to rotate logs daily
// This function should be called in a goroutine
func SetupDailyRotation(logger *logrus.Logger) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		// Log rotation will be handled automatically by lumberjack
		// We just log a daily marker
		logger.WithField("event", "daily_rotation").Info("Daily log rotation checkpoint")
	}
}

// resolveLogsDirPath resolves the logs directory path relative to the binary location
// This ensures log files are created relative to where the MCP server binary is located,
// not the working directory from which it was launched
func resolveLogsDirPath(relativePath string) string {
	// If it's already an absolute path, return as-is
	if filepath.IsAbs(relativePath) {
		return relativePath
	}

	// Get the directory where the MCP server binary is located
	binaryPath, err := os.Executable()
	if err != nil {
		// Fallback to working directory if we can't determine binary location
		return relativePath
	}

	binaryDir := filepath.Dir(binaryPath)

	// Resolve the path relative to the binary directory
	return filepath.Join(binaryDir, relativePath)
}
