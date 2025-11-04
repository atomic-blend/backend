package bulkfiledownloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/atomic-blend/backend/cli/ui/types"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
)

// messages
type progressMsg struct {
	idx     int
	percent float64
}
type errMsg struct {
	idx int
	err error
}
type doneMsg struct{ idx int }

// model holds per-file spinner and state
type model struct {
	names     []string
	spinners  []spinner.Model
	percents  []float64
	errs      []error
	done      []bool
	quitting  bool
	total     int
	completed int
}

func (m model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.spinners))
	for i := range m.spinners {
		cmds = append(cmds, m.spinners[i].Tick)
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmds []tea.Cmd
		for i := range m.spinners {
			sp, cmd := m.spinners[i].Update(msg)
			m.spinners[i] = sp
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)
	case progressMsg:
		if msg.idx >= 0 && msg.idx < len(m.percents) {
			m.percents[msg.idx] = msg.percent
		}
		return m, nil
	case errMsg:
		if msg.idx >= 0 && msg.idx < len(m.errs) {
			m.errs[msg.idx] = msg.err
			if !m.done[msg.idx] {
				m.done[msg.idx] = true
				m.completed++
			}
		}
		if m.completed >= m.total {
			return m, tea.Quit
		}
		return m, nil
	case doneMsg:
		if msg.idx >= 0 && msg.idx < len(m.done) {
			if !m.done[msg.idx] {
				m.done[msg.idx] = true
				m.completed++
			}
			m.percents[msg.idx] = 1.0
		}
		if m.completed >= m.total {
			return m, tea.Quit
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {
	s := fmt.Sprintf("\nDownloading %d files...\n\n", m.total)
	for i := range m.names {
		name := m.names[i]
		if m.errs[i] != nil {
			s += fmt.Sprintf("%s %s ❌\n", " ", name)
			continue
		}
		if m.done[i] {
			s += fmt.Sprintf("%s %s ✅\n", " ", name)
			continue
		}

		// show spinner and percent if available
		spView := m.spinners[i].View()
		pct := ""
		if m.percents[i] > 0 {
			pct = fmt.Sprintf(" %.0f%%", m.percents[i]*100)
		}
		s += fmt.Sprintf("%s %s%s\n", spView, name, pct)
	}
	return s
}

// DownloadBulk downloads the provided list of DownloadableFile concurrently
// and displays a TUI showing one line per file with a spinner and percent.
// Each DownloadableFile.LocalPath should be the final destination path for the
// file (absolute or relative).
func DownloadBulk(files []types.DownloadableFile) error {
	if len(files) == 0 {
		return nil
	}

	// prepare model slices
	n := len(files)
	names := make([]string, n)
	spinners := make([]spinner.Model, n)
	percents := make([]float64, n)
	errs := make([]error, n)
	done := make([]bool, n)

	for i := range files {
		names[i] = path.Base(files[i].LocalPath)
		sp := spinner.New()
		spinners[i] = sp
	}

	m := model{
		names:    names,
		spinners: spinners,
		percents: percents,
		errs:     errs,
		done:     done,
		total:    n,
	}

	p := tea.NewProgram(m)

	var wg sync.WaitGroup
	wg.Add(n)
	for i, f := range files {
		go func(idx int, file types.DownloadableFile) {
			defer wg.Done()

			url := file.URL
			dest := file.LocalPath
			if err := os.MkdirAll(path.Dir(dest), 0o755); err != nil {
				p.Send(errMsg{idx: idx, err: fmt.Errorf("mkdir: %w", err)})
				return
			}

			req, err := http.NewRequest("GET", url, nil) // nolint:gosec
			if err != nil {
				p.Send(errMsg{idx: idx, err: err})
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				p.Send(errMsg{idx: idx, err: err})
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				p.Send(errMsg{idx: idx, err: fmt.Errorf("status %d", resp.StatusCode)})
				return
			}

			total := resp.ContentLength

			tmp := dest + ".download"
			out, err := os.Create(tmp)
			if err != nil {
				p.Send(errMsg{idx: idx, err: err})
				return
			}
			defer out.Close()

			var downloaded int64
			buf := make([]byte, 32*1024)
			for {
				n, rerr := resp.Body.Read(buf)
				if n > 0 {
					if _, werr := out.Write(buf[:n]); werr != nil {
						p.Send(errMsg{idx: idx, err: werr})
						return
					}
					downloaded += int64(n)
					if total > 0 {
						p.Send(progressMsg{idx: idx, percent: float64(downloaded) / float64(total)})
					}
				}
				if rerr != nil {
					if rerr == io.EOF {
						break
					}
					p.Send(errMsg{idx: idx, err: rerr})
					return
				}
			}

			if err := os.Rename(tmp, dest); err != nil {
				p.Send(errMsg{idx: idx, err: err})
				return
			}

			p.Send(progressMsg{idx: idx, percent: 1.0})
			p.Send(doneMsg{idx: idx})
		}(i, f)
	}

	// background waiter; model will quit automatically when all done
	go func() { wg.Wait() }()

	if _, err := p.Run(); err != nil {
		log.Error().Err(err).Msg("bulk downloader run failed")
		return err
	}

	// add a blank line after the TUI so following logs have space
	fmt.Println()

	return nil
}
