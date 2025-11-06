// Package configcmd implements the `self-host config` subcommand.
package configcmd

import (
	"os"
	"os/exec"
	"path"

	configpkg "github.com/atomic-blend/backend/ab-cli/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Open the .env configuration file in vim",
		Long:  "Open the `.env` file of the self-hosted directory in the user's editor (vim).",
		Run: func(cmd *cobra.Command, args []string) {
			dir := configpkg.CliConfig.Directory
			log.Info().Str("directory", dir).Msg("Opening .env config in editor")

			editor := exec.Command("vim", path.Join(dir, ".env"))
			editor.Stdin = os.Stdin
			editor.Stdout = os.Stdout
			editor.Stderr = os.Stderr

			if err := editor.Run(); err != nil {
				log.Fatal().Err(err).Msg("failed to open .env in editor")
			}
		},
	}

	return cmd
}
