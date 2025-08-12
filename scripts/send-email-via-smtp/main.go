package main

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
	"strings"
	"time"

	"github.com/emersion/go-smtp"
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
}

type EmailProgress struct {
	Step     string
	Status   string
	Duration time.Duration
	Details  string
}

func main() {
	fmt.Println("🚀 Interactive SMTP Email Sender")
	fmt.Println("📍 For local development and testing")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	config := &EmailConfig{
		SMTPServer: "localhost:1025",
	}

	// Interactive prompts with better formatting
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Sender
	fmt.Println("📧 STEP 1: Sender Information")
	fmt.Println(strings.Repeat("-", 30))
	senderInput := promptWithDefault(reader, "From (sender email)", "brandon@brandonguigo.com")
	config.Sender = senderInput
	fmt.Printf("✅ Sender set to: %s\n\n", config.Sender)

	// Step 2: Recipients
	fmt.Println("👥 STEP 2: Recipients")
	fmt.Println(strings.Repeat("-", 30))
	recipientsInput := promptWithDefault(reader, "To (recipient email(s), comma-separated)", "user2@brandonguigo.com")
	config.Recipients = parseEmailList(recipientsInput)
	fmt.Printf("✅ Recipients set to: %s\n\n", strings.Join(config.Recipients, ", "))

	// Step 3: CC Recipients
	fmt.Println("📋 STEP 3: CC Recipients (Optional)")
	fmt.Println(strings.Repeat("-", 30))
	ccInput := promptWithDefault(reader, "CC (comma-separated, or press Enter to skip)", "")
	if ccInput != "" {
		config.CCRecipients = parseEmailList(ccInput)
		fmt.Printf("✅ CC recipients set to: %s\n", strings.Join(config.CCRecipients, ", "))
	} else {
		fmt.Println("⏭️  Skipping CC recipients")
	}
	fmt.Println()

	// Step 4: BCC Recipients
	fmt.Println("🔒 STEP 4: BCC Recipients (Optional)")
	fmt.Println(strings.Repeat("-", 30))
	bccInput := promptWithDefault(reader, "BCC (comma-separated, or press Enter to skip)", "")
	if bccInput != "" {
		config.BCCRecipients = parseEmailList(bccInput)
		fmt.Printf("✅ BCC recipients set to: %s\n", strings.Join(config.BCCRecipients, ", "))
	} else {
		fmt.Println("⏭️  Skipping BCC recipients")
	}
	fmt.Println()

	// Step 5: Subject
	fmt.Println("📝 STEP 5: Email Subject")
	fmt.Println(strings.Repeat("-", 30))
	config.Subject = promptRequired(reader, "Subject")
	fmt.Printf("✅ Subject set to: %s\n\n", config.Subject)

	// Step 6: Email Body
	fmt.Println("💬 STEP 6: Email Body")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("Type your email body below. Type 'END' on a new line to finish:")
	fmt.Println("💡 Tip: You can write multiple lines. Type 'END' when done.")
	fmt.Println()

	var bodyLines []string
	lineNumber := 1
	for {
		fmt.Printf("Line %d: ", lineNumber)
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		if line == "END" {
			break
		}
		bodyLines = append(bodyLines, line)
		lineNumber++
	}
	config.Body = strings.Join(bodyLines, "\n")
	fmt.Printf("✅ Email body completed (%d lines)\n\n", len(bodyLines))

	// Step 7: Attachments
	fmt.Println("📎 STEP 7: File Attachments (Optional)")
	fmt.Println(strings.Repeat("-", 30))
	attachInput := promptWithDefault(reader, "Do you want to attach a file? (y/n)", "n")
	if strings.ToLower(attachInput) == "y" || strings.ToLower(attachInput) == "yes" {
		config.AttachmentPath = promptRequired(reader, "Enter file path to attach")

		// Validate file exists
		if _, err := os.Stat(config.AttachmentPath); os.IsNotExist(err) {
			fmt.Printf("❌ Error: File does not exist: %s\n", config.AttachmentPath)
			return
		}

		// Get file info
		fileInfo, _ := os.Stat(config.AttachmentPath)
		fmt.Printf("✅ Attachment: %s (%.2f KB)\n", filepath.Base(config.AttachmentPath), float64(fileInfo.Size())/1024)
	} else {
		fmt.Println("⏭️  Skipping attachments")
	}
	fmt.Println()

	// Step 8: SMTP Configuration
	fmt.Println("⚙️  STEP 8: SMTP Configuration")
	fmt.Println(strings.Repeat("-", 30))
	smtpInput := promptWithDefault(reader, "SMTP Server", "localhost:1025")
	if smtpInput != "" {
		config.SMTPServer = smtpInput
	}
	fmt.Printf("✅ SMTP Server: %s\n\n", config.SMTPServer)

	// Step 9: SMTP Authentication (Optional)
	fmt.Println("🔐 STEP 9: SMTP Authentication (Optional)")
	fmt.Println(strings.Repeat("-", 30))
	authInput := promptWithDefault(reader, "Do you need SMTP authentication? (y/n)", "n")
	if strings.ToLower(authInput) == "y" || strings.ToLower(authInput) == "yes" {
		config.SMTPUsername = promptRequired(reader, "SMTP Username")
		config.SMTPPassword = promptPassword(reader, "SMTP Password")
		fmt.Println("✅ Authentication credentials set")
	} else {
		fmt.Println("⏭️  Skipping authentication")
	}
	fmt.Println()

	// Summary
	fmt.Println("📋 EMAIL SUMMARY")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("📧 From: %s\n", config.Sender)
	fmt.Printf("👥 To: %s\n", strings.Join(config.Recipients, ", "))
	if len(config.CCRecipients) > 0 {
		fmt.Printf("📋 CC: %s\n", strings.Join(config.CCRecipients, ", "))
	}
	if len(config.BCCRecipients) > 0 {
		fmt.Printf("🔒 BCC: %s\n", strings.Join(config.BCCRecipients, ", "))
	}
	fmt.Printf("📝 Subject: %s\n", config.Subject)
	fmt.Printf("💬 Body: %d lines\n", len(strings.Split(config.Body, "\n")))
	if config.AttachmentPath != "" {
		fmt.Printf("📎 Attachment: %s\n", filepath.Base(config.AttachmentPath))
	}
	fmt.Printf("⚙️  SMTP Server: %s\n", config.SMTPServer)
	if config.SMTPUsername != "" {
		fmt.Printf("🔐 Username: %s\n", config.SMTPUsername)
	}
	fmt.Println()

	// Confirm before sending
	confirmInput := promptWithDefault(reader, "Ready to send? (y/n)", "y")
	if strings.ToLower(confirmInput) != "y" && strings.ToLower(confirmInput) != "yes" {
		fmt.Println("❌ Email sending cancelled")
		return
	}

	// Send email with progress
	fmt.Println("\n🚀 SENDING EMAIL")
	fmt.Println(strings.Repeat("=", 50))

	startTime := time.Now()

	// Step 1: Creating message
	fmt.Print("📝 Creating email message... ")
	message, err := createEmailMessage(config)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅")

	// Step 2: Connecting to SMTP
	fmt.Printf("🔌 Connecting to %s... ", config.SMTPServer)
	client, err := smtp.Dial(config.SMTPServer)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer client.Close()
	fmt.Println("✅")

	// Step 3: Setting sender
	fmt.Print("📤 Setting sender... ")
	if err := client.Mail(config.Sender, nil); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅")

	// Step 4: Setting recipients
	fmt.Print("👥 Setting recipients... ")
	allRecipients := append(config.Recipients, config.CCRecipients...)
	allRecipients = append(allRecipients, config.BCCRecipients...)

	for _, recipient := range allRecipients {
		if err := client.Rcpt(recipient, nil); err != nil {
			fmt.Printf("❌ Error setting recipient %s: %v\n", recipient, err)
			return
		}
	}
	fmt.Println("✅")

	// Step 5: Sending message
	fmt.Print("📨 Sending message... ")
	writer, err := client.Data()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	_, err = writer.Write(message)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅")

	duration := time.Since(startTime)

	// Success summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🎉 EMAIL SENT SUCCESSFULLY!")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("📧 Message ID: Generated\n")
	fmt.Printf("👥 Sent to: %d recipients\n", len(allRecipients))
	attachmentStatus := "No"
	if config.AttachmentPath != "" {
		attachmentStatus = "Yes"
	}
	fmt.Printf("📎 Attachment: %s\n", attachmentStatus)
	fmt.Printf("📊 Message size: %.2f KB\n", float64(len(message))/1024)
	fmt.Printf("⏱️  Total time: %.2fs\n", duration.Seconds())
	fmt.Printf("🚀 SMTP Server: %s\n", config.SMTPServer)
	fmt.Println()
	fmt.Println("✅ Your email has been sent successfully!")
}

func promptWithDefault(reader *bufio.Reader, prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func promptRequired(reader *bufio.Reader, prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
		fmt.Println("❌ This field is required. Please try again.")
	}
}

func promptPassword(reader *bufio.Reader, prompt string) string {
	fmt.Printf("%s: ", prompt)
	// For now, just read the password (in production, you'd want to hide input)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func parseEmailList(input string) []string {
	if input == "" {
		return []string{}
	}
	emails := strings.Split(input, ",")
	for i, email := range emails {
		emails[i] = strings.TrimSpace(email)
	}
	return emails
}

func createEmailMessage(config *EmailConfig) ([]byte, error) {
	if config.AttachmentPath != "" {
		return createMultipartMessage(config)
	}
	return createSimpleMessage(config), nil
}

func createSimpleMessage(config *EmailConfig) []byte {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("From: %s\r\n", config.Sender))
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(config.Recipients, ", ")))

	if len(config.CCRecipients) > 0 {
		message.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(config.CCRecipients, ", ")))
	}

	message.WriteString(fmt.Sprintf("Subject: %s\r\n", config.Subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("\r\n")
	message.WriteString(config.Body)

	return []byte(message.String())
}

func createMultipartMessage(config *EmailConfig) ([]byte, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Generate boundary
	boundary := writer.Boundary()

	// Write headers
	buffer.WriteString(fmt.Sprintf("From: %s\r\n", config.Sender))
	buffer.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(config.Recipients, ", ")))

	if len(config.CCRecipients) > 0 {
		buffer.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(config.CCRecipients, ", ")))
	}

	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", config.Subject))
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	buffer.WriteString("\r\n")

	// Text part
	buffer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buffer.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buffer.WriteString("\r\n")
	buffer.WriteString(config.Body)
	buffer.WriteString("\r\n")

	// Attachment part
	file, err := os.Open(config.AttachmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open attachment: %v", err)
	}
	defer file.Close()

	fileName := filepath.Base(config.AttachmentPath)
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	buffer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buffer.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimeType))
	buffer.WriteString("Content-Transfer-Encoding: base64\r\n")
	buffer.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
	buffer.WriteString("\r\n")

	// Read and encode file
	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read attachment: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(fileData)

	// Write encoded data in chunks of 76 characters (RFC 2045)
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
