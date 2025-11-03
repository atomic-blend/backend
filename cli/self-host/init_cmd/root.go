package initcmd

import (
	"fmt"
	"strings"

	"github.com/atomic-blend/backend/cli/config"
	"github.com/atomic-blend/backend/cli/utils/yamlutils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var directory string

// selectorModel is a tiny Bubble Tea model used to select between update channels.
type selectorModel struct {
	choices  []string
	cursor   int
	chosenCh chan string
}

func (m selectorModel) Init() tea.Cmd { return nil }

func (m selectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			// send the chosen value and quit
			m.chosenCh <- m.choices[m.cursor]
			return m, tea.Quit
		case "ctrl+c", "q":
			// cancel/quit -> default to rc
			m.chosenCh <- "rc"
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectorModel) View() string {
	var b strings.Builder
	b.WriteString("Choose update channel (use ↑/↓ and Enter):\n\n")
	for i, c := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, c))
	}
	b.WriteString("\nPress q or Ctrl+C to cancel (defaults to rc)\n")
	return b.String()
}

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

	log.Info().Msg("Self-hosted atomic blend instance initialized successfully")
}

func getOrSetupChannel() string {
	// Check for configured channel in multiple possible keys for compatibility.
	channel := config.CliConfig.Channel

	// check channel, and if not set, ask the user to choose one
	if channel == "" {
		// No configured channel — prompt the user using a Bubble Tea interactive selector.
		chosenCh := make(chan string, 1)
		m := selectorModel{
			choices:  []string{"stable", "rc"},
			cursor:   1, // default to rc
			chosenCh: chosenCh,
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
