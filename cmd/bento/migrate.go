package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"bento/pkg/neta"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [file or directory]",
	Short: "Add version field to existing bento files",
	Long: `Migrate adds the version field to existing .bento.yaml files.

This command:
- Reads existing .bento.yaml files
- Adds version: "1.0" if missing
- Preserves all existing fields
- Can process single files or entire directories

Examples:
  bento migrate workflow.bento.yaml
  bento migrate examples/
  bento migrate .`,
	Args: cobra.ExactArgs(1),
	RunE: runMigrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrate(cmd *cobra.Command, args []string) error {
	path := args[0]

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return migrateDir(path)
	}
	return migrateFile(path)
}

func migrateFile(path string) error {
	def, err := loadAndParseFile(path)
	if err != nil {
		return err
	}

	if !needsMigration(&def) {
		fmt.Printf("✓ %s (already versioned: %s)\n", path, def.Version)
		return nil
	}

	addVersion(&def)

	if err := saveDefinition(path, def); err != nil {
		return err
	}

	fmt.Printf("✓ %s (added version: %s)\n", path, neta.CurrentVersion)
	return nil
}

func loadAndParseFile(path string) (neta.Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return neta.Definition{}, err
	}

	var def neta.Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return neta.Definition{}, fmt.Errorf("parse failed: %w", err)
	}

	return def, nil
}

func saveDefinition(path string, def neta.Definition) error {
	data, err := yaml.Marshal(def)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}

// needsMigration checks if def or any child nodes are missing versions.
func needsMigration(def *neta.Definition) bool {
	if def.Version == "" {
		return true
	}

	for i := range def.Nodes {
		if needsMigration(&def.Nodes[i]) {
			return true
		}
	}

	return false
}

// addVersion recursively adds version to definition and all child nodes.
func addVersion(def *neta.Definition) {
	if def.Version == "" {
		def.Version = neta.CurrentVersion
	}

	for i := range def.Nodes {
		addVersion(&def.Nodes[i])
	}
}

func migrateDir(dir string) error {
	matches, err := findBentoFiles(dir)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		fmt.Printf("No .bento.yaml files found in %s\n", dir)
		return nil
	}

	success, failed := processBentoFiles(matches)
	return reportMigrationResults(success, failed)
}

func findBentoFiles(dir string) ([]string, error) {
	pattern := filepath.Join(dir, "*.bento.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found %d .bento.yaml files in %s\n", len(matches), dir)
	return matches, nil
}

func processBentoFiles(files []string) (success, failed int) {
	for _, path := range files {
		if err := migrateFile(path); err != nil {
			fmt.Printf("✗ %s: %v\n", path, err)
			failed++
		} else {
			success++
		}
	}
	return success, failed
}

func reportMigrationResults(success, failed int) error {
	fmt.Printf("\nMigration complete: %d succeeded, %d failed\n", success, failed)

	if failed > 0 {
		return fmt.Errorf("%d files failed migration", failed)
	}

	return nil
}
