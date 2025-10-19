// Package main implements the recipe command for viewing bento documentation.
//
// The recipe command uses charm glow to render markdown documentation
// in a beautiful, readable format directly in the terminal.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var recipeCmd = &cobra.Command{
	Use:   "recipe [doc-name]",
	Short: "📖 View bento documentation",
	Long: `View bento documentation in beautiful markdown format.

Uses charm glow to render documentation files in the terminal.

Available docs:
  readme          - Project README
  overview        - Project overview and architecture
  packages        - Package naming conventions
  principles      - Bento Box Principle
  emojis          - Approved emoji list
  charm           - Charm stack integration
  nodes           - Complete node inventory
  standards       - Go standards review
  status-words    - Status word guidelines

Examples:
  bento recipe readme
  bento recipe principles
  bento recipe nodes`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRecipe,
}

// docMap maps friendly names to file paths.
var docMap = map[string]string{
	"readme":       "README.md",
	"overview":     ".claude/README.md",
	"packages":     ".claude/PACKAGE_NAMING.md",
	"principles":   ".claude/BENTO_BOX_PRINCIPLE.md",
	"emojis":       ".claude/EMOJIS.md",
	"charm":        ".claude/CHARM_STACK.md",
	"nodes":        ".claude/COMPLETE_NODE_INVENTORY.md",
	"standards":    ".claude/GO_STANDARDS_REVIEW.md",
	"status-words": ".claude/STATUS_WORDS.md",
}

// runRecipe executes the recipe command logic.
func runRecipe(cmd *cobra.Command, args []string) error {
	// Default to README if no doc specified
	docName := "readme"
	if len(args) > 0 {
		docName = args[0]
	}

	// Look up the file path
	filePath, ok := docMap[docName]
	if !ok {
		return fmt.Errorf("unknown doc: %s\nRun 'bento recipe --help' to see available docs", docName)
	}

	// Check if glow is installed
	if !isGlowInstalled() {
		return fmt.Errorf("glow is not installed\n\nInstall with:\n  brew install glow\n\nOr visit: https://github.com/charmbracelet/glow")
	}

	// Get absolute path to the doc
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("doc file not found: %s", absPath)
	}

	// Run glow
	return viewWithGlow(absPath)
}

// isGlowInstalled checks if glow is available.
func isGlowInstalled() bool {
	_, err := exec.LookPath("glow")
	return err == nil
}

// viewWithGlow renders the markdown file with glow.
func viewWithGlow(path string) error {
	glowCmd := exec.Command("glow", path)
	glowCmd.Stdin = os.Stdin
	glowCmd.Stdout = os.Stdout
	glowCmd.Stderr = os.Stderr
	return glowCmd.Run()
}
