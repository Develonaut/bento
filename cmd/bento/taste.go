package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tasteCmd = &cobra.Command{
	Use:   "taste [file.bento.yaml]",
	Short: "Dry run a bento (alias for prepare)",
	Long: `Taste validates a bento without executing it.

This is an alias for 'bento prepare' with more verbose output.`,
	Args: cobra.ExactArgs(1),
	RunE: runTaste,
}

func init() {
	rootCmd.AddCommand(tasteCmd)
}

func runTaste(cmd *cobra.Command, args []string) error {
	fmt.Println("🍱 Tasting your bento...")

	if err := runPrepare(cmd, args); err != nil {
		fmt.Println("\n❌ This bento doesn't taste right!")
		return err
	}

	fmt.Println("\n✨ Delicious! Ready to pack.")
	return nil
}
