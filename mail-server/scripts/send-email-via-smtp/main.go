package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	// SMTP configuration
	smtpServer := "localhost:1025"
	sender := "test@example.com"
	recipient := "recipient@example.com"

	// Email content
	subject := "Test Email from Local SMTP Server"
	body := `This is a test email sent to verify that the local SMTP server is working correctly.

If you receive this email, your SMTP server is functioning properly!

Best regards,
Your Local SMTP Server`

	// Create email message
	message := fmt.Sprintf("From: %s\r\n", sender) +
		fmt.Sprintf("To: %s\r\n", recipient) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		body

	// Send email
	fmt.Printf("Connecting to %s...\n", smtpServer)
	auth := smtp.PlainAuth("", "username", "password", "localhost")
	err := smtp.SendMail(smtpServer, auth, sender, []string{recipient}, []byte(message))

	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
