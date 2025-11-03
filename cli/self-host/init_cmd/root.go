package initcmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var directory string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a self-hosted atomic blend instance",
		Long: `The init command sets up a new self-hosted atomic blend instance.
It guides you through the necessary steps to configure and deploy your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Implementation of the init command goes here
			log.Debug().Str("directory", directory).Msg("Initializing self-hosted atomic blend instance")
		},
	}
	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Specify the directory to initialize the instance in")
	return cmd
}
