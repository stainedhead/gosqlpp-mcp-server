package sqlpp

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stainedhead/gosqlpp-mcp-server/pkg/types"
)

// Executor handles execution of sqlpp commands
type Executor struct {
	executablePath string
	timeout        time.Duration
	logger         *logrus.Logger
}

// NewExecutor creates a new sqlpp executor
func NewExecutor(executablePath string, timeoutSeconds int, logger *logrus.Logger) *Executor {
	return &Executor{
		executablePath: executablePath,
		timeout:        time.Duration(timeoutSeconds) * time.Second,
		logger:         logger,
	}
}

// ExecuteSchemaCommand executes a schema-related command (@schema-*)
func (e *Executor) ExecuteSchemaCommand(schemaType, connection, filter, output string) (*types.SqlppResult, error) {
	// Build the schema command
	schemaCommand := fmt.Sprintf("@schema-%s", schemaType)
	if filter != "" {
		schemaCommand = fmt.Sprintf("%s %s", schemaCommand, filter)
	}
	
	return e.executeStdinCommandWithOptions(schemaCommand, connection, output)
}

// ExecuteSQLCommand executes a SQL command
func (e *Executor) ExecuteSQLCommand(connection, command, output string) (*types.SqlppResult, error) {
	return e.executeStdinCommandWithOptions(command, connection, output)
}

// ListConnections lists available database connections
func (e *Executor) ListConnections() (*types.SqlppResult, error) {
	args := []string{"--list-connections"}
	return e.executeCommand(args)
}

// ListDrivers lists available database drivers
func (e *Executor) ListDrivers() (*types.SqlppResult, error) {
	return e.executeStdinCommand("@drivers")
}

// executeCommand executes a sqlpp command with the given arguments
func (e *Executor) executeCommand(args []string) (*types.SqlppResult, error) {
	e.logger.WithFields(logrus.Fields{
		"executable": e.executablePath,
		"args":       args,
		"timeout":    e.timeout,
	}).Debug("Executing sqlpp command")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, e.executablePath, args...)
	
	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	err := cmd.Run()
	
	result := &types.SqlppResult{
		Success: err == nil,
		Output:  strings.TrimSpace(stdout.String()),
	}

	if err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		if stderrStr != "" {
			result.Error = stderrStr
		} else {
			result.Error = err.Error()
		}
		
		e.logger.WithFields(logrus.Fields{
			"error":  err,
			"stderr": stderrStr,
			"args":   args,
		}).Error("sqlpp command failed")
	} else {
		e.logger.WithFields(logrus.Fields{
			"args":        args,
			"output_size": len(result.Output),
		}).Debug("sqlpp command succeeded")
	}

	return result, nil
}

// executeStdinCommandWithOptions executes a sqlpp command by sending input via stdin with connection and output options
func (e *Executor) executeStdinCommandWithOptions(input, connection, output string) (*types.SqlppResult, error) {
	args := []string{"--stdin"}
	
	if connection != "" {
		args = append(args, "--connection", connection)
	}
	
	if output != "" {
		args = append(args, "--output", output)
	}

	e.logger.WithFields(logrus.Fields{
		"executable": e.executablePath,
		"args":       args,
		"input":      input,
		"timeout":    e.timeout,
	}).Debug("Executing sqlpp command with stdin and options")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// Create command with args
	cmd := exec.CommandContext(ctx, e.executablePath, args...)
	
	// Set up stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Write input to stdin and close it
	if _, err := stdin.Write([]byte(input)); err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to write to stdin: %w", err)
	}
	stdin.Close()

	// Wait for command to complete
	err = cmd.Wait()
	
	result := &types.SqlppResult{
		Success: err == nil,
		Output:  strings.TrimSpace(stdout.String()),
	}

	if err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		if stderrStr != "" {
			result.Error = stderrStr
		} else {
			result.Error = err.Error()
		}
		
		e.logger.WithFields(logrus.Fields{
			"error":  err,
			"stderr": stderrStr,
			"args":   args,
			"input":  input,
		}).Error("sqlpp stdin command with options failed")
	} else {
		e.logger.WithFields(logrus.Fields{
			"args":        args,
			"input":       input,
			"output_size": len(result.Output),
		}).Debug("sqlpp stdin command with options succeeded")
	}

	return result, nil
}

// executeStdinCommand executes a sqlpp command by sending input via stdin
func (e *Executor) executeStdinCommand(input string) (*types.SqlppResult, error) {
	e.logger.WithFields(logrus.Fields{
		"executable": e.executablePath,
		"input":      input,
		"timeout":    e.timeout,
	}).Debug("Executing sqlpp command with stdin")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// Create command with --stdin flag
	cmd := exec.CommandContext(ctx, e.executablePath, "--stdin")
	
	// Set up stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Write input to stdin and close it
	if _, err := stdin.Write([]byte(input)); err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to write to stdin: %w", err)
	}
	stdin.Close()

	// Wait for command to complete
	err = cmd.Wait()
	
	result := &types.SqlppResult{
		Success: err == nil,
		Output:  strings.TrimSpace(stdout.String()),
	}

	if err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		if stderrStr != "" {
			result.Error = stderrStr
		} else {
			result.Error = err.Error()
		}
		
		e.logger.WithFields(logrus.Fields{
			"error":  err,
			"stderr": stderrStr,
			"input":  input,
		}).Error("sqlpp stdin command failed")
	} else {
		e.logger.WithFields(logrus.Fields{
			"input":       input,
			"output_size": len(result.Output),
		}).Debug("sqlpp stdin command succeeded")
	}

	return result, nil
}

// ValidateExecutable checks if the sqlpp executable is available and working
func (e *Executor) ValidateExecutable() error {
	e.logger.WithField("executable", e.executablePath).Debug("Validating sqlpp executable")
	
	// Try to run sqlpp with --version or --help to check if it's available
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.executablePath, "--help")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		if stderrStr != "" {
			return fmt.Errorf("sqlpp executable validation failed: %s", stderrStr)
		}
		return fmt.Errorf("sqlpp executable not found or not working: %w", err)
	}

	e.logger.Info("sqlpp executable validated successfully")
	return nil
}
