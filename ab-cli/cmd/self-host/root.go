// Package selfhost implements the self-host command group for the ab-cli tool.
package selfhost

import (
	deletecmd "github.com/atomic-blend/backend/ab-cli/cmd/self-host/delete"
	initcmd "github.com/atomic-blend/backend/ab-cli/cmd/self-host/init"
	startcmd "github.com/atomic-blend/backend/ab-cli/cmd/self-host/start"
	stopcmd "github.com/atomic-blend/backend/ab-cli/cmd/self-host/stop"
	updatecmd "github.com/atomic-blend/backend/ab-cli/cmd/self-host/update"
	configcmd "github.com/atomic-blend/backend/ab-cli/cmd/self-host/config"
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
	cmd.PersistentFlags().String("docker-project", "ab-backend", "Docker project name (used as docker project in compose args)")

	cmd.AddCommand(initcmd.NewCommand())
	cmd.AddCommand(startcmd.NewCommand())
	cmd.AddCommand(stopcmd.NewCommand())
	cmd.AddCommand(deletecmd.NewCommand())
	cmd.AddCommand(updatecmd.NewCommand())
	cmd.AddCommand(configcmd.NewCommand())
	return cmd
}
