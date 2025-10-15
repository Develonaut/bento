package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bento/pkg/omise"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:     "bento",
	Short:   "🍱 Bento - Organized bento orchestration",
	Version: "0.1.0",
	Long: `Bento is a Go-based CLI orchestration tool.

Run 'bento' without arguments to launch the interactive TUI (Phase 4).
Or use commands directly for scripting and automation.

Available commands:
  prepare - Validate a .bento.yaml file
  pack    - Execute a bento
  pantry  - List/search available neta types
  taste   - Dry run a bento

Also available as 'b3o' alias.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := omise.Launch(); err != nil {
			fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bento.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Cross-platform: Use os.UserHomeDir() (Go 1.12+)
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding home: %v\n", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bento")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
