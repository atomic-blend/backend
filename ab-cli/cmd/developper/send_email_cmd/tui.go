package sendemailcmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	sendemail "github.com/atomic-blend/backend/ab-cli/internal/sendemail"
)

type model struct {
	focusIndex int
	from       textinput.Model
	to         textinput.Model
	cc         textinput.Model
	bcc        textinput.Model
	subject    textinput.Model
	smtp       textinput.Model
	thread     textinput.Model
	body       textarea.Model
	status     string
	done       bool
	cfg        *sendemail.EmailConfig
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func initialModel() model {
	// create separate textinput instances for each field
	from := textinput.New()
	from.Placeholder = "you@domain.com"
	from.CharLimit = 256
	from.Width = 36

	to := textinput.New()
	to.Placeholder = "recipient@domain.com"
	to.CharLimit = 256
	to.Width = 36

	// default recipient
	to.SetValue("brandon@brandonguigo.com")

	cc := textinput.New()
	cc.Placeholder = "cc@domain.com"
	cc.CharLimit = 256
	cc.Width = 36

	bcc := textinput.New()
	bcc.Placeholder = "bcc@domain.com"
	bcc.CharLimit = 256
	bcc.Width = 36

	subj := textinput.New()
	subj.Placeholder = "Email subject"
	subj.Width = 50

	smtp := textinput.New()
	// do not show the default SMTP placeholder in the UI; default is applied on send
	smtp.Placeholder = ""
	smtp.Width = 36

	thread := textinput.New()
	thread.Placeholder = "1"
	thread.CharLimit = 6
	thread.Width = 6
	thread.SetValue("1")

	ta := textarea.New()
	// adjust for VSCode integrated terminal if detected
	if os.Getenv("TERM_PROGRAM") == "vscode" {
		ta.SetWidth(60)
		ta.SetHeight(8)
	} else {
		ta.SetWidth(80)
		ta.SetHeight(10)
	}
	ta.Placeholder = "Write your message here... (press Ctrl+S to send)"

	m := model{
		focusIndex: 0,
		from:       from,
		to:         to,
		cc:         cc,
		bcc:        bcc,
		subject:    subj,
		smtp:       smtp,
		thread:     thread,
		body:       ta,
		status:     "",
		done:       false,
	}

	// generate default values (from, subject, body) using the project's faker helper
	if cfg, err := sendemail.GenerateRandomConfig(sendemail.RandomOptions{}); err == nil && cfg != nil {
		m.from.SetValue(cfg.Sender)
		m.subject.SetValue(cfg.Subject)
		// textarea.SetValue exists on the textarea.Model API
		m.body.SetValue(cfg.Body)
	}

	// give initial focus to the first input
	m.from.Focus()
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		switch k {
		case "ctrl+c", "esc":
			m.done = true
			return m, tea.Quit
		case "ctrl+s":
			// send immediately
			m.status = "Sending..."
			cfg := &sendemail.EmailConfig{
				Sender:        strings.TrimSpace(m.from.Value()),
				Recipients:    parseComma(m.to.Value()),
				CCRecipients:  parseComma(m.cc.Value()),
				BCCRecipients: parseComma(m.bcc.Value()),
				Subject:       strings.TrimSpace(m.subject.Value()),
				Body:          m.body.Value(),
				SMTPServer:    strings.TrimSpace(m.smtp.Value()),
			}
			// parse thread size
			if n, err := strconv.Atoi(strings.TrimSpace(m.thread.Value())); err == nil && n > 0 {
				cfg.ThreadSize = n
			} else {
				cfg.ThreadSize = 1
			}
			m.cfg = cfg
			m.done = true
			return m, tea.Quit
		case "tab":
			m.focusIndex = (m.focusIndex + 1) % 8
		case "shift+tab":
			m.focusIndex = (m.focusIndex + 7) % 8
		case "enter":
			// When body (textarea) has focus, let it handle Enter (create newline).
			if m.focusIndex != 7 {
				m.focusIndex = (m.focusIndex + 1) % 8
			}
		case "shift+enter", "ctrl+enter", "alt+enter", "meta+enter":
			// Terminals frequently don't differentiate Shift+Enter; accept common variants.
			// Send from anywhere in the composer (not only when body has focus).
			m.status = "Sending..."
			cfg := &sendemail.EmailConfig{
				Sender:        strings.TrimSpace(m.from.Value()),
				Recipients:    parseComma(m.to.Value()),
				CCRecipients:  parseComma(m.cc.Value()),
				BCCRecipients: parseComma(m.bcc.Value()),
				Subject:       strings.TrimSpace(m.subject.Value()),
				Body:          m.body.Value(),
				SMTPServer:    strings.TrimSpace(m.smtp.Value()),
			}
			if n, err := strconv.Atoi(strings.TrimSpace(m.thread.Value())); err == nil && n > 0 {
				cfg.ThreadSize = n
			} else {
				cfg.ThreadSize = 1
			}
			m.cfg = cfg
			m.done = true
			return m, tea.Quit
		case "up":
			// if body has focus, let it handle arrow keys
			if m.focusIndex != 7 {
				m.focusIndex = (m.focusIndex + 6) % 8
			}
		case "down":
			if m.focusIndex != 7 {
				m.focusIndex = (m.focusIndex + 1) % 8
			}
		case "left":
			if m.focusIndex != 7 {
				m.focusIndex = (m.focusIndex + 6) % 8
			}
		case "right":
			if m.focusIndex != 7 {
				m.focusIndex = (m.focusIndex + 1) % 8
			}
		}
	}

	// Ensure only the focused component has focus state
	switch m.focusIndex {
	case 0:
		m.from.Focus()
		m.to.Blur()
		m.cc.Blur()
		m.bcc.Blur()
		m.subject.Blur()
		m.smtp.Blur()
		m.thread.Blur()
		m.body.Blur()
	case 1:
		m.to.Focus()
		m.from.Blur()
		m.cc.Blur()
		m.bcc.Blur()
		m.subject.Blur()
		m.smtp.Blur()
		m.thread.Blur()
		m.body.Blur()
	case 2:
		m.cc.Focus()
		m.from.Blur()
		m.to.Blur()
		m.bcc.Blur()
		m.subject.Blur()
		m.smtp.Blur()
		m.thread.Blur()
		m.body.Blur()
	case 3:
		m.bcc.Focus()
		m.from.Blur()
		m.to.Blur()
		m.cc.Blur()
		m.subject.Blur()
		m.smtp.Blur()
		m.thread.Blur()
		m.body.Blur()
	case 4:
		m.subject.Focus()
		m.from.Blur()
		m.to.Blur()
		m.cc.Blur()
		m.bcc.Blur()
		m.smtp.Blur()
		m.thread.Blur()
		m.body.Blur()
	case 5:
		m.smtp.Focus()
		m.from.Blur()
		m.to.Blur()
		m.cc.Blur()
		m.bcc.Blur()
		m.subject.Blur()
		m.thread.Blur()
		m.body.Blur()
	case 6:
		m.thread.Focus()
		m.from.Blur()
		m.to.Blur()
		m.cc.Blur()
		m.bcc.Blur()
		m.subject.Blur()
		m.smtp.Blur()
		m.body.Blur()
	case 7:
		m.body.Focus()
		m.from.Blur()
		m.to.Blur()
		m.cc.Blur()
		m.bcc.Blur()
		m.subject.Blur()
		m.smtp.Blur()
		m.thread.Blur()
	}

	// Update focused component
	switch m.focusIndex {
	case 0:
		m.from, cmd = m.from.Update(msg)
	case 1:
		m.to, cmd = m.to.Update(msg)
	case 2:
		m.cc, cmd = m.cc.Update(msg)
	case 3:
		m.bcc, cmd = m.bcc.Update(msg)
	case 4:
		m.subject, cmd = m.subject.Update(msg)
	case 5:
		m.smtp, cmd = m.smtp.Update(msg)
	case 6:
		m.thread, cmd = m.thread.Update(msg)
	case 7:
		m.body, cmd = m.body.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.done && m.cfg != nil {
		return ""
	}
	// helper to render underscores when empty and not focused
	renderInput := func(t textinput.Model, width int) string {
		if strings.TrimSpace(t.Value()) == "" && !t.Focused() {
			return strings.Repeat("_", width)
		}
		return t.View()
	}
	renderTextarea := func(t textarea.Model, width int) string {
		if strings.TrimSpace(t.Value()) == "" && !t.Focused() {
			return strings.Repeat("_", width)
		}
		v := t.View()
		// remove a single leading tab from the first line so it aligns with other lines
		lines := strings.Split(v, "\n")
		if len(lines) > 0 {
			lines[0] = strings.TrimPrefix(lines[0], "\t")
		}
		return strings.Join(lines, "\n")
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render("Compose Email") + "\n\n")
	b.WriteString(labelStyle.Render("From: ") + renderInput(m.from, 36) + "\n")
	b.WriteString(labelStyle.Render("To:   ") + renderInput(m.to, 36) + "\n")
	b.WriteString(labelStyle.Render("CC:   ") + renderInput(m.cc, 36) + "\n")
	b.WriteString(labelStyle.Render("BCC:  ") + renderInput(m.bcc, 36) + "\n")
	b.WriteString(labelStyle.Render("Subject:") + " " + renderInput(m.subject, 50) + "\n\n")
	// show default SMTP when empty and not focused
	smtpView := renderInput(m.smtp, 36)
	if strings.TrimSpace(m.smtp.Value()) == "" && !m.smtp.Focused() {
		smtpView = "localhost:1025"
	}
	b.WriteString(labelStyle.Render("SMTP: ") + smtpView + "\n\n")
	b.WriteString(labelStyle.Render("Thread: ") + renderInput(m.thread, 6) + "\n\n")
	b.WriteString(labelStyle.Render("Body:\n"))
	b.WriteString(renderTextarea(m.body, func() int {
		// match the textarea width used in initialModel
		if os.Getenv("TERM_PROGRAM") == "vscode" {
			return 60
		}
		return 80
	}()) + "\n")
	b.WriteString("\n" + lipgloss.NewStyle().Italic(true).Render("Tab to navigate, Enter to send when focused on Body, Esc to cancel"))
	if m.status != "" {
		b.WriteString("\n\n" + m.status)
	}
	return b.String()
}

func parseComma(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// RunComposer runs the Bubble Tea composer and returns the filled config.
func RunComposer() (*sendemail.EmailConfig, error) {
	m := initialModel()
	// Prefer alt-screen unless running inside VSCode integrated terminal
	var p *tea.Program
	if os.Getenv("TERM_PROGRAM") == "vscode" {
		p = tea.NewProgram(m)
	} else {
		p = tea.NewProgram(m, tea.WithAltScreen())
	}
	final, err := p.Run()
	if err != nil {
		return nil, err
	}
	fm := final.(model)
	if fm.cfg == nil {
		return nil, fmt.Errorf("composer cancelled")
	}
	// defaults
	if fm.cfg.SMTPServer == "" {
		fm.cfg.SMTPServer = "localhost:1025"
	}
	return fm.cfg, nil
}
