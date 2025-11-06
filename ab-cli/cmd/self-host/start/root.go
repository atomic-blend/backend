package startcmd

import (
	"os"
	"os/exec"

	"github.com/atomic-blend/backend/ab-cli/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for `self-host start`.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the self-hosted atomic blend stack (docker compose up -d)",
		Run: func(cmd *cobra.Command, args []string) {
			dir := config.CliConfig.Directory
			if dir == "" {
				dir = "."
			}

			// check that a docker-compose.yaml file exists in the directory
			if _, err := os.Stat(dir + "/docker-compose.yaml"); os.IsNotExist(err) {
				log.Fatal().Str("directory", dir).Msg("docker-compose.yaml file not found in the specified directory")
			}

			log.Info().Str("directory", dir).Msg("Running `docker compose up -d`")

			composeArgs := []string{"compose"}
			dockerProjectName, _ := cmd.Flags().GetString("docker-project")
			if dockerProjectName != "" {
				composeArgs = append(composeArgs, "-p", dockerProjectName)
			}
			composeArgs = append(composeArgs, "up", "-d")

			c := exec.Command("docker", composeArgs...)
			c.Dir = dir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			if err := c.Run(); err != nil {
				log.Fatal().Err(err).Msg("failed to run docker compose up -d")
			}

			log.Info().Msg("Services started successfully")
		},
	}
	cmd.Flags().String("docker-project", "ab-backend", "Docker project name (used as docker project in compose args)")
	return cmd
}
