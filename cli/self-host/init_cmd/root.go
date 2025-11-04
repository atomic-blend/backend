package initcmd

import (
	"fmt"
	"os"
	"path"

	"github.com/atomic-blend/backend/cli/config"
	"github.com/atomic-blend/backend/cli/self-host/init_cmd/ui"
	"github.com/atomic-blend/backend/cli/utils/yamlutils"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
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
	channel := getOrSetupChannel()
	log.Info().Str("channel", channel).Msg("Using channel")

	// TODO: check that config files (.env, docker-compose.yaml, ...) exists
	// If not, download the files from GitHub
	err := setupSelfHostedDirectory()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set up self-hosted directory")
	}

	log.Info().Msg("Self-hosted atomic blend instance initialized successfully")
}

func getOrSetupChannel() string {
	// Check for configured channel in multiple possible keys for compatibility.
	channel := config.CliConfig.Channel

	// check channel, and if not set, ask the user to choose one
	if channel == "" {
		// No configured channel — prompt the user using a Bubble Tea interactive selector.
		chosenCh := make(chan string, 1)
		m := ui.ChannelSelector{
			Choices:  []string{"stable", "rc"},
			Cursor:   1, // default to rc
			ChosenCh: chosenCh,
		}

		p := tea.NewProgram(m, tea.WithAltScreen())
		// Run the program. If it fails (non-tty or other error), fall back to default "rc".
		if _, err := p.Run(); err != nil {
			log.Error().Err(err).Msg("interactive selector failed, falling back to default 'rc'")
			chosen := "rc"
			viper.Set("channel", chosen)
			if err := yamlutils.PersistInConfigFile("channel", chosen); err != nil {
				log.Error().Err(err).Msg("failed to persist channel to config file")
			} else {
				log.Debug().Str("channel", chosen).Msg("wrote channel to config file")
			}
			return chosen
		}

		// Read chosen value from channel (model sends it before quitting)
		chosen := <-chosenCh

		log.Info().Msgf("User have selected update channel: %s", chosen)

		// Persist both the flattened and nested keys for compatibility.
		viper.Set("channel", chosen)

		// Persist into the config file, preserving comments/formatting.
		if err := yamlutils.PersistInConfigFile("channel", chosen); err != nil {
			log.Error().Err(err).Msg("failed to persist channel to config file")
		} else {
			log.Debug().Str("channel", chosen).Msg("wrote channel to config file")
		}
		return chosen
	} else {
		log.Debug().Str("channel", channel).Msg("Using configured update channel")
		return channel
	}
}

type FileSetup struct {
	LocalPath  string
	GitHubPath string
	Repository string
}

func setupSelfHostedDirectory() error {
	// TODO: Implement the directory setup logic
	files := []FileSetup{
		{LocalPath: ".env", GitHubPath: "docker/.env.example", Repository: "atomic-blend/backend"},
		{LocalPath: "docker-compose.yaml", GitHubPath: "docker/docker-compose.yaml", Repository: "atomic-blend/backend"},
		{LocalPath: "app-nginx.conf", GitHubPath: "docker/app-nginx.conf", Repository: "atomic-blend/backend"},
		{LocalPath: "nginx.conf", GitHubPath: "nginx.conf", Repository: "atomic-blend/backend"},
	}

	for _, file := range files {
		log.Info().Str("file", file.LocalPath).Msg("Setting up file in self-hosted directory")
		filename := path.Join(config.CliConfig.Directory, file.LocalPath)
		if _, err := os.Stat(filename); err != nil {
			log.Info().Str("file", filename).Msg("File does not exist, creating...")
			// Create or download the file
			GitHubPath := file.GitHubPath
			downloadURL := "https://raw.githubusercontent.com/" + file.Repository + "/main/" + GitHubPath
			log.Debug().Str("url", downloadURL).Msg("Downloading file from URL")
			resp, err := ui.GetResponse(downloadURL)
			if err != nil {
				fmt.Println("could not get response", err)
				os.Exit(1)
			}
			defer resp.Body.Close() // nolint:errcheck

			log.Debug().Int64("content-length", resp.ContentLength).Msg("Received response")
			log.Debug().Str("status", resp.Status).Msg("Response status")

			// Don't add TUI if the header doesn't include content size
			// it's impossible see progress without total
			if resp.ContentLength <= 0 {
				fmt.Println("can't parse content length, aborting download")
				os.Exit(1)
			}

			log.Debug().Str("filename", filename).Msg("Creating file")
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("could not create file:", err)
				os.Exit(1)
			}
			defer file.Close() // nolint:errcheck

			pw := &ui.ProgressWriter{
				Total:  int(resp.ContentLength),
				File:   file,
				Reader: resp.Body,
				OnProgress: func(ratio float64) {
					ui.P.Send(ui.ProgressMsg(ratio))
				},
			}

			m := ui.Model{
				Pw:       pw,
				Progress: progress.New(progress.WithDefaultGradient()),
			}
			// Start Bubble Tea
			ui.P = tea.NewProgram(m)

			// Start the download
			go pw.Start()

			if _, err := ui.P.Run(); err != nil {
				fmt.Println("error running program:", err)
				os.Exit(1)
			}
		} else {
			log.Info().Str("file", file.LocalPath).Msg("File already exists, skipping...")
		}

	}
	return nil
}
