package stopcmd

import (
	"os"
	"os/exec"

	"github.com/atomic-blend/backend/ab-cli/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for `self-host stop`.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the self-hosted atomic blend stack (docker compose down)",
		Run: func(cmd *cobra.Command, args []string) {
			dir := config.CliConfig.Directory
			if dir == "" {
				dir = "."
			}

			log.Info().Str("directory", dir).Msg("Running `docker compose down`")

			composeArgs := []string{"compose"}
			dockerProjectName, _ := cmd.Flags().GetString("docker-project")
			if dockerProjectName != "" {
				composeArgs = append(composeArgs, "-p", dockerProjectName)
			}
			composeArgs = append(composeArgs, "down")

			c := exec.Command("docker", composeArgs...)
			c.Dir = dir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			if err := c.Run(); err != nil {
				log.Fatal().Err(err).Msg("failed to run docker compose down")
			}

			log.Info().Msg("Services stopped successfully")
		},
	}

	cmd.Flags().String("docker-project", "ab-backend", "Docker project name (used as docker project in compose args)")
	return cmd
}
