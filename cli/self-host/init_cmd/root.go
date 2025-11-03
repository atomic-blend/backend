package initcmd

import (
	"fmt"

	"github.com/atomic-blend/backend/cli/config"
	"github.com/atomic-blend/backend/cli/utils/yamlutils"
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
			directory = viper.GetString("directory")
			initSelfHost(cmd, args)
		},
	}

	return cmd
}

func initSelfHost(cmd *cobra.Command, args []string) {
	log.Info().Str("directory", directory).Msg("Initializing self-hosted atomic blend instance")
	// Check for configured channel in multiple possible keys for compatibility.
	channel := config.CliConfig.Channel

	if channel == "" {
		// No configured channel — prompt the user to choose one and persist it.
		fmt.Println("No update channel configured. Choose one:")
		fmt.Println("  1) stable")
		fmt.Println("  2) rc")
		fmt.Print("Enter choice [1-2] (default 2): ")

		var choice string
		_, err := fmt.Scanln(&choice)
		if err != nil {
			// If Scanln fails (e.g., EOF), default to rc
			choice = "2"
		}

		chosen := "rc"
		if choice == "1" || choice == "stable" {
			chosen = "stable"
		}

		// Persist both the flattened and nested keys for compatibility.
		viper.Set("channel", chosen)

		// Persist into the config file, preserving comments/formatting.
		if err := yamlutils.PersistInConfigFile("channel", chosen); err != nil {
			log.Error().Err(err).Msg("failed to persist channel to config file")
		} else {
			log.Debug().Str("channel", chosen).Msg("wrote channel to config file")
		}
	} else {
		log.Debug().Str("channel", channel).Msg("Using configured update channel")
	}
	log.Info().Msg("Self-hosted atomic blend instance initialized successfully")
}
