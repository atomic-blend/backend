// Package updatecmd implements the `self-host update` subcommand.
package updatecmd

import (
	"os"

	"github.com/atomic-blend/backend/ab-cli/config"
	platformcomponentupdater "github.com/atomic-blend/backend/ab-cli/ui/platform_component_updater"
	envfilesutils "github.com/atomic-blend/backend/ab-cli/utils/env_files_utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for `self-host update`.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Check and apply updates for platform components in .env",
		Run: func(cmd *cobra.Command, args []string) {
			// check in directory that .env and docker-compose.yaml files exist
			if _, err := os.Stat(config.CliConfig.Directory + "/.env"); os.IsNotExist(err) {
				log.Fatal().Str("directory", config.CliConfig.Directory).Msg(".env file not found in the specified directory")
			}
			if _, err := os.Stat(config.CliConfig.Directory + "/docker-compose.yaml"); os.IsNotExist(err) {
				log.Fatal().Str("directory", config.CliConfig.Directory).Msg("docker-compose.yaml file not found in the specified directory")
			}
			// Read platform components from local files
			platformComponents, err := envfilesutils.GetPlatformConfig(config.CliConfig.Directory+"/.env", config.CliConfig.Directory+"/docker-compose.yaml")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get platform config")
			}

			// Determine RC preference: explicit flag wins, otherwise use configured channel
			rcFlag, _ := cmd.Flags().GetBool("rc")
			isRC := rcFlag || config.CliConfig.Channel == "rc"

			log.Info().Bool("rc", isRC).Msg("Starting platform component updater")
			if err := platformcomponentupdater.GetUpdates(platformComponents, &isRC); err != nil {
				log.Error().Err(err).Msg("platform component updater failed")
			}
		},
	}

	cmd.Flags().Bool("rc", false, "Prefer RC (release-candidate) images when available and newer")
	return cmd
}
