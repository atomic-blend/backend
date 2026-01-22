// Package platformcomponentupdater provides a small TUI to check and apply
// updates for platform components listed in the project's .env.
package platformcomponentupdater

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"unicode"

	"github.com/atomic-blend/backend/ab-cli/config"
	env_files_utils "github.com/atomic-blend/backend/ab-cli/utils/env_files_utils"
	ghutils "github.com/atomic-blend/backend/ab-cli/utils/gh_utils"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

// messages
type latestMsg struct {
	idx    int
	latest string
	image  string
}
type fetchErrMsg struct {
	idx int
	err error
}

// model for table-based updater
type model struct {
	tbl       table.Model
	quitting  bool
	total     int
	completed int
	errors    []string
	// selection phase data
	selecting    bool
	candidateIdx []int
	choices      []string
	selectCursor int
	selected     map[int]bool
	selectTbl    table.Model
	applied      int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If we're in selection phase, consume keys for the selector.
		if m.selecting {
			switch msg.String() {
			case "up", "k", "down", "j":
				// forward navigation keys to the select table so its internal
				// cursor (and its Selected style) moves automatically
				m.selectTbl, _ = m.selectTbl.Update(msg)
				// sync our selectCursor with the table cursor
				m.selectCursor = m.selectTbl.Cursor()
				// rebuild selection table rows to update checkbox column
				origRows := m.tbl.Rows()
				rows := buildSelectionRows(origRows, m.candidateIdx, m.selected)
				m.selectTbl.SetRows(rows)
			case " ":
				if m.selected == nil {
					m.selected = make(map[int]bool)
				}
				// if the table exposes a cursor, use that as the selected
				// index, otherwise use m.selectCursor
				idx := m.selectTbl.Cursor()
				// If idx == 0 it's the "All" row: toggle all candidate rows
				totalRows := len(m.candidateIdx) + 1
				if idx == 0 {
					// toggle all: determine new state
					newState := !m.selected[0]
					m.selected[0] = newState
					for i := 1; i < totalRows; i++ {
						m.selected[i] = newState
					}
				} else {
					// toggle individual
					m.selected[idx] = !m.selected[idx]
					// update All checkbox state: set true only if all are selected
					allSel := true
					for i := 1; i < totalRows; i++ {
						if !m.selected[i] {
							allSel = false
							break
						}
					}
					m.selected[0] = allSel
				}
				// rebuild selection rows after toggling
				origRows := m.tbl.Rows()
				rows := buildSelectionRows(origRows, m.candidateIdx, m.selected)
				m.selectTbl.SetRows(rows)
			case "enter":
				// apply selected
				envPath := path.Join(config.CliConfig.Directory, ".env")
				for i := range m.choices {
					// m.selected is indexed by select table rows; candidates are
					// shifted by +1 because the first row is the "All" option.
					if !m.selected[i+1] {
						continue
					}
					rows := m.tbl.Rows()
					rowsIdx := m.candidateIdx[i]
					r := rows[rowsIdx]
					name := fmt.Sprintf("%v", r[0])
					latest := fmt.Sprintf("%v", r[2])
					envVar, ok := config.EnvComponentVersionMapping[name]
					if !ok {
						log.Error().Str("component", name).Msg("no env var mapping for component, skipping")
						fmt.Printf("Skipping %s: no env var mapping\n", name)
						continue
					}
					if err := env_files_utils.SetEnvVarValue(envPath, envVar, latest); err != nil {
						log.Error().Err(err).Str("component", name).Msg("failed to write env var")
						fmt.Printf("Failed to update %s in %s: %v\n", envVar, envPath, err)
						m.errors = append(m.errors, fmt.Sprintf("apply %s: %v", name, err))
						continue
					}
					fmt.Printf("Updated %s -> %s in %s\n", envVar, latest, envPath)
					m.applied++
				}
				m.quitting = true
				return m, tea.Quit
			case "q", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}
		// Not selecting: disable interactive navigation while checks are
		// running. Only accept quit keys, and Enter once checks complete
		// to enter selection mode. All other key events are ignored so
		// the table remains read-only during background fetches.
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			// only allow entering selection phase when checks complete
			if m.completed >= m.total {
				// build candidates from rows
				rows := m.tbl.Rows()
				candidateIdx, choices := getCandidatesFromRows(rows)
				if len(choices) == 0 {
					if len(m.errors) > 0 {
						fmt.Println("Errors encountered while fetching latest versions:")
						for _, e := range m.errors {
							fmt.Printf(" - %s\n", e)
						}
					}
					fmt.Println("No components to update in .env")
					m.quitting = true
					return m, tea.Quit
				}

				selectTbl := createSelectionTable(rows, candidateIdx)

				m.selecting = true
				m.candidateIdx = candidateIdx
				m.choices = choices
				m.selectCursor = 0
				m.selected = make(map[int]bool)
				// ensure All checkbox exists at key 0
				m.selected[0] = false
				m.selectTbl = selectTbl
				return m, nil
			}
			// ignore Enter while checks are incomplete
			return m, nil
		default:
			// ignore all other keys while not selecting
			return m, nil
		}
	case latestMsg:
		if msg.idx >= 0 && msg.idx < len(m.tbl.Rows()) {
			rows := m.tbl.Rows()
			rows = updateTableRowWithLatest(rows, msg.idx, msg.latest, msg.image)
			m.tbl.SetRows(rows)
			m.completed++
		}
		return m, nil
	case fetchErrMsg:
		if msg.idx >= 0 && msg.idx < len(m.tbl.Rows()) {
			rows := m.tbl.Rows()
			r := rows[msg.idx]
			rows[msg.idx] = table.Row{r[0], r[1], "-", statusError}
			m.tbl.SetRows(rows)
			// store full error message so it can be displayed below the table
			m.errors = append(m.errors, fmt.Sprintf("%s: %v", r[0], msg.err))
			m.completed++
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.tbl, cmd = m.tbl.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.selecting {
		v := m.selectTbl.View()
		// Trim trailing whitespace from the table view so appended
		// instructions don't create large blank areas when terminal
		// height is limited. Some terminals/tables pad with spaces
		// rather than pure newlines, so trim all trailing whitespace.
		v = strings.TrimRightFunc(v, unicode.IsSpace)
		v = v + "\nApply updates to .env (use ↑/↓, space to toggle, Enter to confirm, q to cancel)."
		if len(m.errors) > 0 {
			v = v + "\n\nErrors:\n"
			for _, e := range m.errors {
				v = v + " - " + e + "\n"
			}
		}
		return v
	}

	v := m.tbl.View()
	// Trim trailing whitespace from the table view for compact rendering
	// on terminals with small height.
	v = strings.TrimRightFunc(v, unicode.IsSpace)
	if m.completed >= m.total {
		v = v + "\nAll checks complete. Press Enter to continue or q to cancel."
	}
	if len(m.errors) > 0 {
		v = v + "\n\nErrors:\n"
		for _, e := range m.errors {
			v = v + " - " + e + "\n"
		}
	}
	return v
}

// GetUpdates displays an alt-screen table and fetches the latest tag for each
// provided platform component. It uses a background goroutine per component
// (same concurrency/daemon pattern as the bulk downloader) and updates the
// table as results come in.
func GetUpdates(components []env_files_utils.PlatformComponent, rc *bool) error {
	if len(components) == 0 {
		return nil
	}

	cols := []table.Column{
		{Title: "Name", Width: 24},
		{Title: "Current", Width: 14},
		{Title: "Latest", Width: 14},
		{Title: "Status", Width: 20},
	}

	rows := buildInitialRows(components)

	tbl := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(false),
	)
	// Start from the default styles but neutralize the Selected style so
	// no row appears highlighted when the table is not in selection mode.
	styles := table.DefaultStyles()
	styles.Selected = lipgloss.NewStyle()
	tbl.SetStyles(styles)

	m := model{tbl: tbl, total: len(components)}
	p := tea.NewProgram(m, tea.WithAltScreen())

	for i, c := range components {
		go func(idx int, comp env_files_utils.PlatformComponent) {
			latest, err := fetchLatestTagForImage(comp.Image, rc)
			if err != nil {
				p.Send(fetchErrMsg{idx: idx, err: err})
				return
			}
			p.Send(latestMsg{idx: idx, latest: latest, image: comp.Image})
		}(i, c)
	}

	finalModel, err := p.Run()
	if err != nil {
		log.Error().Err(err).Msg("updater run failed")
		return err
	}

	// extra newline after TUI
	fmt.Println()

	// cast to our model to inspect rows
	fm, ok := finalModel.(model)
	if !ok {
		log.Error().Msg("unexpected model type after run")
		return nil
	}

	rows = fm.tbl.Rows()
	// recompute candidate list to decide on messages after run
	_, choices := getCandidatesFromRows(rows)

	if len(choices) == 0 {
		if len(fm.errors) > 0 {
			fmt.Println("Errors encountered while fetching latest versions:")
			for _, e := range fm.errors {
				fmt.Printf(" - %s\n", e)
			}
		}
		fmt.Println("No components to update in .env")
		return nil
	}

	// If there were candidates but nothing applied, notify the user.
	if fm.applied == 0 {
		fmt.Println("No updates applied")
		return nil
	}

	return nil
}

// fetchLatestTagForImage tries to obtain the most recently updated tag for a
// docker hub image. If the image is empty or hosted on an unsupported
// registry, returns an error.
func fetchLatestTagForImage(image string, rc *bool) (string, error) {
	image = strings.TrimSpace(image)
	if image == "" {
		return "-", fmt.Errorf("no image specified")
	}

	// Only support ghcr.io for now
	if !strings.HasPrefix(image, "ghcr.io/") {
		return "-", fmt.Errorf("unsupported registry for image %s", image)
	}

	latest, err := ghutils.GetLatestImageVersion(image, rc)
	if err != nil {
		return "-", err
	}
	return latest, nil
}

// determineStatus returns one of: "OK", "OUT-OF-DATE", or "ERROR" based
// on a comparison between current and latest version strings.
func determineStatus(current, latest, image string) string {
	if latest == "-" {
		return "ERROR"
	}
	if current == latest {
		return "OK"
	}
	if cmp, ok := compareVersions(current, latest); ok {
		// compareVersions returns -1 if current < latest
		if cmp == -1 {
			return "OUT-OF-DATE"
		}
		// If numeric parts are equal but tags differ (e.g. RCs with different
		// commit hashes), use the tag creation timestamps from GHCR to decide.
		if cmp == 0 {
			if current != latest && strings.Contains(current, "-rc-") && strings.Contains(latest, "-rc-") {
				tcur, errCur := ghutils.GetTagCreatedAt(image, current)
				tlat, errLat := ghutils.GetTagCreatedAt(image, latest)
				if errCur != nil {
					log.Error().Err(errCur).Str("image", image).Str("tag", current).Msg("failed to get created_at for current tag")
				}
				if errLat != nil {
					log.Error().Err(errLat).Str("image", image).Str("tag", latest).Msg("failed to get created_at for latest tag")
				}
				if errCur != nil || errLat != nil {
					// If we cannot determine timestamps, treat as out-of-date
					// so the user is prompted to update rather than silently
					// accepting a potentially newer RC.
					return "OUT-OF-DATE"
				}
				if tcur.Before(tlat) {
					return "OUT-OF-DATE"
				}
				return "OK"
			}
		}
		return "OK"
	}
	// fallback: if strings differ, treat as out-of-date
	if current != latest {
		return "OUT-OF-DATE"
	}
	return "OK"
}

const (
	statusPending  = "pending"
	statusError    = "ERROR"
	statusOK       = "OK"
	statusOutdated = "OUT-OF-DATE"
)

// buildInitialRows constructs the default table rows for the provided
// platform components.
func buildInitialRows(components []env_files_utils.PlatformComponent) []table.Row {
	rows := make([]table.Row, len(components))
	for i, c := range components {
		rows[i] = table.Row{c.Name, c.Version, "-", statusPending}
	}
	return rows
}

// getCandidatesFromRows scans table rows and returns indices and user-visible
// choice strings for rows that are out-of-date.
func getCandidatesFromRows(rows []table.Row) ([]int, []string) {
	var candidateIdx []int
	var choices []string
	for i, r := range rows {
		if len(r) >= 4 {
			name := fmt.Sprintf("%v", r[0])
			current := fmt.Sprintf("%v", r[1])
			latest := fmt.Sprintf("%v", r[2])
			status := fmt.Sprintf("%v", r[3])
			if latest != "-" && latest != current && status == statusOutdated {
				candidateIdx = append(candidateIdx, i)
				choices = append(choices, fmt.Sprintf("%s: %s → %s", name, current, latest))
			}
		}
	}
	return candidateIdx, choices
}

// buildSelectionRows builds the rows for the selection table from the
// original rows and candidate indices. The `selected` map may be nil to
// indicate nothing is selected yet.
func buildSelectionRows(origRows []table.Row, candidateIdx []int, selected map[int]bool) []table.Row {
	rows := make([]table.Row, len(candidateIdx)+1)
	allChecked := false
	if selected != nil && selected[0] {
		allChecked = true
	}
	allCheck := "  [ ]"
	if allChecked {
		allCheck = "  [x]"
	}
	rows[0] = table.Row{allCheck, "All", "", "", ""}
	for i, origIdx := range candidateIdx {
		r := origRows[origIdx]
		check := "  [ ]"
		if selected != nil && selected[i+1] {
			check = "  [x]"
		}
		rows[i+1] = table.Row{check, r[0], r[1], r[2], r[3]}
	}
	return rows
}

// createSelectionTable builds a focused table.Model for selection using the
// original rows and candidate indices.
func createSelectionTable(origRows []table.Row, candidateIdx []int) table.Model {
	cols := []table.Column{
		{Title: "", Width: 8},
		{Title: "Name", Width: 24},
		{Title: "Current", Width: 14},
		{Title: "Latest", Width: 14},
		{Title: "Status", Width: 20},
	}
	rows2 := buildSelectionRows(origRows, candidateIdx, nil)
	selectTbl := table.New(
		table.WithColumns(cols),
		table.WithRows(rows2),
		table.WithFocused(true),
	)
	selectTbl.SetStyles(table.DefaultStyles())
	return selectTbl
}

// updateTableRowWithLatest updates the provided rows slice at index idx with
// the supplied latest value and computes the status.
func updateTableRowWithLatest(rows []table.Row, idx int, latest string, image string) []table.Row {
	if idx >= 0 && idx < len(rows) {
		r := rows[idx]
		current := fmt.Sprintf("%v", r[1])
		status := determineStatus(current, latest, image)
		rows[idx] = table.Row{r[0], r[1], latest, status}
	}
	return rows
}

// compareVersions attempts to compare two version strings by parsing
// numeric dot-separated segments. Returns -1 if a<b, 0 if equal, 1 if a>b.
// The bool indicates whether a numeric comparison was possible.
func compareVersions(a, b string) (int, bool) {
	ap, aok := parseNumericParts(a)
	bp, bok := parseNumericParts(b)
	if !aok || !bok {
		return 0, false
	}
	// compare segment-wise
	max := len(ap)
	if len(bp) > max {
		max = len(bp)
	}
	for i := 0; i < max; i++ {
		var av, bv int
		if i < len(ap) {
			av = ap[i]
		}
		if i < len(bp) {
			bv = bp[i]
		}
		if av < bv {
			return -1, true
		}
		if av > bv {
			return 1, true
		}
	}
	return 0, true
}

// parseNumericParts extracts leading numeric parts from a version string.
// It strips a leading 'v' or 'V', splits on '.', and for each segment
// takes the leading digit sequence (if any). Returns the slice and a
// boolean indicating whether at least one numeric part was found.
func parseNumericParts(s string) ([]int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, false
	}
	if s[0] == 'v' || s[0] == 'V' {
		s = s[1:]
	}
	parts := strings.Split(s, ".")
	nums := make([]int, 0, len(parts))
	any := false
	for _, p := range parts {
		digits := ""
		for _, r := range p {
			if r >= '0' && r <= '9' {
				digits += string(r)
			} else {
				break
			}
		}
		if digits == "" {
			nums = append(nums, 0)
			continue
		}
		v, err := strconv.Atoi(digits)
		if err != nil {
			nums = append(nums, 0)
			continue
		}
		nums = append(nums, v)
		any = true
	}
	return nums, any
}
