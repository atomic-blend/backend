// Package sendemailcmd provides the `send-email` developper subcommand.
package sendemailcmd

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	sendemail "github.com/atomic-blend/backend/ab-cli/internal/sendemail"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	ical "github.com/emersion/go-ical"
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
				// determine effective thread size and pre-generate IDs
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

				// If composer requested a calendar invite, generate invite + reply .ics files
				var invitePath string
				var replyPaths []string
				var attendees []string
				var uid string
				var replyMethods []string
				if cfg.IncludeCalendar {
					attendees = buildAttendeesList(cfg, effectiveThread)
					// create invite (METHOD=REQUEST)
					p, uid, err := generateICalInvite(cfg, attendees)
					if err != nil {
						return err
					}
					invitePath = p
					replyPaths = make([]string, effectiveThread)
					replyMethods = make([]string, effectiveThread)
					// generate per-attendee replies for messages 1..N-1
					for j := 1; j < effectiveThread; j++ {
						method := randomReplyMethod()
						rp, err := generateICalReply(cfg, uid, attendees[j], method)
						if err != nil {
							return err
						}
						replyPaths[j] = rp
						replyMethods[j] = method
					}
					defer func() {
						_ = os.Remove(invitePath)
						for _, r := range replyPaths {
							if r != "" {
								_ = os.Remove(r)
							}
						}
					}()
				}

				if invitePath != "" {
					fmt.Printf("📎 Calendar invite will be sent (uid: %s)\n", uid)
				} else if cfg.AttachmentPath != "" {
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
				// include calendar details when present
				if invitePath != "" {
					calInfo := fmt.Sprintf("Calendar UID: %s\nAttendees: %s\n", uid, strings.Join(attendees, ","))
					for i := 0; i < effectiveThread; i++ {
						if i == 0 {
							ctype, brief := summarizeICal(invitePath)
							calInfo += fmt.Sprintf("Email %d Attachment: invite (Type=%s) %s\n", i+1, ctype, brief)
						} else {
							m := replyMethods[i]
							if m == "" {
								m = "REPLY"
							}
							p := replyPaths[i]
							ctype, brief := summarizeICal(p)
							calInfo += fmt.Sprintf("Email %d Calendar Reply: %s (from %s) Type=%s %s\n", i+1, m, attendees[i], ctype, brief)
						}
					}
					summary += calInfo
				}
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
				origSender := cfg.Sender
				origRecipients := append([]string{}, cfg.Recipients...)
				origCC := append([]string{}, cfg.CCRecipients...)
				origBCC := append([]string{}, cfg.BCCRecipients...)
				origSubject := cfg.Subject
				for i := 0; i < effectiveThread; i++ {
					if i > 0 {
						cfg.InReplyTo = prevID
						cfg.References = append(cfg.References, prevID)
					}
					cfg.MessageID = ids[i]
					// attach appropriate calendar file and, for replies, set sender/recipient
					if invitePath != "" {
						if i == 0 {
							cfg.AttachmentPath = invitePath
							cfg.Sender = origSender
							cfg.Recipients = origRecipients
							cfg.CCRecipients = origCC
							cfg.BCCRecipients = origBCC
							cfg.Subject = origSubject
						} else {
							// send reply from attendee[i] to original sender
							att := attendees[i]
							cfg.Sender = att
							cfg.Recipients = []string{origSender}
							cfg.CCRecipients = nil
							cfg.BCCRecipients = nil
							cfg.AttachmentPath = replyPaths[i]
							cfg.Subject = "Re: " + origSubject
						}
					}
					msg, err := sendemail.CreateRFCMessage(cfg)
					if err != nil {
						return err
					}
					if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
						return err
					}
					fmt.Printf("📨 Sent message id %s to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
					prevID = cfg.MessageID
					// restore original sender/recipients/subject for next iteration if needed
					cfg.Sender = origSender
					cfg.Recipients = append([]string{}, origRecipients...)
					cfg.CCRecipients = append([]string{}, origCC...)
					cfg.BCCRecipients = append([]string{}, origBCC...)
					cfg.Subject = origSubject
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
				// Ask user if they want to include a calendar invitation
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Include calendar invitation (.ics)? (y/n) [n]: ")
				choice, _ := reader.ReadString('\n')
				choice = strings.TrimSpace(strings.ToLower(choice))
				if choice == "y" || choice == "yes" {
					cfg.IncludeCalendar = true
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

				// If calendar requested, generate invite + replies
				var invitePath string
				var replyPaths []string
				var attendees []string
				var uid string
				var replyMethods []string
				if cfg.IncludeCalendar {
					attendees = buildAttendeesList(cfg, effectiveThread)
					p, uid, err := generateICalInvite(cfg, attendees)
					if err != nil {
						return err
					}
					invitePath = p
					replyPaths = make([]string, effectiveThread)
					replyMethods = make([]string, effectiveThread)
					for j := 1; j < effectiveThread; j++ {
						method := randomReplyMethod()
						rp, err := generateICalReply(cfg, uid, attendees[j], method)
						if err != nil {
							return err
						}
						replyPaths[j] = rp
						replyMethods[j] = method
					}
					defer func() {
						_ = os.Remove(invitePath)
						for _, r := range replyPaths {
							if r != "" {
								_ = os.Remove(r)
							}
						}
					}()
					fmt.Printf("📎 Calendar invite will be sent (uid: %s)\n", uid)
				}

				// build and show summary
				recipients := strings.Join(cfg.Recipients, ",")
				if len(cfg.CCRecipients) > 0 {
					recipients += ", cc:" + strings.Join(cfg.CCRecipients, ",")
				}
				summary := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nBody length: %d chars\nAttachment: %v\nSMTP: %s\nThread size: %d\n", cfg.Sender, recipients, cfg.Subject, len(cfg.Body), cfg.AttachmentPath != "", cfg.SMTPServer, effectiveThread)
				if invitePath != "" {
					calInfo := fmt.Sprintf("Calendar UID: %s\nAttendees: %s\n", uid, strings.Join(attendees, ","))
					for i := 0; i < effectiveThread; i++ {
						if i == 0 {
							ctype, brief := summarizeICal(invitePath)
							calInfo += fmt.Sprintf("Email %d Attachment: invite (Type=%s) %s\n", i+1, ctype, brief)
						} else {
							m := replyMethods[i]
							if m == "" {
								m = "REPLY"
							}
							p := replyPaths[i]
							ctype, brief := summarizeICal(p)
							calInfo += fmt.Sprintf("Email %d Calendar Reply: %s (from %s) Type=%s %s\n", i+1, m, attendees[i], ctype, brief)
						}
					}
					summary += calInfo
				}
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

			// Ask user if they want to include a calendar invitation
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Include calendar invitation (.ics)? (y/n) [n]: ")
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(strings.ToLower(choice))
			if choice == "y" || choice == "yes" {
				cfg.IncludeCalendar = true
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

			// If calendar requested, generate invite + replies
			var invitePath string
			var replyPaths []string
			var attendees []string
			var uid string
			var replyMethods []string
			if cfg.IncludeCalendar {
				attendees = buildAttendeesList(cfg, effectiveThread)
				p, uid, err := generateICalInvite(cfg, attendees)
				if err != nil {
					return err
				}
				invitePath = p
				replyPaths = make([]string, effectiveThread)
				replyMethods = make([]string, effectiveThread)
				for j := 1; j < effectiveThread; j++ {
					method := randomReplyMethod()
					rp, err := generateICalReply(cfg, uid, attendees[j], method)
					if err != nil {
						return err
					}
					replyPaths[j] = rp
					replyMethods[j] = method
				}
				defer func() {
					_ = os.Remove(invitePath)
					for _, r := range replyPaths {
						if r != "" {
							_ = os.Remove(r)
						}
					}
				}()
				fmt.Printf("📎 Calendar invite will be sent (uid: %s)\n", uid)
			}

			recipients := strings.Join(cfg.Recipients, ",")
			if len(cfg.CCRecipients) > 0 {
				recipients += ", cc:" + strings.Join(cfg.CCRecipients, ",")
			}
			summary := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nBody length: %d chars\nAttachment: %v\nSMTP: %s\nThread size: %d\n", cfg.Sender, recipients, cfg.Subject, len(cfg.Body), cfg.AttachmentPath != "", cfg.SMTPServer, effectiveThread)
			if invitePath != "" {
				calInfo := fmt.Sprintf("Calendar UID: %s\nAttendees: %s\n", uid, strings.Join(attendees, ","))
				for i := 0; i < effectiveThread; i++ {
					if i == 0 {
						ctype, brief := summarizeICal(invitePath)
						calInfo += fmt.Sprintf("Email %d Attachment: invite (Type=%s) %s\n", i+1, ctype, brief)
					} else {
						m := replyMethods[i]
						if m == "" {
							m = "REPLY"
						}
						p := replyPaths[i]
						ctype, brief := summarizeICal(p)
						calInfo += fmt.Sprintf("Email %d Calendar Reply: %s (from %s) Type=%s %s\n", i+1, m, attendees[i], ctype, brief)
					}
				}
				summary += calInfo
			}
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
			origSender := cfg.Sender
			origRecipients := append([]string{}, cfg.Recipients...)
			origCC := append([]string{}, cfg.CCRecipients...)
			origBCC := append([]string{}, cfg.BCCRecipients...)
			for i := 0; i < effectiveThread; i++ {
				if i > 0 {
					cfg.InReplyTo = prevID
					cfg.References = append(cfg.References, prevID)
				}
				cfg.MessageID = ids[i]
				// attach appropriate calendar file and, for replies, set sender/recipient
				if invitePath != "" {
					if i == 0 {
						cfg.AttachmentPath = invitePath
						cfg.Sender = origSender
						cfg.Recipients = origRecipients
						cfg.CCRecipients = origCC
						cfg.BCCRecipients = origBCC
					} else {
						// send reply from attendee[i] to original sender
						att := attendees[i]
						cfg.Sender = att
						cfg.Recipients = []string{origSender}
						cfg.CCRecipients = nil
						cfg.BCCRecipients = nil
						cfg.AttachmentPath = replyPaths[i]
					}
				}
				msg, err := sendemail.CreateRFCMessage(cfg)
				if err != nil {
					return err
				}
				if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
					return err
				}
				fmt.Printf("📨 Sent message id %s to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
				prevID = cfg.MessageID
				// restore original sender/recipients for next iteration if needed
				cfg.Sender = origSender
				cfg.Recipients = append([]string{}, origRecipients...)
				cfg.CCRecipients = append([]string{}, origCC...)
				cfg.BCCRecipients = append([]string{}, origBCC...)
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

// buildAttendeesList returns a list of attendees with length equal to threadSize.
// It prefers provided recipients and fills remaining slots with generated addresses.
func buildAttendeesList(cfg *sendemail.EmailConfig, threadSize int) []string {
	out := make([]string, threadSize)
	// derive domain from sender
	domain := "example.com"
	if parts := strings.Split(cfg.Sender, "@"); len(parts) == 2 {
		domain = parts[1]
	}
	for i := 0; i < threadSize; i++ {
		if i < len(cfg.Recipients) {
			out[i] = cfg.Recipients[i]
			continue
		}
		out[i] = fmt.Sprintf("attendee+%d@%s", i+1, domain)
	}
	return out
}

// generateICalInvite creates an .ics file for an initial invite (METHOD=REQUEST)
// and returns the file path and generated UID.
func generateICalInvite(cfg *sendemail.EmailConfig, attendees []string) (string, string, error) {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//atomic-blend//EN")
	cal.Props.SetText(ical.PropMethod, "REQUEST")

	ev := ical.NewEvent()
	uid := genMessageID()
	ev.Props.SetText(ical.PropUID, uid)
	now := time.Now().UTC()
	ev.Props.SetDateTime(ical.PropDateTimeStamp, now)
	start := now.Add(1 * time.Hour)
	ev.Props.SetDateTime(ical.PropDateTimeStart, start)
	ev.Props.SetDateTime(ical.PropDateTimeEnd, start.Add(1*time.Hour))
	ev.Props.SetText(ical.PropSummary, cfg.Subject)
	ev.Props.SetText(ical.PropDescription, cfg.Body)
	ev.Props.SetText(ical.PropOrganizer, "MAILTO:"+cfg.Sender)
	for _, r := range attendees {
		p := ical.NewProp(ical.PropAttendee)
		p.SetText("MAILTO:" + r)
		ev.Props.Add(p)
	}

	cal.Children = append(cal.Children, ev.Component)

	var buf bytes.Buffer
	enc := ical.NewEncoder(&buf)
	if err := enc.Encode(cal); err != nil {
		return "", "", err
	}

	f, err := os.CreateTemp("", "invite-*.ics")
	if err != nil {
		return "", "", err
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		f.Close()
		return "", "", err
	}
	if err := f.Close(); err != nil {
		return "", "", err
	}
	return f.Name(), uid, nil
}

// generateICalReply creates an .ics file for a reply/cancel/publish with the
// given UID and attendee. method should be one of "REPLY","CANCEL","PUBLISH".
func generateICalReply(cfg *sendemail.EmailConfig, uid, attendee, method string) (string, error) {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//atomic-blend//EN")
	cal.Props.SetText(ical.PropMethod, method)

	ev := ical.NewEvent()
	ev.Props.SetText(ical.PropUID, uid)
	now := time.Now().UTC()
	ev.Props.SetDateTime(ical.PropDateTimeStamp, now)
	// replies typically include DTSTART/DTEND as reference to original
	start := now.Add(1 * time.Hour)
	ev.Props.SetDateTime(ical.PropDateTimeStart, start)
	ev.Props.SetDateTime(ical.PropDateTimeEnd, start.Add(1*time.Hour))
	ev.Props.SetText(ical.PropSummary, cfg.Subject)
	ev.Props.SetText(ical.PropOrganizer, "MAILTO:"+cfg.Sender)

	// attendee for this reply
	p := ical.NewProp(ical.PropAttendee)
	p.SetText("MAILTO:" + attendee)
	// for replies/cancels, set PARTSTAT and optionally STATUS
	switch method {
	case "CANCEL":
		ev.Props.SetText(ical.PropStatus, "CANCELLED")
		// mark as declined
		if p.Params == nil {
			p.Params = ical.Params{}
		}
		p.Params.Set(ical.ParamParticipationStatus, "DECLINED")
	case "REPLY":
		if p.Params == nil {
			p.Params = ical.Params{}
		}
		// random accepted/tentative
		if time.Now().UnixNano()%2 == 0 {
			p.Params.Set(ical.ParamParticipationStatus, "ACCEPTED")
		} else {
			p.Params.Set(ical.ParamParticipationStatus, "TENTATIVE")
		}
	}
	ev.Props.Add(p)

	cal.Children = append(cal.Children, ev.Component)

	var buf bytes.Buffer
	enc := ical.NewEncoder(&buf)
	if err := enc.Encode(cal); err != nil {
		return "", err
	}

	f, err := os.CreateTemp("", "reply-*.ics")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

// randomReplyMethod returns either "REPLY" or "CANCEL" at random.
func randomReplyMethod() string {
	if time.Now().UnixNano()%2 == 0 {
		return "REPLY"
	}
	return "CANCEL"
}

// summarizeICal reads an .ics file and returns a short content type (METHOD)
// and a brief summary (SUMMARY or DESCRIPTION plus attendees) for display.
func summarizeICal(path string) (string, string) {
	if path == "" {
		return "", ""
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", ""
	}
	s := string(b)
	lines := strings.Split(s, "\n")
	var method, summary string
	var attendees []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if method == "" && strings.HasPrefix(l, "METHOD:") {
			method = strings.TrimSpace(strings.TrimPrefix(l, "METHOD:"))
			continue
		}
		if summary == "" && strings.HasPrefix(l, "SUMMARY:") {
			summary = strings.TrimSpace(strings.TrimPrefix(l, "SUMMARY:"))
			continue
		}
		if summary == "" && strings.HasPrefix(l, "DESCRIPTION:") {
			summary = strings.TrimSpace(strings.TrimPrefix(l, "DESCRIPTION:"))
			continue
		}
		if strings.HasPrefix(l, "ATTENDEE:") {
			a := strings.TrimSpace(strings.TrimPrefix(l, "ATTENDEE:"))
			a = strings.TrimPrefix(a, "MAILTO:")
			attendees = append(attendees, a)
		}
	}
	if method == "" {
		method = "UNKNOWN"
	}
	brief := summary
	if brief == "" {
		// fallback to first 80 characters of file
		if len(s) > 80 {
			brief = strings.TrimSpace(s[:80]) + "..."
		} else {
			brief = strings.TrimSpace(s)
		}
	}
	if len(attendees) > 0 {
		// show up to 3 attendees
		n := 3
		if len(attendees) < n {
			n = len(attendees)
		}
		brief += " (Attendees: " + strings.Join(attendees[:n], ",") + ")"
	}
	return method, brief
}
