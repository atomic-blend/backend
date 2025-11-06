package deletecmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/atomic-blend/backend/ab-cli/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewCommand returns the cobra command for `self-host delete`.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete the self-hosted atomic blend stack and downloaded files (including volumes)",
		Run: func(cmd *cobra.Command, args []string) {
			dir := config.CliConfig.Directory
			if dir == "" {
				dir = "."
			}

			dockerProjectName, _ := cmd.Flags().GetString("docker-project")

			// Confirmation prompt
			reader := bufio.NewReader(os.Stdin)
			fmt.Println()
			fmt.Println("docker project " + dockerProjectName)
			fmt.Println()
			errorColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
			fmt.Println(errorColor.Render("WARNING: This will delete ALL data associated with the self-hosted atomic blend instance, including Docker volumes!"))
			fmt.Println()
			fmt.Print("Are you sure you want to delete everything including data? (y/N): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal().Err(err).Msg("failed to read user input")
			}
			input = strings.TrimSpace(input)
			if input == "" || !(strings.ToLower(input) == "y" || strings.ToLower(input) == "yes") {
				fmt.Println("")
				log.Info().Msg("Delete cancelled by user")
				return
			}

			log.Info().Str("directory", dir).Msg("Running `docker compose down -v`")

			composeArgs := []string{"compose"}
			if dockerProjectName != "" {
				composeArgs = append(composeArgs, "-p", dockerProjectName)
			}
			composeArgs = append(composeArgs, "down", "-v")

			c := exec.Command("docker", composeArgs...)
			c.Dir = dir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin

			if err := c.Run(); err != nil {
				log.Fatal().Err(err).Msg("failed to run docker compose down -v")
			}

			// Remove files downloaded by `init` command
			files := []string{".env", "docker-compose.yaml", "app-nginx.conf", "nginx.conf"}
			for _, f := range files {
				filename := path.Join(dir, f)
				if _, err := os.Stat(filename); err == nil {
					if err := os.Remove(filename); err != nil {
						log.Error().Err(err).Str("file", filename).Msg("failed to remove file")
					} else {
						log.Info().Str("file", filename).Msg("Removed file")
					}
				} else {
					log.Debug().Str("file", filename).Msg("File not found, skipping")
				}
			}

			log.Info().Msg("Deletion completed")
		},
	}
	return cmd
}
