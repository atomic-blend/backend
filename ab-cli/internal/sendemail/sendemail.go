// Package sendemail provides utilities to build RFC-5322 compliant
// test emails and send them via SMTP for local development.
package sendemail

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	smtpclient "github.com/emersion/go-smtp"
)

type EmailConfig struct {
	Sender         string
	Recipients     []string
	CCRecipients   []string
	BCCRecipients  []string
	Subject        string
	Body           string
	AttachmentPath string
	SMTPServer     string
	SMTPUsername   string
	SMTPPassword   string
	MessageID      string
	InReplyTo      string
	References     []string
	Date           time.Time
	ThreadSize     int
}

func PromptInteractiveConfig(in io.Reader, out io.Writer) (*EmailConfig, error) {
	reader := bufio.NewReader(in)
	cfg := &EmailConfig{SMTPServer: "localhost:1025"}

	fmt.Fprint(out, "From (sender email) [brandon@brandonguigo.com]: ")
	senderInput, _ := reader.ReadString('\n')
	senderInput = strings.TrimSpace(senderInput)
	if senderInput == "" {
		senderInput = "brandon@brandonguigo.com"
	}
	cfg.Sender = senderInput

	fmt.Fprint(out, "To (recipient email(s), comma-separated) [user2@brandonguigo.com]: ")
	toInput, _ := reader.ReadString('\n')
	toInput = strings.TrimSpace(toInput)
	if toInput == "" {
		toInput = "user2@brandonguigo.com"
	}
	cfg.Recipients = parseEmailList(toInput)

	fmt.Fprint(out, "CC (comma-separated, or press Enter to skip): ")
	ccInput, _ := reader.ReadString('\n')
	ccInput = strings.TrimSpace(ccInput)
	if ccInput != "" {
		cfg.CCRecipients = parseEmailList(ccInput)
	}

	fmt.Fprint(out, "BCC (comma-separated, or press Enter to skip): ")
	bccInput, _ := reader.ReadString('\n')
	bccInput = strings.TrimSpace(bccInput)
	if bccInput != "" {
		cfg.BCCRecipients = parseEmailList(bccInput)
	}

	fmt.Fprint(out, "Subject: ")
	subj, _ := reader.ReadString('\n')
	cfg.Subject = strings.TrimSpace(subj)

	fmt.Fprintln(out, "Type your email body. End with a single line 'END'. Press Enter immediately to auto-generate random body:")

	// Read first line specially: if user just presses Enter (empty line),
	// generate a random body and continue to next step.
	firstLine, _ := reader.ReadString('\n')
	firstLine = strings.TrimRight(firstLine, "\r\n")
	if firstLine == "" {
		gofakeit.Seed(time.Now().UnixNano())
		cfg.Body = gofakeit.Paragraph(1, 3, 12, " ")
	} else {
		lines := []string{firstLine}
		for {
			l, _ := reader.ReadString('\n')
			l = strings.TrimRight(l, "\r\n")
			if l == "END" {
				break
			}
			lines = append(lines, l)
		}
		cfg.Body = strings.Join(lines, "\n")
	}

	// fill missing subject with random content (body handled above)
	if strings.TrimSpace(cfg.Subject) == "" {
		gofakeit.Seed(time.Now().UnixNano())
		cfg.Subject = gofakeit.Sentence(6)
	}

	// Attachments
	fmt.Fprint(out, "Do you want to attach a file? (y/n) [n]: ")
	attachInput, _ := reader.ReadString('\n')
	attachInput = strings.TrimSpace(strings.ToLower(attachInput))
	if attachInput == "y" || attachInput == "yes" {
		fmt.Fprint(out, "Enter file path to attach: ")
		pathInput, _ := reader.ReadString('\n')
		pathInput = strings.TrimSpace(pathInput)
		if pathInput != "" {
			if _, err := os.Stat(pathInput); os.IsNotExist(err) {
				return nil, fmt.Errorf("attachment does not exist: %s", pathInput)
			}
			cfg.AttachmentPath = pathInput
		}
	}

	// Threading
	fmt.Fprint(out, "Do you want to generate a thread (multiple messages)? (y/n) [n]: ")
	threadInput, _ := reader.ReadString('\n')
	threadInput = strings.TrimSpace(strings.ToLower(threadInput))
	if threadInput == "y" || threadInput == "yes" {
		fmt.Fprint(out, "Thread size (including first message) [2]: ")
		sizeInput, _ := reader.ReadString('\n')
		sizeInput = strings.TrimSpace(sizeInput)
		if sizeInput == "" {
			cfg.ThreadSize = 2
		} else {
			n, err := strconv.Atoi(sizeInput)
			if err != nil {
				return nil, fmt.Errorf("invalid thread size: %v", err)
			}
			cfg.ThreadSize = n
		}
	}

	// SMTP server
	fmt.Fprint(out, "SMTP Server [localhost:1025]: ")
	smtpInput, _ := reader.ReadString('\n')
	smtpInput = strings.TrimSpace(smtpInput)
	if smtpInput != "" {
		cfg.SMTPServer = smtpInput
	}

	// SMTP auth (optional)
	fmt.Fprint(out, "Do you need SMTP authentication? (y/n) [n]: ")
	authInput, _ := reader.ReadString('\n')
	authInput = strings.TrimSpace(strings.ToLower(authInput))
	if authInput == "y" || authInput == "yes" {
		fmt.Fprint(out, "SMTP Username: ")
		u, _ := reader.ReadString('\n')
		cfg.SMTPUsername = strings.TrimSpace(u)
		fmt.Fprint(out, "SMTP Password: ")
		p, _ := reader.ReadString('\n')
		cfg.SMTPPassword = strings.TrimSpace(p)
	}

	return cfg, nil
}

func parseEmailList(input string) []string {
	if input == "" {
		return []string{}
	}
	parts := strings.Split(input, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func CreateRFCMessage(cfg *EmailConfig) ([]byte, error) {
	if cfg.Date.IsZero() {
		cfg.Date = time.Now()
	}
	if cfg.MessageID == "" {
		cfg.MessageID = generateMessageIDFromSender(cfg.Sender)
	}

	if cfg.AttachmentPath != "" {
		return createMultipartMessage(cfg)
	}
	return createSimpleMessage(cfg), nil
}

func createSimpleMessage(cfg *EmailConfig) []byte {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("From: %s\r\n", cfg.Sender))
	if len(cfg.Recipients) > 0 {
		b.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(cfg.Recipients, ", ")))
	}
	if len(cfg.CCRecipients) > 0 {
		b.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cfg.CCRecipients, ", ")))
	}
	b.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.QEncoding.Encode("utf-8", cfg.Subject)))
	b.WriteString(fmt.Sprintf("Date: %s\r\n", formatDateRFC5322(cfg.Date)))
	b.WriteString(fmt.Sprintf("Message-ID: <%s>\r\n", cfg.MessageID))
	if cfg.InReplyTo != "" {
		b.WriteString(fmt.Sprintf("In-Reply-To: <%s>\r\n", cfg.InReplyTo))
	}
	if len(cfg.References) > 0 {
		refs := make([]string, len(cfg.References))
		for i, r := range cfg.References {
			refs[i] = fmt.Sprintf("<%s>", r)
		}
		b.WriteString(fmt.Sprintf("References: %s\r\n", strings.Join(refs, " ")))
	}
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(cfg.Body)
	return []byte(b.String())
}

func createMultipartMessage(cfg *EmailConfig) ([]byte, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	boundary := writer.Boundary()

	buffer.WriteString(fmt.Sprintf("From: %s\r\n", cfg.Sender))
	if len(cfg.Recipients) > 0 {
		buffer.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(cfg.Recipients, ", ")))
	}
	if len(cfg.CCRecipients) > 0 {
		buffer.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cfg.CCRecipients, ", ")))
	}
	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.QEncoding.Encode("utf-8", cfg.Subject)))
	buffer.WriteString(fmt.Sprintf("Date: %s\r\n", formatDateRFC5322(cfg.Date)))
	buffer.WriteString(fmt.Sprintf("Message-ID: <%s>\r\n", cfg.MessageID))
	if cfg.InReplyTo != "" {
		buffer.WriteString(fmt.Sprintf("In-Reply-To: <%s>\r\n", cfg.InReplyTo))
	}
	if len(cfg.References) > 0 {
		refs := make([]string, len(cfg.References))
		for i, r := range cfg.References {
			refs[i] = fmt.Sprintf("<%s>", r)
		}
		buffer.WriteString(fmt.Sprintf("References: %s\r\n", strings.Join(refs, " ")))
	}
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	buffer.WriteString("\r\n")

	// text part
	buffer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buffer.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buffer.WriteString("\r\n")
	buffer.WriteString(cfg.Body)
	buffer.WriteString("\r\n")

	// attachment
	file, err := os.Open(cfg.AttachmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open attachment: %v", err)
	}
	defer file.Close()

	fileName := filepath.Base(cfg.AttachmentPath)
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	buffer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buffer.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimeType))
	buffer.WriteString("Content-Transfer-Encoding: base64\r\n")
	buffer.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
	buffer.WriteString("\r\n")

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read attachment: %v", err)
	}
	encoded := base64.StdEncoding.EncodeToString(fileData)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		buffer.WriteString(encoded[i:end])
		buffer.WriteString("\r\n")
	}

	buffer.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	return buffer.Bytes(), nil
}

func SendUsingSMTP(cfg *EmailConfig, message []byte) error {
	client, err := smtpclient.Dial(cfg.SMTPServer)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Mail(cfg.Sender, nil); err != nil {
		return err
	}
	allRecipients := append(cfg.Recipients, cfg.CCRecipients...)
	allRecipients = append(allRecipients, cfg.BCCRecipients...)
	for _, r := range allRecipients {
		if err := client.Rcpt(r, nil); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(message)
	if err != nil {
		return err
	}
	return w.Close()
}
