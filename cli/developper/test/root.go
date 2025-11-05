package test

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests for the atomic blend platform",
		Long: `The test command group provides various commands to manage
testing of atomic blend instances, including initialization, configuration,
and deployment.`,
		Run: func(cmd *cobra.Command, args []string) {
			runTests()
		},
	}
	return cmd
}

func runTests() {
	// TODO: throw err if dev mode is not enabled
	// TODO: run tests all (args 1 = all)
	// TODO: run tests for a specific service (args 1 = service name)
}
