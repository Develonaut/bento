package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"bento/pkg/neta"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run yaml2json.go <input.yaml> <output.json>")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("  go run yaml2json.go input.bento.yaml output.bento.json")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	fmt.Printf("Converting %s → %s\n", inputPath, outputPath)

	// Read YAML file
	yamlData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("❌ Error reading YAML file: %v\n", err)
		os.Exit(1)
	}

	// Parse YAML into neta.Definition
	var def neta.Definition
	if err := yaml.Unmarshal(yamlData, &def); err != nil {
		fmt.Printf("❌ Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Parsed YAML successfully\n")
	fmt.Printf("   Name: %s\n", def.Name)
	fmt.Printf("   Type: %s\n", def.Type)
	fmt.Printf("   Nodes: %d\n", len(def.Nodes))

	// Convert to JSON with 2-space indentation
	jsonData, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error marshaling to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write JSON file
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		fmt.Printf("❌ Error writing JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Wrote JSON to %s\n", outputPath)
	fmt.Printf("✅ Conversion complete!\n")
}
