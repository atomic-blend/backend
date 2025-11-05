package initcmd

import (
	"fmt"
	"os"
	"path"

	"github.com/atomic-blend/backend/cli/config"
	bulkfiledownloader "github.com/atomic-blend/backend/cli/ui/bulk_file_downloader"
	channelselector "github.com/atomic-blend/backend/cli/ui/channel_selector"
	"github.com/atomic-blend/backend/cli/ui/types/file_types"
	envfilesutils "github.com/atomic-blend/backend/cli/utils/env_files_utils"
	"github.com/atomic-blend/backend/cli/utils/yamlutils"
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

	log.Debug().Msg("Reading configuration from local files")
	platformComponents, err := envfilesutils.GetPlatformConfig(path.Join(config.CliConfig.Directory, ".env"), path.Join(config.CliConfig.Directory, "docker-compose.yaml"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get platform config")
	}

	log.Info().Msg("Platform components and versions:")
	for _, comp := range platformComponents {
		log.Info().Str("component", comp.Name).Str("version", comp.Version)
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
		m := channelselector.ChannelSelector{
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

func setupSelfHostedDirectory() error {
	// TODO: Implement the directory setup logic
	files := []filetypes.DownloadableGithubFile{
		{LocalPath: ".env", GitHubPath: "docker/.env.example", Repository: "atomic-blend/backend"},
		{LocalPath: "docker-compose.yaml", GitHubPath: "docker/docker-compose.yaml", Repository: "atomic-blend/backend"},
		{LocalPath: "app-nginx.conf", GitHubPath: "docker/app-nginx.conf", Repository: "atomic-blend/backend"},
		{LocalPath: "nginx.conf", GitHubPath: "nginx.conf", Repository: "atomic-blend/backend"},
	}

	var toDownload []filetypes.DownloadableFile
	for _, file := range files {
		filename := path.Join(config.CliConfig.Directory, file.LocalPath)
		if _, err := os.Stat(filename); err == nil {
			log.Info().Str("file", file.LocalPath).Msg("File already exists, skipping...")
			continue
		}

		log.Debug().Str("file", file.LocalPath).Msg("File does not exist, scheduling download")

		// Ensure target directory exists
		destDir := path.Dir(filename)
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			fmt.Println("could not create directory:", err)
			os.Exit(1)
		}

		downloadURL := "https://raw.githubusercontent.com/" + file.Repository + "/main/" + file.GitHubPath
		toDownload = append(toDownload, filetypes.DownloadableFile{
			URL:       downloadURL,
			LocalPath: filename,
		})
	}

	if len(toDownload) == 0 {
		return nil
	}

	if err := bulkfiledownloader.DownloadBulk(toDownload); err != nil {
		fmt.Println("error downloading files:", err)
		os.Exit(1)
	}

	log.Debug().Msg("All necessary files downloaded successfully")
	return nil
}
