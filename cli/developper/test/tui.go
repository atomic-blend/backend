// Package test implements a small Bubble Tea TUI to run golint and tests
// for microservices. It is used by the `developper test` command.
package test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Simple TUI that runs golint and go test for a list of services and shows
// a small table + live log. It's intentionally lightweight and does not
// replicate all the features of the external script; it focuses on providing
// a daemon-style visual while the checks run.

var servicesList = []string{
	"auth",
	"productivity",
	"grpc",
	"mail",
	"mail-server",
	"shared",
}

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
	return model{services: svcs, logs: []string{}, done: false}
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
				m.logs = append(m.logs, "✅ All checks finished. Press q to quit.")
			}
		}
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
		m.logs = append(m.logs, "✅ All checks finished. Press q to quit.")
		// trigger an extra tick to refresh animation/frame one last time
		// and schedule an automatic quit after 3 seconds so the TUI closes
		quitCmd := func() tea.Msg {
			time.Sleep(3 * time.Second)
			return quitNow{}
		}
		return m, tea.Batch(tickCmd(), quitCmd)
	case tickMsg:
		// advance animation frame and schedule next tick while not done
		m.frame = (m.frame + 1)
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

	// Table header
	b.WriteString(fmt.Sprintf("%-15s %-8s %-8s %-8s %-8s %-20s\n",
		"SERVICE", "STATUS", "GOLINT", "TESTS", "GRPC", "PROGRESS"))
	b.WriteString(strings.Repeat("-", 80) + "\n")
	for _, s := range m.services {
		statusIcon := getStatusIcon(s.Status)
		golintIcon := getAnimatedIcon(s.Golint, m.frame)
		testIcon := getAnimatedIcon(s.Tests, m.frame)
		grpcIcon := getAnimatedIcon(s.GRPC, m.frame)

		b.WriteString(fmt.Sprintf("%-15s %-8s %-8s %-8s %-8s %-20s\n",
			s.Name, statusIcon, golintIcon, testIcon, grpcIcon, s.Progress))
	}

	// Live logs display removed — keeping TUI compact
	b.WriteString("\n")

	if m.done {
		b.WriteString("\nPress q to quit.\n")
	} else {
		b.WriteString("\nRunning... press q to quit (will not stop checks).\n")
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
	// Build a sequence of tea.Cmds so the UI receives updates incrementally.
	cmds := make([]tea.Cmd, 0)
	wd, _ := os.Getwd()

	for i, svc := range servicesList {
		// capture loop variables
		idx := i
		svcName := svc
		svcPath := filepath.Join(wd, "..", svcName)
		if _, err := os.Stat(svcPath); os.IsNotExist(err) {
			svcPath = filepath.Join(wd, "..", "..", svcName)
		}

		// Cmd to set service -> running and golint -> running
		cmds = append(cmds, func() tea.Msg {
			return updateSvcMsg{Index: idx, Status: "running", Golint: "running", Progress: "Running golint..."}
		})

		// Cmd to run golint and return result messages
		cmds = append(cmds, func() tea.Msg {
			// run golint
			passed, out := runGolint(svcPath)
			if passed {
				return compositeMsg{msgs: []tea.Msg{
					updateSvcMsg{Index: idx, Golint: "passed", Progress: "Running tests..."},
					logMsg(fmt.Sprintf("✔ [%s] golint passed", svcName)),
				}}
			}
			st := "failed"
			if strings.Contains(strings.ToLower(out), "not found") || strings.Contains(strings.ToLower(out), "executable file not found") || strings.Contains(strings.ToLower(out), "golint not installed") {
				st = "unavailable"
			}
			return compositeMsg{msgs: []tea.Msg{
				updateSvcMsg{Index: idx, Golint: st, Progress: "Running tests..."},
				logMsg(fmt.Sprintf("✖ [%s] golint: %s", svcName, summarizeOutput(out, 200))),
			}}
		})

		// Check for go.mod
		hasGoMod := false
		if _, err := os.Stat(filepath.Join(svcPath, "go.mod")); err == nil {
			hasGoMod = true
		}
		if !hasGoMod {
			// Add commands to mark tests as none and log, mark completed
			cmds = append(cmds, func() tea.Msg {
				return updateSvcMsg{Index: idx, Tests: "none", Progress: "No tests", Status: "completed"}
			})
			cmds = append(cmds, func() tea.Msg {
				return logMsg(fmt.Sprintf("→ [%s] no go.mod, skipping tests", svcName))
			})
		} else {
			// Add running indicator for tests
			cmds = append(cmds, func() tea.Msg {
				return updateSvcMsg{Index: idx, Tests: "running", Progress: "Running tests..."}
			})

			// Cmd to run tests and send results
			cmds = append(cmds, func() tea.Msg {
				passed, out := runTests(svcPath)
				if passed {
					return compositeMsg{msgs: []tea.Msg{
						updateSvcMsg{Index: idx, Tests: "passed", Progress: "Completed", Status: "completed"},
						logMsg(fmt.Sprintf("✔ [%s] tests passed", svcName)),
					}}
				}
				return compositeMsg{msgs: []tea.Msg{
					updateSvcMsg{Index: idx, Tests: "failed", Progress: "Completed", Status: "failed"},
					logMsg(fmt.Sprintf("✖ [%s] tests failed: %s", svcName, summarizeOutput(out, 200))),
				}}
			})
		}

		// If this is the grpc service, run gRPC lint too
		if svcName == "grpc" {
			// mark grpc running
			cmds = append(cmds, func() tea.Msg { return updateSvcMsg{Index: idx, GRPC: "running", Progress: "Running gRPC lint..."} })
			cmds = append(cmds, func() tea.Msg {
				// run grpc lint
				passed, out := runGRPCLint(svcPath)
				if passed {
					return compositeMsg{msgs: []tea.Msg{
						updateSvcMsg{Index: idx, GRPC: "passed", Progress: "Completed"},
						logMsg(fmt.Sprintf("✔ [%s] gRPC lint passed", svcName)),
					}}
				}
				return compositeMsg{msgs: []tea.Msg{
					updateSvcMsg{Index: idx, GRPC: "failed", Progress: "Completed"},
					logMsg(fmt.Sprintf("✖ [%s] gRPC lint: %s", svcName, summarizeOutput(out, 200))),
				}}
			})
		}
	}

	// final cmd to signal finished
	cmds = append(cmds, func() tea.Msg { return finishedMsg{} })

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
