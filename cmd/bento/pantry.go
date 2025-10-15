package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"bento/pkg/itamae"
	"bento/pkg/neta/conditional"
	"bento/pkg/neta/group"
	"bento/pkg/neta/http"
	"bento/pkg/neta/loop"
	"bento/pkg/neta/transform"
	"bento/pkg/pantry"
)

var pantryCmd = &cobra.Command{
	Use:   "pantry [search]",
	Short: "List or search available neta types",
	Long: `Pantry shows all registered node types.

Optionally provide a search term to filter results.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPantry,
}

func init() {
	rootCmd.AddCommand(pantryCmd)
}

func runPantry(cmd *cobra.Command, args []string) error {
	p := initializePantry()

	types := p.List()
	sort.Strings(types)

	if len(args) > 0 {
		types = filterTypes(types, args[0])
	}

	printTypes(types)
	return nil
}

func filterTypes(types []string, search string) []string {
	search = strings.ToLower(search)
	filtered := []string{}

	for _, t := range types {
		if strings.Contains(strings.ToLower(t), search) {
			filtered = append(filtered, t)
		}
	}

	return filtered
}

func printTypes(types []string) {
	fmt.Printf("🍱 Available neta types (%d):\n\n", len(types))

	for _, t := range types {
		fmt.Printf("  • %s\n", t)
	}

	if len(types) == 0 {
		fmt.Println("  (no types found)")
	}
}

// initializePantry creates and populates the pantry with all node types.
func initializePantry() *pantry.Pantry {
	p := pantry.New()

	// We need to create a dummy itamae for nodes that require an executor.
	// This is safe because pantry command only lists types, never executes them.
	dummyItamae := itamae.New(p)

	_ = p.Register("http", http.New())
	_ = p.Register("jq", transform.NewJQ())
	_ = p.Register("sequence", group.NewSequence(dummyItamae))
	_ = p.Register("parallel", group.NewParallel(dummyItamae))
	_ = p.Register("if", conditional.NewIf(dummyItamae))
	_ = p.Register("for", loop.NewFor(dummyItamae))

	return p
}
