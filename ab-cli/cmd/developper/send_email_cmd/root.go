// Package sendemailcmd provides the `send-email` developper subcommand.
package sendemailcmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	sendemail "github.com/atomic-blend/backend/ab-cli/internal/sendemail"
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
			// show summary and confirm before sending
			confirm := func(summary string) (bool, error) {
				fmt.Print(summary)
				fmt.Print("Ready to send? (y/n) [y]: ")
				r := bufio.NewReader(os.Stdin)
				in, err := r.ReadString('\n')
				if err != nil {
					return false, err
				}
				in = strings.TrimSpace(strings.ToLower(in))
				if in == "" || in == "y" || in == "yes" {
					return true, nil
				}
				return false, nil
			}

			// Default behavior: run Bubble Tea composer UI unless --classic is set
			if !classic {
				cfg, err := RunComposer()
				if err != nil {
					return err
				}
				// after composer, send as interactive flow
				if cfg.AttachmentPath != "" {
					fmt.Printf("📎 Attachment: %s\n", cfg.AttachmentPath)
				}
				if cfg.ThreadSize > 1 {
					fmt.Printf("🔁 Sending thread of %d messages\n", cfg.ThreadSize)
				}
				recipients := strings.Join(cfg.Recipients, ",")
				if len(cfg.CCRecipients) > 0 {
					recipients += ", cc:" + strings.Join(cfg.CCRecipients, ",")
				}
				summary := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nBody length: %d chars\nAttachment: %v\nSMTP: %s\nThread size: %d\n", cfg.Sender, recipients, cfg.Subject, len(cfg.Body), cfg.AttachmentPath != "", cfg.SMTPServer, cfg.ThreadSize)
				// reuse confirm helper
				ok, err := func(summary string) (bool, error) {
					fmt.Print(summary)
					fmt.Print("Ready to send? (y/n) [y]: ")
					r := bufio.NewReader(os.Stdin)
					in, err := r.ReadString('\n')
					if err != nil {
						return false, err
					}
					in = strings.TrimSpace(strings.ToLower(in))
					if in == "" || in == "y" || in == "yes" {
						return true, nil
					}
					return false, nil
				}(summary)
				if err != nil {
					return err
				}
				if !ok {
					fmt.Println("❌ Cancelled by user")
					return nil
				}
				// send thread or single
				if cfg.ThreadSize > 1 {
					var prevID string
					for i := 0; i < cfg.ThreadSize; i++ {
						if i > 0 {
							cfg.InReplyTo = prevID
							cfg.References = append(cfg.References, prevID)
						}
						if cfg.MessageID == "" {
							cfg.MessageID = strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "")
						}
						msg, err := sendemail.CreateRFCMessage(cfg)
						if err != nil {
							return err
						}
						if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
							return err
						}
						fmt.Printf("📨 Sent message id <%s> to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
						prevID = cfg.MessageID
						if delayMs > 0 {
							time.Sleep(time.Duration(delayMs) * time.Millisecond)
						}
					}
					return nil
				}
				msg, err := sendemail.CreateRFCMessage(cfg)
				if err != nil {
					return err
				}
				if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
					return err
				}
				fmt.Printf("📨 Sent message id <%s> to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
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
				// build and show summary
				recipients := strings.Join(cfg.Recipients, ",")
				if len(cfg.CCRecipients) > 0 {
					recipients += ", cc:" + strings.Join(cfg.CCRecipients, ",")
				}
				summary := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nBody length: %d chars\nAttachment: %v\nSMTP: %s\nThread size: %d\n", cfg.Sender, recipients, cfg.Subject, len(cfg.Body), cfg.AttachmentPath != "", cfg.SMTPServer, threadSize)
				ok, err := confirm(summary)
				if err != nil {
					return err
				}
				if !ok {
					fmt.Println("❌ Cancelled by user")
					return nil
				}
				// handle threading
				var prevID string
				for i := 0; i < threadSize; i++ {
					if i > 0 {
						cfg.InReplyTo = prevID
						cfg.References = append(cfg.References, prevID)
					}
					if cfg.MessageID == "" {
						cfg.MessageID = strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "")
					}
					msg, err := sendemail.CreateRFCMessage(cfg)
					if err != nil {
						return err
					}
					if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
						return err
					}
					fmt.Printf("📨 Sent message id <%s> to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
					prevID = cfg.MessageID
					if delayMs > 0 {
						time.Sleep(time.Duration(delayMs) * time.Millisecond)
					}
				}
				if threadSize == 0 {
					// single send
					msg, err := sendemail.CreateRFCMessage(cfg)
					if err != nil {
						return err
					}
					if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
						return err
					}
					fmt.Printf("📨 Sent message id <%s> to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
					return nil
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

			// interactive summary + confirmation
			recipients := strings.Join(cfg.Recipients, ",")
			if len(cfg.CCRecipients) > 0 {
				recipients += ", cc:" + strings.Join(cfg.CCRecipients, ",")
			}
			summary := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nBody length: %d chars\nAttachment: %v\nSMTP: %s\nThread size: %d\n", cfg.Sender, recipients, cfg.Subject, len(cfg.Body), cfg.AttachmentPath != "", cfg.SMTPServer, cfg.ThreadSize)
			ok, err := confirm(summary)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println("❌ Cancelled by user")
				return nil
			}

			if cfg.ThreadSize > 1 {
				var prevID string
				for i := 0; i < cfg.ThreadSize; i++ {
					if i > 0 {
						cfg.InReplyTo = prevID
						cfg.References = append(cfg.References, prevID)
					}
					if cfg.MessageID == "" {
						cfg.MessageID = strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "")
					}
					msg, err := sendemail.CreateRFCMessage(cfg)
					if err != nil {
						return err
					}
					if err := sendemail.SendUsingSMTP(cfg, msg); err != nil {
						return err
					}
					fmt.Printf("📨 Sent message id <%s> to %d recipients\n", cfg.MessageID, len(cfg.Recipients)+len(cfg.CCRecipients)+len(cfg.BCCRecipients))
					prevID = cfg.MessageID
					if delayMs > 0 {
						time.Sleep(time.Duration(delayMs) * time.Millisecond)
					}
				}
				return nil
			}

			msg, err := sendemail.CreateRFCMessage(cfg)
			if err != nil {
				return err
			}
			return sendemail.SendUsingSMTP(cfg, msg)
		},
	}

	cmd.Flags().BoolVar(&random, "random", false, "generate random From/To/Subject/Body")
	cmd.Flags().IntVar(&threadSize, "thread-size", 0, "generate a thread of N messages (first + N-1 replies)")
	cmd.Flags().IntVar(&delayMs, "delay-ms", 0, "delay between thread messages in milliseconds")
	cmd.Flags().BoolVar(&classic, "classic", false, "use classic prompt mode (non-TUI)")

	return cmd
}
