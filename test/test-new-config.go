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
	// Test the new configuration logic with binary-relative path resolution

	fmt.Println("Testing new configuration logic with binary-relative paths...")
	fmt.Println()

	// Get current binary directory for reference
	binaryPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	binaryDir := filepath.Dir(binaryPath)
	fmt.Printf("Current binary directory: %s\n", binaryDir)
	fmt.Println()

	// Test 1: Relative path resolution
	fmt.Println("Test 1: Relative path resolution")
	testDir := "test-bin"
	fullTestDir := filepath.Join(binaryDir, testDir)

	err = os.MkdirAll(fullTestDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(fullTestDir)

	// Create mock sqlpp executable
	sqlppPath := filepath.Join(fullTestDir, "sqlpp")
	err = os.WriteFile(sqlppPath, []byte("#!/bin/bash\necho 'mock sqlpp'\n"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Test configuration with relative path
	v := viper.New()
	v.Set("sqlpp.executable_path", testDir) // Relative path

	cfg := &config.Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	execPath := cfg.Sqlpp.GetSqlppExecutablePath()
	fmt.Printf("  Configuration executable_path: %s\n", cfg.Sqlpp.ExecutablePath)
	fmt.Printf("  Resolved executable path: %s\n", execPath)
	fmt.Printf("  Expected path: %s\n", sqlppPath)

	// Verify the file exists
	if _, err := os.Stat(execPath); err != nil {
		log.Fatalf("  ❌ Sqlpp executable not found at %s: %v", execPath, err)
	}

	if execPath == sqlppPath {
		fmt.Println("  ✓ Relative path resolved correctly relative to binary directory!")
	} else {
		log.Fatalf("  ❌ Path mismatch: expected %s, got %s", sqlppPath, execPath)
	}
	fmt.Println()

	// Test 2: Absolute path (should be unchanged)
	fmt.Println("Test 2: Absolute path preservation")
	absolutePath := "/usr/local/bin"
	v2 := viper.New()
	v2.Set("sqlpp.executable_path", absolutePath)

	cfg2 := &config.Config{}
	err = v2.Unmarshal(cfg2)
	if err != nil {
		log.Fatal(err)
	}

	execPath2 := cfg2.Sqlpp.GetSqlppExecutablePath()
	expectedAbsolute := filepath.Join(absolutePath, "sqlpp")
	fmt.Printf("  Configuration executable_path: %s\n", cfg2.Sqlpp.ExecutablePath)
	fmt.Printf("  Resolved executable path: %s\n", execPath2)
	fmt.Printf("  Expected path: %s\n", expectedAbsolute)

	if execPath2 == expectedAbsolute {
		fmt.Println("  ✓ Absolute path preserved correctly!")
	} else {
		log.Fatalf("  ❌ Path mismatch: expected %s, got %s", expectedAbsolute, execPath2)
	}
	fmt.Println()

	// Test 3: Default path (.bin)
	fmt.Println("Test 3: Default path resolution")
	v3 := viper.New() // No executable_path set, should use default

	cfg3 := &config.Config{}
	err = v3.Unmarshal(cfg3)
	if err != nil {
		log.Fatal(err)
	}

	execPath3 := cfg3.Sqlpp.GetSqlppExecutablePath()
	expectedDefault := filepath.Join(binaryDir, ".bin", "sqlpp")
	fmt.Printf("  Configuration executable_path: '%s' (empty, using default)\n", cfg3.Sqlpp.ExecutablePath)
	fmt.Printf("  Resolved executable path: %s\n", execPath3)
	fmt.Printf("  Expected path: %s\n", expectedDefault)

	if execPath3 == expectedDefault {
		fmt.Println("  ✓ Default path resolved correctly relative to binary directory!")
	} else {
		log.Fatalf("  ❌ Path mismatch: expected %s, got %s", expectedDefault, execPath3)
	}

	fmt.Println()
	fmt.Println("✅ All tests passed! New configuration logic works correctly!")
	fmt.Println("   - Relative paths are resolved relative to the MCP server binary directory")
	fmt.Println("   - Absolute paths are preserved as-is")
	fmt.Println("   - Default .bin directory is resolved relative to binary directory")
}
