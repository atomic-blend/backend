package developper

import (
	sendemailcmd "github.com/atomic-blend/backend/ab-cli/cmd/developper/send_email_cmd"
	test "github.com/atomic-blend/backend/ab-cli/cmd/developper/test_cmd"
	testall "github.com/atomic-blend/backend/ab-cli/cmd/developper/testall_cmd"
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
	cmd.AddCommand(sendemailcmd.NewCommand())
	return cmd
}
