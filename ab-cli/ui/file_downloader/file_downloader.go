package filedownloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

// Download is the single public entrypoint: it downloads the file at url
// to the current working directory (using the URL basename) and displays
// a TUI progress bar while downloading. It returns an error on failure.
func Download(url string) error {
	resp, err := getResponse(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck

	// Need content length to show progress
	if resp.ContentLength <= 0 {
		return fmt.Errorf("can't parse content length, aborting download")
	}

	// Determine filename
	filename := path.Base(resp.Request.URL.Path)
	if filename == "" || filename == "/" || strings.HasSuffix(filename, "/") {
		filename = "download"
	}

	// Create destination file in current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not determine working directory: %w", err)
	}
	dstPath := path.Join(cwd, filename)
	file, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close() // nolint:errcheck

	pw := &progressWriter{
		Total:  int(resp.ContentLength),
		File:   file,
		Reader: resp.Body,
	}

	m := model{
		pw:       pw,
		Progress: progress.New(progress.WithDefaultGradient()),
	}

	p := tea.NewProgram(m)

	// wire the program into the progress writer so it can send messages
	pw.program = p

	// Start the download in background
	go pw.start()

	// Run the TUI; returns when user quits or download completes
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	log.Printf("downloaded to %s", dstPath)
	return nil
}

type progressWriter struct {
	Total      int
	Downloaded int
	File       *os.File
	Reader     io.Reader
	program    *tea.Program
}

func (pw *progressWriter) start() {
	// TeeReader calls pw.Write() each time a new response chunk is received
	_, err := io.Copy(pw.File, io.TeeReader(pw.Reader, pw))
	if err != nil {
		if pw.program != nil {
			pw.program.Send(progressErrMsg{err})
		}
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.Downloaded += len(p)
	if pw.Total > 0 && pw.program != nil {
		// send the ratio as a progress message
		pw.program.Send(progressMsg(float64(pw.Downloaded) / float64(pw.Total)))
	}
	return len(p), nil
}

type progressMsg float64

type progressErrMsg struct{ err error }

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg { return nil })
}

type model struct {
	pw       *progressWriter
	Progress progress.Model
	Err      error
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.Progress.Width = msg.Width - padding*2 - 4
		if m.Progress.Width > maxWidth {
			m.Progress.Width = maxWidth
		}
		return m, nil
	case progressErrMsg:
		m.Err = msg.err
		return m, tea.Quit
	case progressMsg:
		var cmds []tea.Cmd
		if msg >= 1.0 {
			cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
		}
		cmds = append(cmds, m.Progress.SetPercent(float64(msg)))
		return m, tea.Batch(cmds...)
	case progress.FrameMsg:
		progressModel, cmd := m.Progress.Update(msg)
		m.Progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.Err != nil {
		return "Error downloading: " + m.Err.Error() + "\n"
	}
	pad := strings.Repeat(" ", padding)
	return "\n" + pad + m.Progress.View() + "\n\n" + pad + helpStyle("Press any key to quit")
}

func getResponse(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil) // nolint:gosec
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "None")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	return resp, nil
}
