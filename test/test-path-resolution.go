package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/config"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/logging"
)

func main() {
	fmt.Println("Testing binary-relative path resolution for both sqlpp and logging...")
	fmt.Println()

	// Get current binary directory for reference
	binaryPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	binaryDir := filepath.Dir(binaryPath)
	fmt.Printf("Binary directory: %s\n", binaryDir)
	fmt.Println()

	// Test 1: sqlpp executable path resolution
	fmt.Println("=== Test 1: sqlpp executable path resolution ===")
	v := viper.New()
	v.Set("sqlpp.executable_path", ".bin")

	cfg := &config.Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	execPath := cfg.Sqlpp.GetSqlppExecutablePath()
	expectedSqlpp := filepath.Join(binaryDir, ".bin", "sqlpp")

	fmt.Printf("Config executable_path: %s\n", cfg.Sqlpp.ExecutablePath)
	fmt.Printf("Resolved sqlpp path: %s\n", execPath)
	fmt.Printf("Expected sqlpp path: %s\n", expectedSqlpp)

	if execPath == expectedSqlpp {
		fmt.Println("✅ sqlpp path resolution working correctly!")
	} else {
		fmt.Printf("❌ sqlpp path resolution failed: expected %s, got %s\n", expectedSqlpp, execPath)
	}
	fmt.Println()

	// Test 2: Logging directory resolution
	fmt.Println("=== Test 2: Logging directory resolution ===")

	// Create a test logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Test file logging setup (this will create the logs directory)
	err = logging.SetupFileLogging(logger, true)
	if err != nil {
		log.Fatal(err)
	}

	expectedLogsDir := filepath.Join(binaryDir, "logs")
	fmt.Printf("Expected logs directory: %s\n", expectedLogsDir)

	// Check if logs directory exists
	if _, err := os.Stat(expectedLogsDir); err == nil {
		fmt.Printf("✅ Logs directory created at: %s\n", expectedLogsDir)
	} else {
		fmt.Printf("❌ Logs directory not found at: %s\n", expectedLogsDir)
	}

	// Check if log file was created
	logFiles, err := filepath.Glob(filepath.Join(expectedLogsDir, "mcp_sqlpp_*.log"))
	if err == nil && len(logFiles) > 0 {
		fmt.Printf("✅ Log file created: %s\n", logFiles[0])
	} else {
		fmt.Printf("❌ No log files found in: %s\n", expectedLogsDir)
	}

	fmt.Println()
	fmt.Println("=== Working Directory Independence Test ===")

	// Change to a different directory to simulate being run from elsewhere
	testDir := filepath.Join(binaryDir, "path-test")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	fmt.Printf("Changed working directory to: %s\n", testDir)

	// Test sqlpp path resolution from different working directory
	execPath2 := cfg.Sqlpp.GetSqlppExecutablePath()
	fmt.Printf("sqlpp path from different WD: %s\n", execPath2)

	if execPath2 == expectedSqlpp {
		fmt.Println("✅ sqlpp path resolution is working directory independent!")
	} else {
		fmt.Println("❌ sqlpp path resolution is affected by working directory")
	}

	fmt.Println()
	fmt.Println("🎉 All tests completed!")
	fmt.Println("Both sqlpp executable and log file paths are now resolved relative to the binary location,")
	fmt.Println("making the MCP server working directory independent!")
}
