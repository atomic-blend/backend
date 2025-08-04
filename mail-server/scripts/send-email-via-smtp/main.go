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

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

func main() {
	// SMTP configuration
	smtpServer := "localhost:1025"

	// Interactive prompts
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Interactive SMTP Email Sender ===")
	fmt.Println()

	// Get sender email
	fmt.Print("From (sender email) [brandon@brandonguigo.com]: ")
	senderInput, _ := reader.ReadString('\n')
	senderInput = strings.TrimSpace(senderInput)
	sender := "brandon@brandonguigo.com"
	if senderInput != "" {
		sender = senderInput
	}

	// Get recipient(s)
	fmt.Print("To (recipient email(s), comma-separated for multiple) [user1@brandonguigo.com]: ")
	recipientsInput, _ := reader.ReadString('\n')
	recipientsInput = strings.TrimSpace(recipientsInput)
	var recipients []string
	if recipientsInput != "" {
		recipients = strings.Split(recipientsInput, ",")
		for i, r := range recipients {
			recipients[i] = strings.TrimSpace(r)
		}
	} else {
		recipients = []string{"user1@brandonguigo.com"}
	}

	// Get CC recipients
	fmt.Print("CC (comma-separated, or press Enter to skip): ")
	ccInput, _ := reader.ReadString('\n')
	ccInput = strings.TrimSpace(ccInput)
	var ccRecipients []string
	if ccInput != "" {
		ccRecipients = strings.Split(ccInput, ",")
		for i, r := range ccRecipients {
			ccRecipients[i] = strings.TrimSpace(r)
		}
	}

	// Get BCC recipients
	fmt.Print("BCC (comma-separated, or press Enter to skip): ")
	bccInput, _ := reader.ReadString('\n')
	bccInput = strings.TrimSpace(bccInput)
	var bccRecipients []string
	if bccInput != "" {
		bccRecipients = strings.Split(bccInput, ",")
		for i, r := range bccRecipients {
			bccRecipients[i] = strings.TrimSpace(r)
		}
	}

	// Get subject
	fmt.Print("Subject: ")
	subject, _ := reader.ReadString('\n')
	subject = strings.TrimSpace(subject)

	// Get email body
	fmt.Println("Email body (type 'END' on a new line to finish):")
	var bodyLines []string
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		if line == "END" {
			break
		}
		bodyLines = append(bodyLines, line)
	}
	body := strings.Join(bodyLines, "\n")

	// Ask for attachment
	fmt.Print("Do you want to attach a file? (y/n): ")
	attachInput, _ := reader.ReadString('\n')
	attachInput = strings.TrimSpace(strings.ToLower(attachInput))

	var attachmentPath string
	if attachInput == "y" || attachInput == "yes" {
		fmt.Print("Enter file path to attach: ")
		attachmentPath, _ = reader.ReadString('\n')
		attachmentPath = strings.TrimSpace(attachmentPath)
	}

	// Get custom SMTP server if needed
	fmt.Print("SMTP Server (press Enter for default localhost:1025): ")
	smtpInput, _ := reader.ReadString('\n')
	smtpInput = strings.TrimSpace(smtpInput)
	if smtpInput != "" {
		smtpServer = smtpInput
	}

	// Get SMTP credentials if needed
	fmt.Print("SMTP Username (press Enter to skip authentication): ")
	usernameInput, _ := reader.ReadString('\n')
	usernameInput = strings.TrimSpace(usernameInput)

	var password string
	if usernameInput != "" {
		fmt.Print("SMTP Password: ")
		passwordInput, _ := reader.ReadString('\n')
		password = strings.TrimSpace(passwordInput)
	}

	// Create email message
	var message []byte
	var err error

	if attachmentPath != "" {
		message, err = createMultipartMessage(sender, recipients, ccRecipients, bccRecipients, subject, body, attachmentPath)
		if err != nil {
			fmt.Printf("Error creating multipart message: %v\n", err)
			return
		}
	} else {
		message = createSimpleMessage(sender, recipients, ccRecipients, bccRecipients, subject, body)
	}

	// Combine all recipients for sending
	allRecipients := append(recipients, ccRecipients...)
	allRecipients = append(allRecipients, bccRecipients...)

	// Send email using go-smtp
	fmt.Printf("\nConnecting to %s...\n", smtpServer)

	// Create SMTP client
	client, err := smtp.Dial(smtpServer)
	if err != nil {
		fmt.Printf("Error connecting to SMTP server: %v\n", err)
		return
	}
	defer client.Close()

	// Authenticate - use anonymous if no credentials provided
	if usernameInput != "" {
		auth := sasl.NewPlainClient("", usernameInput, password)
		if err := client.Auth(auth); err != nil {
			fmt.Printf("Error authenticating: %v\n", err)
			return
		}
	} else {
		// Use anonymous authentication
		auth := sasl.NewAnonymousClient("anonymous")
		if err := client.Auth(auth); err != nil {
			fmt.Printf("Error with anonymous authentication: %v\n", err)
			return
		}
	}

	// Set sender
	if err := client.Mail(sender, nil); err != nil {
		fmt.Printf("Error setting sender: %v\n", err)
		return
	}

	// Set recipients
	for _, recipient := range allRecipients {
		if err := client.Rcpt(recipient, nil); err != nil {
			fmt.Printf("Error setting recipient %s: %v\n", recipient, err)
			return
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		fmt.Printf("Error starting message: %v\n", err)
		return
	}

	_, err = writer.Write(message)
	if err != nil {
		fmt.Printf("Error writing message: %v\n", err)
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing message: %v\n", err)
		return
	}

	fmt.Println("Email sent successfully!")
}

func createSimpleMessage(sender string, to, cc, bcc []string, subject, body string) []byte {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("From: %s\r\n", sender))
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))

	if len(cc) > 0 {
		message.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ", ")))
	}

	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("\r\n")
	message.WriteString(body)

	return []byte(message.String())
}

func createMultipartMessage(sender string, to, cc, bcc []string, subject, body, attachmentPath string) ([]byte, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Generate boundary
	boundary := writer.Boundary()

	// Write headers
	buffer.WriteString(fmt.Sprintf("From: %s\r\n", sender))
	buffer.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))

	if len(cc) > 0 {
		buffer.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ", ")))
	}

	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	buffer.WriteString("\r\n")

	// Text part
	buffer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buffer.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buffer.WriteString("\r\n")
	buffer.WriteString(body)
	buffer.WriteString("\r\n")

	// Attachment part
	file, err := os.Open(attachmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open attachment: %v", err)
	}
	defer file.Close()

	fileName := filepath.Base(attachmentPath)
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
