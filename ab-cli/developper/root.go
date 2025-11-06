package developper

import (
	"github.com/atomic-blend/backend/ab-cli/developper/test"
	"github.com/atomic-blend/backend/ab-cli/developper/testall"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "developper",
		Short: "Commands for developper atomic blend instances",
		Long:  `The developper command group provides various commands to various actions while developping the platform.`,
	}
	cmd.AddCommand(testall.NewCommand())
	cmd.AddCommand(test.NewCommand())
	return cmd
}
