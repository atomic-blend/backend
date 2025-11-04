package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ChannelSelector is a tiny Bubble Tea model used to select between update channels.
type ChannelSelector struct {
	Choices  []string
	Cursor   int
	ChosenCh chan string
}

func (m ChannelSelector) Init() tea.Cmd { return nil }

func (m ChannelSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter":
			// send the chosen value and quit
			m.ChosenCh <- m.Choices[m.Cursor]
			return m, tea.Quit
		case "ctrl+c", "q":
			// cancel/quit -> default to rc
			m.ChosenCh <- "rc"
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ChannelSelector) View() string {
	var b strings.Builder
	b.WriteString("Choose update channel (use ↑/↓ and Enter):\n\n")
	for i, c := range m.Choices {
		cursor := "  "
		if m.Cursor == i {
			cursor = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, c))
	}
	b.WriteString("\nPress q or Ctrl+C to cancel (defaults to rc)\n")
	return b.String()
}
