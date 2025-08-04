package rspamd

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Example usage of the Rspamd client
func Example() {
	// Create client with default configuration (uses environment variables)
	client := NewClient(nil)

	// Example email message in RFC 2822 format
	emailMessage := []byte(`From: sender@example.com
To: recipient@example.com
Subject: Test Email
Date: Mon, 01 Jan 2024 12:00:00 +0000
Message-ID: <test@example.com>

This is a test email message.
`)

	// Create check request
	req := &CheckRequest{
		From:     "sender@example.com",
		Rcpt:     []string{"recipient@example.com"},
		IP:       "192.168.1.1",
		Helo:     "mail.example.com",
		Hostname: "mail.example.com",
		Message:  emailMessage,
	}

	// Check message for spam
	resp, err := client.CheckMessage(req)
	if err != nil {
		log.Fatalf("Failed to check message: %v", err)
	}

	// Process results
	fmt.Printf("Spam Score: %.2f\n", resp.GetScore())
	fmt.Printf("Action: %s\n", resp.GetAction())
	fmt.Printf("Is Spam: %t\n", resp.IsSpam())

	// Check if specific symbols were triggered
	if len(resp.Symbols) > 0 {
		fmt.Println("Triggered symbols:")
		for symbol, details := range resp.Symbols {
			fmt.Printf("  - %s: %v\n", symbol, details)
		}
	}

	// Check for URLs found in the message
	if len(resp.URLs) > 0 {
		fmt.Println("URLs found:")
		for _, url := range resp.URLs {
			fmt.Printf("  - %s\n", url)
		}
	}

	// Check for email addresses found in the message
	if len(resp.Emails) > 0 {
		fmt.Println("Email addresses found:")
		for _, email := range resp.Emails {
			fmt.Printf("  - %s\n", email)
		}
	}
}

// ExampleWithCustomConfig demonstrates using custom configuration
func ExampleWithCustomConfig() {
	// Create custom configuration
	config := &Config{
		BaseURL:    "http://rspamd.internal:11333",
		Password:   "mysecretpassword",
		Timeout:    60 * time.Second,
		MaxRetries: 5,
	}

	// Create client with custom configuration
	client := NewClient(config)

	// Test connection
	if err := client.Ping(); err != nil {
		log.Fatalf("Failed to ping Rspamd: %v", err)
	}

	fmt.Println("Successfully connected to Rspamd")
}

// ExampleWithEnvironmentVariables demonstrates using environment variables
func ExampleWithEnvironmentVariables() {
	// Set environment variables (in real usage, these would be set in your environment)
	os.Setenv("RSPAMD_BASE_URL", "http://rspamd.production:11333")
	os.Setenv("RSPAMD_PASSWORD", "prodpassword")
	os.Setenv("RSPAMD_TIMEOUT_SECONDS", "45")
	os.Setenv("RSPAMD_MAX_RETRIES", "3")

	// Create client with environment-based configuration
	client := NewClient(nil)

	// Example message
	message := []byte(`From: spammer@evil.com
To: user@company.com
Subject: Make money fast!
Date: Mon, 01 Jan 2024 12:00:00 +0000

CLICK HERE TO MAKE MILLIONS FAST!!!
http://evil.com/scam
`)

	req := &CheckRequest{
		From:     "spammer@evil.com",
		Rcpt:     []string{"user@company.com"},
		IP:       "1.2.3.4",
		Helo:     "evil.com",
		Hostname: "evil.com",
		Message:  message,
	}

	resp, err := client.CheckMessage(req)
	if err != nil {
		log.Fatalf("Failed to check message: %v", err)
	}

	fmt.Printf("Spam Score: %.2f\n", resp.GetScore())
	fmt.Printf("Action: %s\n", resp.GetAction())
	fmt.Printf("Is Spam: %t\n", resp.IsSpam())
}
