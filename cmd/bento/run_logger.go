// Package main implements logger creation for the run command.
//
// This file contains functions for creating file loggers and dual loggers
// (file + stdout) for bento execution.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Develonaut/bento/pkg/logs"
	"github.com/Develonaut/bento/pkg/shoyu"
)

// createFileLogger creates a logger that writes to ~/.bento/logs/
// Returns the logger, the log file (for cleanup), and any error.
func createFileLogger() (*shoyu.Logger, *os.File, error) {
	// Ensure logs directory exists
	if err := logs.EnsureLogsDirectory(); err != nil {
		return nil, nil, err
	}

	// Get logs directory path
	logsDir, err := logs.GetLogsDirectory()
	if err != nil {
		return nil, nil, err
	}

	// Generate log file name
	logFileName := logs.GenerateLogFileName()
	logPath := filepath.Join(logsDir, logFileName)

	// Trim log file if it exceeds threshold (10,000 lines → 5,000 lines)
	// This keeps the log file bounded and prevents unlimited growth
	if err := logs.TrimLogFile(logPath, 10000, 5000); err != nil { //nolint:staticcheck // Intentionally ignoring trim errors - file might not exist yet
		// Don't fail on trim errors - just continue and try to log
		// (file might not exist yet, which is fine)
	}

	// Open in append mode (create if doesn't exist)
	// This allows multiple bento runs to append to the same log
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}

	level := shoyu.LevelInfo
	if verboseFlag {
		level = shoyu.LevelDebug
	}

	logger := shoyu.New(shoyu.Config{
		Level:  level,
		Output: logFile,
	})

	// Add blank line before execution for separation
	logger.Info("")

	return logger, logFile, nil
}

// createDualLogger creates a logger that writes to both file and stdout.
func createDualLogger(fileLogger *shoyu.Logger) *shoyu.Logger {
	level := shoyu.LevelInfo
	if verboseFlag {
		level = shoyu.LevelDebug
	}

	return shoyu.New(shoyu.Config{
		Level: level,
		// Enable streaming output for long-running processes
		// This outputs lines from shell-command neta in real-time
		OnStream: func(line string) {
			fmt.Println(line)
		},
	})
}
