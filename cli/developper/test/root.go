package test

import (
	"context"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests for the atomic blend platform",
		Long: `The test command group provides various commands to manage
testing of atomic blend instances, including initialization, configuration,
and deployment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Start the Bubble Tea TUI which will run lint and tests and display results
			ctx := context.Background()
			return StartTUI(ctx)
		},
	}
	return cmd
}
