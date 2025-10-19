// Package main implements the bento CLI.
//
// Bento is a high-performance workflow automation CLI written in Go.
// It uses playful sushi-themed commands to make automation fun.
//
// Commands:
//   - taste: Execute a bento workflow
//   - sniff: Validate a bento without executing
//   - menu: List available bentos
//   - pack: Create a new bento template
//
// Learn more: https://github.com/Develonaut/bento
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev" // Set by build process

var rootCmd = &cobra.Command{
	Use:   "bento",
	Short: "🍱 High-performance workflow automation",
	Long: `Bento - Workflow automation with a taste of sushi 🍱

Bento lets you build powerful automation workflows using composable
"neta" (ingredients) that can be connected together like a carefully
crafted bento box.

Commands are playfully themed:
  • taste - Execute a bento (taste it to see if it works!)
  • sniff - Validate without executing (sniff to check if it's fresh)
  • menu  - List available bentos (restaurant menu)
  • pack  - Create a new bento template (pack ingredients into a box)`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(tasteCmd)
	rootCmd.AddCommand(sniffCmd)
	rootCmd.AddCommand(menuCmd)
	rootCmd.AddCommand(packCmd)
	rootCmd.AddCommand(versionCmd)
}
