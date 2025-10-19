// Package main provides output formatting utilities for the bento CLI.
//
// This file contains helpers for creating user-friendly output with:
//   - Duration formatting (ms, seconds, minutes)
//   - Success/error boxes for visual emphasis
//   - Consistent emoji usage
package main

import (
	"fmt"
	"time"
)

// formatDuration formats a duration for human readability.
//
// Examples:
//   - 45ms -> "45ms"
//   - 1.5s -> "1.5s"
//   - 3m 45s -> "3m 45s"
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		secs := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm %ds", mins, secs)
	}
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
}

// printSuccess prints a success message with bento emoji.
func printSuccess(message string) {
	fmt.Printf("\n🍱 %s\n", message)
}

// printError prints an error message with emoji.
func printError(message string) {
	fmt.Printf("\n❌ %s\n", message)
}

// printInfo prints an info message with bento emoji.
func printInfo(message string) {
	fmt.Printf("🍱 %s\n", message)
}

// printProgress prints a progress message with neta emoji.
func printProgress(message string) {
	fmt.Printf("🍙 %s\n", message)
}

// printCheck prints a check mark for completed items.
func printCheck(message string) {
	fmt.Printf("✓ %s\n", message)
}
