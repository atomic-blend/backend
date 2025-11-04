package ui

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var P *tea.Program

type ProgressWriter struct {
	Total      int
	Downloaded int
	File       *os.File
	Reader     io.Reader
	OnProgress func(float64)
}

func (pw *ProgressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw.File, io.TeeReader(pw.Reader, pw))
	if err != nil {
		P.Send(ProgressErrMsg{err})
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	pw.Downloaded += len(p)
	if pw.Total > 0 && pw.OnProgress != nil {
		pw.OnProgress(float64(pw.Downloaded) / float64(pw.Total))
	}
	return len(p), nil
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

type ProgressMsg float64

type ProgressErrMsg struct{ err error }

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

type Model struct {
	Pw       *ProgressWriter
	Progress progress.Model
	Err      error
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.Progress.Width = msg.Width - padding*2 - 4
		if m.Progress.Width > maxWidth {
			m.Progress.Width = maxWidth
		}
		return m, nil

	case ProgressErrMsg:
		m.Err = msg.err
		return m, tea.Quit

	case ProgressMsg:
		var cmds []tea.Cmd

		if msg >= 1.0 {
			cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
		}

		cmds = append(cmds, m.Progress.SetPercent(float64(msg)))
		return m, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.Progress.Update(msg)
		m.Progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m Model) View() string {
	if m.Err != nil {
		return "Error downloading: " + m.Err.Error() + "\n"
	}

	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.Progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func GetResponse(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil) // nolint:gosec
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "None")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	return resp, nil
}
