// Package main implements the wasabi command for secrets management.
//
// The wasabi command provides secure storage and retrieval of sensitive data
// like API tokens, passwords, and credentials using OS-native keychain systems.
package main

import (
	"fmt"
	"os"

	"github.com/Develonaut/bento/pkg/wasabi"
	"github.com/spf13/cobra"
)

var wasabiCmd = &cobra.Command{
	Use:   "wasabi",
	Short: "🟢 Manage secrets (handle the spicy stuff)",
	Long: `Manage secrets securely using OS-native keychain.

Wasabi stores sensitive data like API tokens and passwords in your
system's keychain (macOS Keychain, Windows Credential Manager, or
Linux Secret Service) so they never appear in bento files or git history.

Secrets are accessed in bentos using {{SECRETS.X}} syntax:
  "Authorization": "Bearer {{SECRETS.FIGMA_TOKEN}}"

Commands:
  set    - Store a secret
  get    - Retrieve a secret
  list   - List all secret keys
  delete - Remove a secret`,
}

var wasabiSetCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Store a secret in the keychain",
	Long: `Store a secret securely in the OS keychain.

The secret will be encrypted and stored using your system's
native keychain service. It can then be used in bentos via
the {{SECRETS.KEY}} syntax.

Examples:
  bento wasabi set FIGMA_TOKEN figd_abc123xyz
  bento wasabi set API_KEY sk-1234567890abcdef`,
	Args: cobra.ExactArgs(2),
	RunE: runWasabiSet,
}

var wasabiGetCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Retrieve a secret from the keychain",
	Long: `Retrieve a secret from the OS keychain.

SECURITY WARNING: This prints the secret to stdout!
Only use this for debugging - secrets are automatically
resolved in bentos without needing to print them.

Examples:
  bento wasabi get FIGMA_TOKEN`,
	Args: cobra.ExactArgs(1),
	RunE: runWasabiGet,
}

var wasabiListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secret keys (not values)",
	Long: `List all secret keys stored in the keychain.

This shows only the KEY names, not the secret values.
Use 'bento wasabi get KEY' to retrieve a specific secret.

Example:
  bento wasabi list`,
	Args: cobra.NoArgs,
	RunE: runWasabiList,
}

var wasabiDeleteCmd = &cobra.Command{
	Use:   "delete KEY",
	Short: "Remove a secret from the keychain",
	Long: `Delete a secret from the OS keychain.

This permanently removes the secret. You'll need to
use 'bento wasabi set' to store it again.

Examples:
  bento wasabi delete FIGMA_TOKEN
  bento wasabi delete OLD_API_KEY`,
	Args: cobra.ExactArgs(1),
	RunE: runWasabiDelete,
}

func init() {
	wasabiCmd.AddCommand(wasabiSetCmd)
	wasabiCmd.AddCommand(wasabiGetCmd)
	wasabiCmd.AddCommand(wasabiListCmd)
	wasabiCmd.AddCommand(wasabiDeleteCmd)
}

// runWasabiSet stores a secret in the keychain.
func runWasabiSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Validate key format (uppercase letters, numbers, underscores)
	if !isValidSecretKey(key) {
		printError("Invalid secret key. Use uppercase letters, numbers, and underscores only.")
		printInfo("Example: FIGMA_TOKEN, API_KEY_V2")
		return fmt.Errorf("invalid secret key format")
	}

	mgr, err := wasabi.NewManager()
	if err != nil {
		printError(fmt.Sprintf("Failed to initialize secrets manager: %v", err))
		return err
	}

	if err := mgr.Set(key, value); err != nil {
		printError(fmt.Sprintf("Failed to store secret: %v", err))
		return err
	}

	printSuccess(fmt.Sprintf("Secret '%s' stored securely", key))
	printInfo("Use {{SECRETS." + key + "}} in your bentos to access it")
	return nil
}

// runWasabiGet retrieves a secret from the keychain.
func runWasabiGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	mgr, err := wasabi.NewManager()
	if err != nil {
		printError(fmt.Sprintf("Failed to initialize secrets manager: %v", err))
		return err
	}

	value, err := mgr.Get(key)
	if err != nil {
		printError(fmt.Sprintf("Failed to retrieve secret: %v", err))
		return err
	}

	// Print warning about security
	printWarning("⚠️  Security Warning: Printing secret to stdout!")
	fmt.Println(value)
	return nil
}

// runWasabiList lists all secret keys.
func runWasabiList(cmd *cobra.Command, args []string) error {
	mgr, err := wasabi.NewManager()
	if err != nil {
		printError(fmt.Sprintf("Failed to initialize secrets manager: %v", err))
		return err
	}

	keys, err := mgr.List()
	if err != nil {
		printError(fmt.Sprintf("Failed to list secrets: %v", err))
		return err
	}

	if len(keys) == 0 {
		printInfo("No secrets stored yet")
		printInfo("Use 'bento wasabi set KEY VALUE' to store a secret")
		return nil
	}

	printSuccess(fmt.Sprintf("Found %d secret(s):", len(keys)))
	for _, key := range keys {
		fmt.Printf("  • %s (use {{SECRETS.%s}} in bentos)\n", key, key)
	}

	return nil
}

// runWasabiDelete removes a secret from the keychain.
func runWasabiDelete(cmd *cobra.Command, args []string) error {
	key := args[0]

	mgr, err := wasabi.NewManager()
	if err != nil {
		printError(fmt.Sprintf("Failed to initialize secrets manager: %v", err))
		return err
	}

	if err := mgr.Delete(key); err != nil {
		printError(fmt.Sprintf("Failed to delete secret: %v", err))
		return err
	}

	printSuccess(fmt.Sprintf("Secret '%s' deleted", key))
	return nil
}

// isValidSecretKey validates secret key format.
// Keys must be uppercase with underscores (e.g., FIGMA_TOKEN, API_KEY).
func isValidSecretKey(key string) bool {
	if key == "" {
		return false
	}

	// Must start with uppercase letter or underscore
	if !((key[0] >= 'A' && key[0] <= 'Z') || key[0] == '_') {
		return false
	}

	// Rest must be uppercase letters, numbers, or underscores
	for _, ch := range key {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}

	return true
}

// printWarning prints a warning message.
func printWarning(msg string) {
	fmt.Fprintf(os.Stderr, "⚠️  %s\n", msg)
}
