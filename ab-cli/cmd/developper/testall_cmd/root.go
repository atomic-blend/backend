package testall

import (
	"context"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testall",
		Short: "Run all tests for the atomic blend backend",
		Long: `The testall command runs all tests and linters for the atomic blend backend services.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Start the Bubble Tea TUI which will run lint and tests and display results
			ctx := context.Background()
			return StartTUI(ctx)
		},
	}
	return cmd
}
