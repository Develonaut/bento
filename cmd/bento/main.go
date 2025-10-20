// Package main implements the bento CLI.
//
// Bento is a high-performance workflow automation CLI written in Go.
//
// Commands:
//   - run: Execute a bento workflow
//   - validate: Validate a bento without executing
//   - list: List available bentos
//   - new: Create a new bento template
//   - docs: View documentation
//   - secrets: Manage secrets
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
	Short: "High-performance workflow automation",
	Long: `Bento - High-performance workflow automation

Bento lets you build powerful automation workflows using composable
nodes that can be connected together.

Available Commands:
  • run      - Execute a bento workflow
  • validate - Validate a workflow without executing
  • list     - List available bento workflows
  • new      - Create a new bento workflow template
  • docs     - View documentation
  • secrets  - Manage secrets securely
  • version  - Show version information`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(docsCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(secretsCmd)
}
