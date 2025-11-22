// Package sendemailcmd provides the `send-email` developper subcommand.
package sendemailcmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	sendemail "github.com/atomic-blend/backend/ab-cli/internal/sendemail"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var random = true
	var threadSize int
	var delayMs int
	var classic bool

	cmd := &cobra.Command{
		Use:   "send-email",
		Short: "Send a test email (interactive or random)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// show summary and confirm before sending using a Bubble Tea table

			// Default behavior: run Bubble Tea composer UI unless --classic is set
			if !classic {
				cfg, err := RunComposer()
				if err != nil {
					return err
				}
				// after composer, determine effective thread size and pre-generate IDs
				effectiveThread := threadSize
				if effectiveThread <= 0 {
					effectiveThread = cfg.ThreadSize
				}
				if effectiveThread <= 0 {
					effectiveThread = 1
				}

				ids := make([]string, effectiveThread)
				for i := 0; i < effectiveThread; i++ {
					ids[i] = genMessageID()
				}

				if cfg.AttachmentPath != "" {
					fmt.Printf("📎 Attachment: %s\n", cfg.AttachmentPath)
				}
				if effectiveThread > 1 {
					fmt.Printf("🔁 Sending thread of %d messages\n", effectiveThread)
				}
				recipients := strings.Join(cfg.Recipients, ",")
				if len(cfg.CCRecipients) > 0 {
					recipients += ", cc:" + strings.Join(cfg.CCRecipients, ",")
				}
				summary := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nBody length: %d chars\nAttachment: %v\nSMTP: %s\nThread size: %d\n", cfg.Sender, recipients, cfg.Subject, len(cfg.Body), cfg.AttachmentPath != "", cfg.SMTPServer, effectiveThread)
				// append full headers for each email in the thread so user sees exact headers
				summary += buildThreadHeaders(cfg, ids)

				ok, err := showSummaryTUI(summary)
				if err != nil {
					return err
				}
				if !ok {
					fmt.Println("❌ Cancelled by user")
					return nil
				}
				// send thread using the pre-generated IDs
				var prevID string
				for i := 0; i < effectiveThread; i++ {
					if i > 0 {
						cfg.InReplyTo = prevID
						cfg.References = append(cfg.References, prevID)
					}
					cfg.MessageID = ids[i]
					msg, err := sendemail.CreateRFCMessage(cfg)
					if err != nil {
						return err
					}
					if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
						return err
					}
					fmt.Printf("📨 Sent message id %s to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
					prevID = cfg.MessageID
					if delayMs > 0 {
						time.Sleep(time.Duration(delayMs) * time.Millisecond)
					}
				}
				return nil
			}

			if random {
				opts := sendemail.RandomOptions{}
				cfg, err := sendemail.GenerateRandomConfig(opts)
				if err != nil {
					return err
				}
				fmt.Printf("🔀 Using random content: From=%s To=%s Subject=%s\n", cfg.Sender, strings.Join(cfg.Recipients, ","), cfg.Subject)
				if cfg.AttachmentPath != "" {
					fmt.Printf("📎 Attachment: %s\n", cfg.AttachmentPath)
				}
				// determine effective thread size and pre-generate IDs
				effectiveThread := threadSize
				if effectiveThread <= 0 {
					effectiveThread = 1
				}
				ids := make([]string, effectiveThread)
				for i := 0; i < effectiveThread; i++ {
					ids[i] = genMessageID()
				}

				// build and show summary
				recipients := strings.Join(cfg.Recipients, ",")
				if len(cfg.CCRecipients) > 0 {
					recipients += ", cc:" + strings.Join(cfg.CCRecipients, ",")
				}
				summary := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nBody length: %d chars\nAttachment: %v\nSMTP: %s\nThread size: %d\n", cfg.Sender, recipients, cfg.Subject, len(cfg.Body), cfg.AttachmentPath != "", cfg.SMTPServer, effectiveThread)
				summary += buildThreadHeaders(cfg, ids)
				ok, err := showSummaryTUI(summary)
				if err != nil {
					return err
				}
				if !ok {
					fmt.Println("❌ Cancelled by user")
					return nil
				}
				// send thread using the pre-generated IDs
				var prevID string
				for i := 0; i < effectiveThread; i++ {
					if i > 0 {
						cfg.InReplyTo = prevID
						cfg.References = append(cfg.References, prevID)
					}
					cfg.MessageID = ids[i]
					msg, err := sendemail.CreateRFCMessage(cfg)
					if err != nil {
						return err
					}
					if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
						return err
					}
					fmt.Printf("📨 Sent message id %s to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
					prevID = cfg.MessageID
					if delayMs > 0 {
						time.Sleep(time.Duration(delayMs) * time.Millisecond)
					}
				}
				return nil
			}
			// interactive
			cfg, err := sendemail.PromptInteractiveConfig(os.Stdin, os.Stdout)
			if err != nil {
				return err
			}
			if cfg.AttachmentPath != "" {
				fmt.Printf("📎 Attachment: %s\n", cfg.AttachmentPath)
			}
			if cfg.ThreadSize > 1 {
				fmt.Printf("🔁 Sending thread of %d messages\n", cfg.ThreadSize)
			}

			// interactive summary + confirmation (compute effective thread and ids first)
			effectiveThread := threadSize
			if effectiveThread <= 0 {
				effectiveThread = cfg.ThreadSize
			}
			if effectiveThread <= 0 {
				effectiveThread = 1
			}
			ids := make([]string, effectiveThread)
			for i := 0; i < effectiveThread; i++ {
				ids[i] = genMessageID()
			}

			recipients := strings.Join(cfg.Recipients, ",")
			if len(cfg.CCRecipients) > 0 {
				recipients += ", cc:" + strings.Join(cfg.CCRecipients, ",")
			}
			summary := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nBody length: %d chars\nAttachment: %v\nSMTP: %s\nThread size: %d\n", cfg.Sender, recipients, cfg.Subject, len(cfg.Body), cfg.AttachmentPath != "", cfg.SMTPServer, effectiveThread)
			summary += buildThreadHeaders(cfg, ids)
			ok, err := showSummaryTUI(summary)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println("❌ Cancelled by user")
				return nil
			}
			// send thread using the pre-generated IDs
			var prevID string
			for i := 0; i < effectiveThread; i++ {
				if i > 0 {
					cfg.InReplyTo = prevID
					cfg.References = append(cfg.References, prevID)
				}
				cfg.MessageID = ids[i]
				msg, err := sendemail.CreateRFCMessage(cfg)
				if err != nil {
					return err
				}
				if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
					return err
				}
				fmt.Printf("📨 Sent message id %s to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
				prevID = cfg.MessageID
				if delayMs > 0 {
					time.Sleep(time.Duration(delayMs) * time.Millisecond)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&random, "random", false, "generate random From/To/Subject/Body")
	cmd.Flags().IntVar(&threadSize, "thread-size", 0, "generate a thread of N messages (first + N-1 replies)")
	cmd.Flags().IntVar(&delayMs, "delay-ms", 0, "delay between thread messages in milliseconds")
	cmd.Flags().BoolVar(&classic, "classic", false, "use classic prompt mode (non-TUI)")

	return cmd
}

// showSummaryTUI presents the summary as a Bubble Tea table and
// returns true if the user confirms (Y or Enter), false otherwise.
func showSummaryTUI(summary string) (bool, error) {
	// parse summary lines into table rows
	lines := strings.Split(strings.TrimSpace(summary), "\n")
	// Give a very wide Value column so nothing is truncated; user can scroll.
	cols := []table.Column{
		{Title: "Field", Width: 24},
		{Title: "Value", Width: 2000},
	}
	rows := make([]table.Row, 0, len(lines))
	for _, l := range lines {
		parts := strings.SplitN(l, ":", 2)
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			rows = append(rows, table.Row{field, val})
		} else {
			rows = append(rows, table.Row{"", strings.TrimSpace(l)})
		}
	}

	tbl := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	m := &summaryModel{table: tbl, confirmed: false}

	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		return false, err
	}
	return m.confirmed, nil
}

// summaryModel is the Bubble Tea model used to render the summary table
// and capture a Y/N confirmation.
type summaryModel struct {
	table     table.Model
	confirmed bool
	xOffset   int
}

func (m *summaryModel) Init() tea.Cmd { return nil }

func (m *summaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := strings.ToLower(msg.String())
		if k == "y" || k == "enter" {
			m.confirmed = true
			return m, tea.Quit
		}
		if k == "n" || k == "esc" || k == "q" {
			m.confirmed = false
			return m, tea.Quit
		}
		// horizontal pan: left/right adjust xOffset
		if k == "left" {
			if m.xOffset >= 4 {
				m.xOffset -= 4
			} else {
				m.xOffset = 0
			}
			return m, nil
		}
		if k == "right" {
			m.xOffset += 4
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *summaryModel) View() string {
	// Render the table and apply horizontal offset by slicing each line.
	raw := m.table.View()
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		if m.xOffset >= len(line) {
			lines[i] = ""
		} else if m.xOffset > 0 {
			// slice by rune index to avoid breaking multi-byte chars
			// but for performance and simplicity assume ASCII table output
			lines[i] = line[m.xOffset:]
		}
	}
	out := strings.Join(lines, "\n")
	return out + "\n\nUse ←/→ to scroll horizontally. Press Y/Enter to confirm, N/ESC/Q to cancel"
}

// genMessageID produces a unique, RFC-like Message-ID value.
func genMessageID() string {
	host, _ := os.Hostname()
	if host == "" {
		host = "localhost"
	}
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%s@%s", strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", ""), host)
	}
	return fmt.Sprintf("%s-%s@%s", hex.EncodeToString(b), strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", ""), host)
}

// buildThreadHeaders returns a string containing lines representing the headers
// for each message in the thread. For each message it emits one row with the
// field "Email N Headers" and the first header as the value, then additional
// rows with an empty field and one header per row so the TUI shows one header
// per line in the value column.
func buildThreadHeaders(orig *sendemail.EmailConfig, ids []string) string {
	var sb strings.Builder
	// For each id, build a copy of cfg with appropriate InReplyTo/References
	for i, id := range ids {
		cfg := *orig
		// set MessageID raw (no <>)
		cfg.MessageID = id
		if i > 0 {
			cfg.InReplyTo = ids[i-1]
			// References should include all previous IDs
			cfg.References = make([]string, i)
			copy(cfg.References, ids[:i])
		} else {
			cfg.InReplyTo = ""
			cfg.References = nil
		}
		// Generate the raw message bytes and extract headers
		msg, _ := sendemail.CreateRFCMessage(&cfg)
		headerEnd := -1
		if idx := strings.Index(string(msg), "\r\n\r\n"); idx != -1 {
			headerEnd = idx
		} else if idx := strings.Index(string(msg), "\n\n"); idx != -1 {
			headerEnd = idx
		}
		var headers []string
		if headerEnd != -1 {
			rawHeaders := string(msg[:headerEnd])
			// split on CRLF or LF
			for _, line := range strings.Split(rawHeaders, "\r\n") {
				if strings.TrimSpace(line) == "" {
					continue
				}
				headers = append(headers, line)
			}
			if len(headers) == 0 {
				// fallback to LF-only
				for _, line := range strings.Split(rawHeaders, "\n") {
					if strings.TrimSpace(line) == "" {
						continue
					}
					headers = append(headers, line)
				}
			}
		}

		// Emit rows: first row with field, subsequent with empty field
		field := fmt.Sprintf("Email %d Headers", i+1)
		if len(headers) == 0 {
			sb.WriteString(fmt.Sprintf("%s: <no-headers>\n", field))
			continue
		}
		for j, h := range headers {
			if j == 0 {
				sb.WriteString(fmt.Sprintf("%s: %s\n", field, h))
			} else {
				sb.WriteString(fmt.Sprintf(": %s\n", h))
			}
		}
	}
	return sb.String()
}
