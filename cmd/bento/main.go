// Package main implements the bento CLI.
//
// Bento is a high-performance workflow automation CLI written in Go.
// It uses playful sushi-themed commands to make automation fun.
//
// Commands:
//   - savor: Execute a bento workflow
//   - sample: Validate a bento without executing
//   - menu: List available bentos
//   - box: Create a new bento template
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
  • savor  - Execute a bento (savor every delicious bite!)
  • sample - Validate without executing (sample to check if it tastes right)
  • menu   - List available bentos (restaurant menu)
  • box    - Create a new bento template (box up fresh ingredients)`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(savorCmd)
	rootCmd.AddCommand(sampleCmd)
	rootCmd.AddCommand(menuCmd)
	rootCmd.AddCommand(boxCmd)
	rootCmd.AddCommand(versionCmd)
}
