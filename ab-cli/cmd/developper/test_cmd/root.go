package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/atomic-blend/backend/ab-cli/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests for an atomic blend backend service",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Prompt the user to select one component using bubbletea
			component, err := selectComponent(config.PlatformComponents)
			if err != nil {
				fmt.Printf("selection cancelled: %v\n", err)
				return nil
			}

			// Locate repository root and build absolute path to component.
			repoRoot, err := findRepoRoot()
			if err != nil {
				fmt.Printf("could not locate repo root: %v\n", err)
				return nil
			}
			componentDir := filepath.Join(repoRoot, component)

			// Ensure directory exists
			if _, err := os.Stat(componentDir); os.IsNotExist(err) {
				fmt.Printf("component directory not found: %s\n", componentDir)
				return nil
			}

			// Run linter: prefer `golint`, fallback to `golangci-lint run` if available.
			var lintCmd *exec.Cmd
			if _, err := exec.LookPath("golint"); err == nil {
				lintCmd = exec.Command("golint", "./...")
			} else if _, err := exec.LookPath("golangci-lint"); err == nil {
				lintCmd = exec.Command("golangci-lint", "run")
			} else {
				fmt.Printf("no linter (golint or golangci-lint) found in PATH — skipping lint step\n")
				lintCmd = nil
			}

			if lintCmd != nil {
				lintCmd.Dir = componentDir
				out, err := runWithSpinner(lintCmd, "linting "+component, true)
				// Treat any linter output as error (warnings are failures)
				if err != nil || len(out) > 0 {
					fmt.Printf("\n--- LINT FAILED for %s ---\n%s\n", component, strings.TrimSpace(string(out)))
					// Don't return an error to Cobra (avoid showing help); stop the command gracefully.
					return nil
				}
			}

			// Run tests
			testCmd := exec.Command("go", "test", "./...")
			testCmd.Dir = componentDir
			out, err := runWithSpinner(testCmd, "testing "+component, false)
			if err != nil {
				fmt.Printf("\n--- TESTS FAILED for %s ---\n%s\n", component, strings.TrimSpace(string(out)))
				// Don't return an error to Cobra (avoid showing help); stop the command gracefully.
				return nil
			}

			fmt.Printf("tests passed for %s\n", component)
			return nil
		},
	}
	return cmd
}

// findRepoRoot walks up from the current working directory until it finds
// a marker of the repository root (either `go.work` or `.git`). Returns an
// error if no root is found.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("repo root not found (go.work or .git)")
}

// selectorModel is a simple Bubbletea model used for selecting a string from a list.
type selectorModel struct {
	choices  []string
	cursor   int
	selected bool
}

func (m selectorModel) Init() tea.Cmd { return nil }

func (m selectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectorModel) View() string {
	var b strings.Builder
	b.WriteString("Use ↑/↓ or j/k to navigate, Enter to select, q to quit\n\n")
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	return b.String()
}

// selectComponent shows a simple Bubbletea list and returns the chosen item.
func selectComponent(choices []string) (string, error) {
	m := selectorModel{choices: choices}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if fm, ok := finalModel.(selectorModel); ok {
		if fm.selected || len(fm.choices) > 0 {
			return fm.choices[fm.cursor], nil
		}
	}
	return "", fmt.Errorf("no selection made")
}

// runWithSpinner runs the provided command while showing a simple terminal spinner.
// It captures combined stdout+stderr and returns it along with the command error.
func runWithSpinner(cmd *exec.Cmd, label string, treatOutputAsError bool) ([]byte, error) {
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	spinner := []rune{'|', '/', '-', '\\'}
	i := 0
	for {
		select {
		case err := <-done:
			// Determine success: if the process errored, it's a failure. If configured,
			// non-empty output also counts as failure (used for lint warnings).
			outputPresent := buf.Len() > 0
			failed := err != nil || (treatOutputAsError && outputPresent)
			if !failed {
				fmt.Printf("\r%s ✅\n", label)
			} else {
				fmt.Printf("\r%s ❌\n", label)
			}
			return buf.Bytes(), err
		default:
			fmt.Printf("\r%s %c", label, spinner[i%len(spinner)])
			time.Sleep(120 * time.Millisecond)
			i++
		}
	}
}
