// Package testall implements a small Bubble Tea TUI to run golint and tests
// for microservices. It is used by the `developper test` command.
package testall

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/atomic-blend/backend/ab-cli/config"
	btable "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Simple TUI that runs golint and go test for a list of services and shows
// a small table + live log. It's intentionally lightweight and does not
// replicate all the features of the external script; it focuses on providing
// a daemon-style visual while the checks run.

var servicesList = config.BackendServices

type svcStatus struct {
	Name     string
	Status   string // queued/running/completed/failed
	Golint   string // pending/running/passed/failed/unavailable
	Tests    string // pending/running/passed/failed/none
	GRPC     string // pending/running/passed/failed/none
	Progress string
}

type model struct {
	services []svcStatus
	logs     []string
	done     bool
	frame    int
	table    btable.Model
}

type updateSvcMsg struct {
	Index                                 int
	Status, Golint, Tests, GRPC, Progress string
}
type logMsg string
type finishedMsg struct{}
type tickMsg time.Time
type quitNow struct{}

// compositeMsg lets a background Cmd return several messages at once. The
// Update loop will unpack and apply them sequentially.
type compositeMsg struct{ msgs []tea.Msg }

func initialModel() model {
	svcs := make([]svcStatus, len(servicesList))
	for i, s := range servicesList {
		grpcState := "none"
		if s == "grpc" {
			grpcState = "pending"
		}
		svcs[i] = svcStatus{Name: s, Status: "queued", Golint: "pending", Tests: "pending", GRPC: grpcState, Progress: "Waiting..."}
	}
	// Build table
	cols := []btable.Column{
		{Title: "SERVICE", Width: 15},
		{Title: "STATUS", Width: 8},
		{Title: "GOLINT", Width: 8},
		{Title: "TESTS", Width: 8},
		{Title: "GRPC", Width: 8},
		{Title: "PROGRESS", Width: 20},
	}
	rows := servicesToRows(svcs, 0)
	// Set table height to fit rows + header to avoid the table filling the
	// remaining terminal height and creating large empty space below.
	// Also neutralize the selected style so the first row isn't highlighted
	// (this table is immutable / informational only).
	styles := btable.DefaultStyles()
	styles.Selected = lipgloss.NewStyle()

	tbl := btable.New(
		btable.WithColumns(cols),
		btable.WithRows(rows),
		btable.WithHeight(len(rows)+1),
		btable.WithFocused(false),
		btable.WithStyles(styles),
	)

	return model{services: svcs, logs: []string{}, done: false, table: tbl}
}

// servicesToRows converts the svcStatus slice into table rows.
func servicesToRows(svcs []svcStatus, frame int) []btable.Row {
	rows := make([]btable.Row, len(svcs))
	for i, s := range svcs {
		rows[i] = btable.Row{s.Name, getStatusIcon(s.Status), getAnimatedIcon(s.Golint, frame), getAnimatedIcon(s.Tests, frame), getAnimatedIcon(s.GRPC, frame), s.Progress}
	}
	return rows
}

func (m model) Init() tea.Cmd {
	// Start the background worker
	// start background work and tick for animation
	return tea.Batch(runAllChecks(), tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case compositeMsg:
		// Unpack and apply all inner messages sequentially
		for _, inner := range v.msgs {
			switch im := inner.(type) {
			case updateSvcMsg:
				if im.Index >= 0 && im.Index < len(m.services) {
					if im.Status != "" {
						m.services[im.Index].Status = im.Status
					}
					if im.Progress != "" {
						m.services[im.Index].Progress = im.Progress
					}
					if im.Golint != "" {
						m.services[im.Index].Golint = im.Golint
					}
					if im.Tests != "" {
						m.services[im.Index].Tests = im.Tests
					}
					if im.GRPC != "" {
						m.services[im.Index].GRPC = im.GRPC
					}
				}
			case logMsg:
				m.logs = append(m.logs, string(im))
				if len(m.logs) > 200 {
					m.logs = m.logs[len(m.logs)-200:]
				}
			case finishedMsg:
				m.done = true
				m.logs = append(m.logs, "✅ All checks finished.")
			}
		}
		// update table rows after applying all inner messages and set height
		rows := servicesToRows(m.services, m.frame)
		m.table.SetRows(rows)
		m.table.SetHeight(len(rows) + 1)
		return m, nil
	case updateSvcMsg:
		if v.Index >= 0 && v.Index < len(m.services) {
			if v.Status != "" {
				m.services[v.Index].Status = v.Status
			}
			if v.Progress != "" {
				m.services[v.Index].Progress = v.Progress
			}
			if v.Golint != "" {
				m.services[v.Index].Golint = v.Golint
			}
			if v.Tests != "" {
				m.services[v.Index].Tests = v.Tests
			}
			if v.GRPC != "" {
				m.services[v.Index].GRPC = v.GRPC
			}
		}
		// sync table with updated service state and adjust height
		rows := servicesToRows(m.services, m.frame)
		m.table.SetRows(rows)
		m.table.SetHeight(len(rows) + 1)
		return m, nil
	case logMsg:
		// keep logs bounded to last 200 lines
		m.logs = append(m.logs, string(v))
		if len(m.logs) > 200 {
			m.logs = m.logs[len(m.logs)-200:]
		}
		return m, nil
	case finishedMsg:
		// Finalize UI state: convert lingering running/pending states to clear final values
		for i := range m.services {
			// Convert running -> passed/completed where appropriate
			if m.services[i].Golint == "running" {
				m.services[i].Golint = "passed"
			}
			if m.services[i].Tests == "running" {
				m.services[i].Tests = "passed"
			}
			if m.services[i].GRPC == "running" {
				m.services[i].GRPC = "passed"
			}

			// Convert pending -> none so there is no spinner left
			if m.services[i].Golint == "pending" {
				m.services[i].Golint = "none"
			}
			if m.services[i].Tests == "pending" {
				m.services[i].Tests = "none"
			}
			if m.services[i].GRPC == "pending" {
				m.services[i].GRPC = "none"
			}

			// Normalize progress text
			if m.services[i].Progress == "Running tests..." || m.services[i].Progress == "Running golint..." || m.services[i].Progress == "Running gRPC lint..." {
				m.services[i].Progress = "Completed"
			}

			// Compute overall status: failed if any check failed, else completed
			if m.services[i].Golint == "failed" || m.services[i].Tests == "failed" || m.services[i].GRPC == "failed" {
				m.services[i].Status = "failed"
			} else {
				m.services[i].Status = "completed"
			}
		}
		m.done = true
		m.logs = append(m.logs, "✅ All checks finished.")
		// trigger an extra tick to refresh animation/frame one last time
		// and schedule an automatic quit after 3 seconds so the TUI closes
		quitCmd := func() tea.Msg {
			time.Sleep(3 * time.Second)
			return quitNow{}
		}
		// refresh table rows to show final icons and adjust height
		rows := servicesToRows(m.services, m.frame)
		m.table.SetRows(rows)
		m.table.SetHeight(len(rows) + 1)
		return m, tea.Batch(tickCmd(), quitCmd)
	case tickMsg:
		// advance animation frame and schedule next tick while not done
		m.frame = (m.frame + 1)
		// update animated icons in the table and keep height fixed to rows
		rows := servicesToRows(m.services, m.frame)
		m.table.SetRows(rows)
		m.table.SetHeight(len(rows) + 1)
		if !m.done {
			return m, tickCmd()
		}
		return m, nil
	case tea.KeyMsg:
		if v.String() == "q" || v.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case quitNow:
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	headerStyle := lipgloss.NewStyle().Bold(true)
	b.WriteString(headerStyle.Render("Atomic Blend — Tests & Linting Daemon"))
	b.WriteString("\n\n")

	// Table view (uses bubbles/table)
	b.WriteString(m.table.View())

	// Live logs display removed — keeping TUI compact
	// Keep exactly one blank line between table and footer
	b.WriteString("\n")

	if m.done {
		b.WriteString("Press q to quit.\n")
	} else {
		b.WriteString("Running... press q to quit (will not stop checks).\n")
	}

	return b.String()
}

// tickCmd returns a command that fires a tickMsg after 200ms
func tickCmd() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func getStatusIcon(status string) string {
	switch status {
	case "queued":
		return "⏳"
	case "running":
		return "🔄"
	case "completed":
		return "✅"
	case "failed":
		return "❌"
	default:
		return "❓"
	}
}

func getAnimatedIcon(status string, frame int) string {
	if status == "running" {
		spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		return spinners[frame%len(spinners)]
	}
	switch status {
	case "pending":
		return "⏳"
	case "passed":
		return "✅"
	case "failed":
		return "❌"
	case "none":
		return "—"
	case "unavailable":
		return "⚠️"
	default:
		return "?"
	}
}

// StartTUI starts the Bubble Tea application. The context can be used to
// cancel long-running commands in the future (not wired to keypresses here).
func StartTUI(ctx context.Context) error {
	p := tea.NewProgram(initialModel())
	return p.Start()
}

// runAllChecks returns a tea.Cmd that launches the background workflow and
// sends messages back into the Bubble Tea program.
func runAllChecks() tea.Cmd {
	// Build per-service commands that run the checks sequentially for that
	// service and use a WaitGroup so the final finishedMsg is emitted only
	// after all services are done. This avoids the previous behaviour where
	// the final finished message could be emitted immediately when using
	// `tea.Batch` on all cmds, leaving some long-running checks showing a
	// spinner.
	wd, _ := os.Getwd()
	var wg sync.WaitGroup
	cmds := make([]tea.Cmd, 0, len(servicesList)+1)

	for i, svc := range servicesList {
		idx := i
		svcName := svc
		svcPath := filepath.Join(wd, "..", svcName)
		if _, err := os.Stat(svcPath); os.IsNotExist(err) {
			svcPath = filepath.Join(wd, "..", "..", svcName)
		}

		wg.Add(1)

		// Create a command that runs golint, tests and grpc lint (for grpc)
		// sequentially for this service and returns a composite message with
		// the sequence of update/log messages for the service.
		cmds = append(cmds, func() tea.Msg {
			defer wg.Done()
			msgs := make([]tea.Msg, 0)

			// mark golint running
			msgs = append(msgs, updateSvcMsg{Index: idx, Status: "running", Golint: "running", Progress: "Running golint..."})

			// run golint
			passed, out := runGolint(svcPath)
			if passed {
				msgs = append(msgs, updateSvcMsg{Index: idx, Golint: "passed", Progress: "Running tests..."})
				msgs = append(msgs, logMsg(fmt.Sprintf("✔ [%s] golint passed", svcName)))
			} else {
				st := "failed"
				if strings.Contains(strings.ToLower(out), "not found") || strings.Contains(strings.ToLower(out), "executable file not found") || strings.Contains(strings.ToLower(out), "golint not installed") {
					st = "unavailable"
				}
				msgs = append(msgs, updateSvcMsg{Index: idx, Golint: st, Progress: "Running tests..."})
				msgs = append(msgs, logMsg(fmt.Sprintf("✖ [%s] golint: %s", svcName, summarizeOutput(out, 200))))
			}

			// Check for go.mod
			hasGoMod := false
			if _, err := os.Stat(filepath.Join(svcPath, "go.mod")); err == nil {
				hasGoMod = true
			}

			if !hasGoMod {
				msgs = append(msgs, updateSvcMsg{Index: idx, Tests: "none", Progress: "No tests", Status: "completed"})
				msgs = append(msgs, logMsg(fmt.Sprintf("→ [%s] no go.mod, skipping tests", svcName)))
				return compositeMsg{msgs: msgs}
			}

			// mark tests running
			msgs = append(msgs, updateSvcMsg{Index: idx, Tests: "running", Progress: "Running tests..."})

			passed, out = runTests(svcPath)
			if passed {
				msgs = append(msgs, updateSvcMsg{Index: idx, Tests: "passed", Progress: "Completed", Status: "completed"})
				msgs = append(msgs, logMsg(fmt.Sprintf("✔ [%s] tests passed", svcName)))
			} else {
				msgs = append(msgs, updateSvcMsg{Index: idx, Tests: "failed", Progress: "Completed", Status: "failed"})
				msgs = append(msgs, logMsg(fmt.Sprintf("✖ [%s] tests failed: %s", svcName, summarizeOutput(out, 200))))
			}

			// If this is the grpc service, run gRPC lint too
			if svcName == "grpc" {
				msgs = append(msgs, updateSvcMsg{Index: idx, GRPC: "running", Progress: "Running gRPC lint..."})
				passed, out := runGRPCLint(svcPath)
				if passed {
					msgs = append(msgs, updateSvcMsg{Index: idx, GRPC: "passed", Progress: "Completed"})
					msgs = append(msgs, logMsg(fmt.Sprintf("✔ [%s] gRPC lint passed", svcName)))
				} else {
					msgs = append(msgs, updateSvcMsg{Index: idx, GRPC: "failed", Progress: "Completed"})
					msgs = append(msgs, logMsg(fmt.Sprintf("✖ [%s] gRPC lint: %s", svcName, summarizeOutput(out, 200))))
				}
			}

			return compositeMsg{msgs: msgs}
		})
	}

	// final cmd: wait for all per-service cmds to finish, then send finishedMsg
	cmds = append(cmds, func() tea.Msg {
		wg.Wait()
		return finishedMsg{}
	})

	return tea.Batch(cmds...)
}

// runGRPCLint runs `buf lint` in the given service path if proto files exist.
func runGRPCLint(servicePath string) (bool, string) {
	protoFiles, err := filepath.Glob(filepath.Join(servicePath, "**/*.proto"))
	if err != nil || len(protoFiles) == 0 {
		return true, "No proto files found"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "buf", "lint")
	cmd.Dir = servicePath

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	combined := out.String() + "\n" + stderr.String()
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "executable file not found") || strings.Contains(strings.ToLower(combined), "not found") {
			return false, "buf not installed or not in PATH"
		}
		return false, combined
	}
	return true, combined
}

func summarizeOutput(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func runGolint(servicePath string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// If golint is not present, try to call `golangci-lint` as fallback
	// Prefer simple golint invocation for parity with original script.
	cmd := exec.CommandContext(ctx, "golint", "-set_exit_status", "./...")
	cmd.Dir = servicePath

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	combined := out.String() + "\n" + stderr.String()
	if err != nil {
		// If golint not installed, return failure with hint
		if strings.Contains(strings.ToLower(err.Error()), "executable file not found") || strings.Contains(strings.ToLower(combined), "not found") {
			return false, "golint not installed or not in PATH"
		}
		// golint returns non-zero when issues found; treat that as failed
		return false, combined
	}
	// no error => passed
	return true, combined
}

func runTests(servicePath string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-json")
	cmd.Dir = servicePath

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	combined := out.String() + "\n" + stderr.String()
	if err != nil {
		return false, combined
	}
	return true, combined
}
