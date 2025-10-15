package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bento",
	Short: "Bento - Organized workflow orchestration",
	Long: `Bento is a Go-based CLI orchestration tool.

Run 'bento' without arguments to launch the interactive TUI.
Or use commands directly: prepare, pack, pantry, taste.

Also available as 'b3o' alias.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Phase 4 will launch TUI here
		// For now, show help
		cmd.Help()
	},
}

func init() {
	// Phase 3 will add subcommands here
}
