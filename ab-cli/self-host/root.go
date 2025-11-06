package selfhost

import (
	initcmd "github.com/atomic-blend/backend/ab-cli/self-host/init_cmd"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-host",
		Short: "Commands for self-hosting atomic blend instances",
		Long: `The self-host command group provides various commands to manage
self-hosted atomic blend instances, including initialization, configuration,
and deployment.`,
	}
	cmd.AddCommand(initcmd.NewCommand())
	return cmd
}