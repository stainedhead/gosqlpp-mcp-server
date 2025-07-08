package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/stainedhead/gosqlpp-mcp-server/internal/config"
)

func main() {
	// Test the new configuration logic

	// Create a test directory with mock sqlpp
	testDir := "test-bin"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	// Create mock sqlpp executable
	sqlppPath := filepath.Join(testDir, "sqlpp")
	err = os.WriteFile(sqlppPath, []byte("#!/bin/bash\necho 'mock sqlpp'\n"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Test configuration
	v := viper.New()
	v.Set("sqlpp.executable_path", testDir)

	cfg := &config.Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Test the new method
	execPath := cfg.Sqlpp.GetSqlppExecutablePath()
	fmt.Printf("Configuration executable_path: %s\n", cfg.Sqlpp.ExecutablePath)
	fmt.Printf("Resolved executable path: %s\n", execPath)

	// Verify the file exists
	if _, err := os.Stat(execPath); err != nil {
		log.Fatalf("Sqlpp executable not found at %s: %v", execPath, err)
	}

	fmt.Println("✓ New configuration logic works correctly!")
}
