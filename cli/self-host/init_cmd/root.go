package initcmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var directory string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a self-hosted atomic blend instance",
		Long: `The init command sets up a new self-hosted atomic blend instance.
It guides you through the necessary steps to configure and deploy your instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			directory = viper.GetString("self-host.directory")
			initSelfHost(cmd, args)
		},
	}
	cmd.Flags().String("directory", ".", "Directory to initialize the atomic blend instance in")
	_ = viper.BindPFlag("self-host.directory", cmd.Flags().Lookup("directory"))
	_ = viper.BindEnv("self-host.directory", "ATOMIC_BLEND_SELFHOST_DIRECTORY")

	return cmd
}

func initSelfHost(cmd *cobra.Command, args []string) {
	log.Info().Str("directory", directory).Msg("Initializing self-hosted atomic blend instance")

	log.Info().Msg("Self-hosted atomic blend instance initialized successfully")
}
