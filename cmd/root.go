package cmd

import (
	"fmt"
	"os"
	"reshell/pkg/tui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "reshell",
	Short: "reshell: The Reproducible Developer Terminal Manager",
	Long: `reshell is a CLI utility and TUI dashboard designed to manage snippets,
aliases, custom shell functions, environment variables, git settings, script logs,
and project templates in a portable, reproducible environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Launch the TUI dashboard when run without arguments
		return tui.Start()
	},
}

// Execute parses commands and handles errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
