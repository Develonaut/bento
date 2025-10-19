// Package main implements the menu command for listing bentos.
//
// The menu command scans a directory for .bento.json files and displays
// them in a user-friendly format with names and metadata.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Develonaut/bento/pkg/neta"
	"github.com/spf13/cobra"
)

var (
	recursiveFlag bool
	jsonFlag      bool
)

var menuCmd = &cobra.Command{
	Use:   "menu [directory]",
	Short: "🍱 List available bentos",
	Long: `List all available bentos in a directory.

Like a restaurant menu, this shows you all the bentos you can taste.

Examples:
  bento menu
  bento menu ~/workflows
  bento menu ~/workflows --recursive`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMenu,
}

func init() {
	menuCmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", false, "Search subdirectories")
	menuCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output as JSON")
}

// runMenu executes the menu command logic.
func runMenu(cmd *cobra.Command, args []string) error {
	dir := getDir(args)
	bentos, err := findBentos(dir)
	if err != nil {
		printError(fmt.Sprintf("Failed to scan directory: %v", err))
		return err
	}

	displayBentos(bentos)
	return nil
}

// getDir returns the directory to scan (default: current).
func getDir(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "."
}

// displayBentos displays the list of found bentos.
func displayBentos(bentos []bentoInfo) {
	if len(bentos) == 0 {
		fmt.Println("🍱 No bentos found")
		return
	}

	printInfo("Available Bentos:\n")
	for _, bento := range bentos {
		printBento(bento)
	}
	fmt.Printf("\n%d bentos found\n", len(bentos))
}

// bentoInfo contains metadata about a bento file.
type bentoInfo struct {
	Path     string
	FileName string
	Name     string
	NumNodes int
}

// findBentos finds all .bento.json files in the given directory.
func findBentos(dir string) ([]bentoInfo, error) {
	var bentos []bentoInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldSkipDir(info, path, dir) {
			return filepath.SkipDir
		}

		if isBentoFile(info) {
			bentos = append(bentos, extractBentoInfo(path))
		}

		return nil
	})

	return bentos, err
}

// shouldSkipDir checks if directory should be skipped.
func shouldSkipDir(info os.FileInfo, path, dir string) bool {
	return info.IsDir() && path != dir && !recursiveFlag
}

// isBentoFile checks if file is a bento JSON file.
func isBentoFile(info os.FileInfo) bool {
	return !info.IsDir() && strings.HasSuffix(info.Name(), ".bento.json")
}

// extractBentoInfo extracts metadata from a bento file.
func extractBentoInfo(path string) bentoInfo {
	info := bentoInfo{
		Path:     path,
		FileName: filepath.Base(path),
	}

	// Try to load the bento to get name and node count
	def, err := loadBento(path)
	if err == nil {
		info.Name = def.Name
		info.NumNodes = countNodes(def)
	}

	return info
}

// countNodes counts the number of nodes in a bento.
func countNodes(def *neta.Definition) int {
	return len(def.Nodes)
}

// printBento prints a single bento entry.
func printBento(b bentoInfo) {
	fmt.Printf("  %s\n", b.FileName)
	if b.Name != "" {
		fmt.Printf("    %s\n", b.Name)
		fmt.Printf("    %d neta\n", b.NumNodes)
	}
	fmt.Println()
}
