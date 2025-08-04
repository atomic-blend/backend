# Interactive SMTP Email Sender

This script provides an interactive command-line interface to send emails to your local email server for development and testing purposes.

**⚠️ This script is intended for local development only and should not be used in production environments.**

## Features

- **Interactive prompts** for all email fields
- **Multiple recipients** support (To, CC, BCC)
- **File attachments** with proper MIME handling
- **Custom SMTP server** configuration
- **Default values** for quick testing
- **Multi-line email body** input
- **Base64 encoding** for attachments (RFC 2045 compliant)

## Prerequisites

- Go 1.19 or higher
- Local email server running (default: `localhost:1025`)
- Docker Compose setup with email services

## Usage

### Building the Script

```bash
cd mail-server/scripts/send-email-via-smtp
go build -o send-email main.go
```

### Running the Script

```bash
./send-email
```

Or run directly with Go:

```bash
go run main.go
```

## Interactive Prompts

When you run the script, it will prompt you for the following information:

### 1. Sender Email
```
From (sender email) [brandon@brandonguigo.com]: 
```
- **Default**: `brandon@brandonguigo.com`
- **Override**: Type a different email address
- **Skip**: Press Enter to use default

### 2. Recipients
```
To (recipient email(s), comma-separated for multiple) [user1@brandonguigo.com]: 
```
- **Default**: `user1@brandonguigo.com`
- **Multiple**: `user1@example.com, user2@example.com`
- **Skip**: Press Enter to use default

### 3. CC Recipients (Optional)
```
CC (comma-separated, or press Enter to skip): 
```
- **Skip**: Press Enter to leave empty
- **Multiple**: `cc1@example.com, cc2@example.com`

### 4. BCC Recipients (Optional)
```
BCC (comma-separated, or press Enter to skip): 
```
- **Skip**: Press Enter to leave empty
- **Multiple**: `bcc1@example.com, bcc2@example.com`

### 5. Subject Line
```
Subject: 
```
- **Required**: Must provide a subject

### 6. Email Body
```
Email body (type 'END' on a new line to finish):
Hello,
This is a test email.
END
```
- **Multi-line**: Type your message line by line
- **Finish**: Type `END` on a new line to complete

### 7. Attachments (Optional)
```
Do you want to attach a file? (y/n): y
Enter file path to attach: /path/to/document.pdf
```
- **Yes/No**: Answer `y`, `yes`, `n`, or `no`
- **File path**: Provide absolute or relative path to file
- **Supported**: All file types with automatic MIME detection

### 8. SMTP Server (Optional)
```
SMTP Server (press Enter for default localhost:1025): 
```
- **Default**: `localhost:1025`
- **Custom**: Type different server address and port

## Example Session

```
=== Interactive SMTP Email Sender ===

From (sender email) [brandon@brandonguigo.com]: 
To (recipient email(s), comma-separated for multiple) [user1@brandonguigo.com]: 
CC (comma-separated, or press Enter to skip): 
BCC (comma-separated, or press Enter to skip): 
Subject: Test Interactive Email
Email body (type 'END' on a new line to finish):
Hello,

This is a test email sent from the interactive SMTP script.
It supports multiple lines and attachments.

Best regards,
Brandon
END
Do you want to attach a file? (y/n): y
Enter file path to attach: /tmp/test.txt
SMTP Server (press Enter for default localhost:1025): 

Connecting to localhost:1025...
Email sent successfully!
```

## Technical Details

### Email Format Support
- **Simple emails**: Text-only messages
- **Multipart emails**: Messages with attachments
- **MIME compliance**: RFC 2045 and RFC 2822 compliant
- **Character encoding**: UTF-8 support

### Attachment Handling
- **Base64 encoding**: Automatic encoding for binary files
- **MIME type detection**: Based on file extension
- **Filename preservation**: Original filename maintained
- **Size handling**: No artificial size limits

### SMTP Configuration
- **Authentication**: Plain text authentication
- **Default credentials**: `username`/`password` (for local testing)
- **Connection**: TCP connection to specified server
- **Error handling**: Detailed error messages

## Troubleshooting

### Common Issues

1. **Connection refused**
   - Ensure your email server is running
   - Check if the port is correct (default: 1025)
   - Verify Docker Compose services are up

2. **File not found**
   - Use absolute paths for attachments
   - Ensure file permissions allow reading
   - Check file exists before running script

3. **SMTP authentication errors**
   - Default credentials are for local testing only
   - Update authentication for production servers

### Debug Mode

To see more detailed information, you can modify the script to add debug logging:

```go
// Add before smtp.SendMail call
fmt.Printf("Sending to: %v\n", allRecipients)
fmt.Printf("Message size: %d bytes\n", len(message))
```

## Development Notes

- **Local testing only**: This script is designed for development environments
- **No encryption**: Uses plain text SMTP (suitable for local testing)
- **No validation**: Minimal input validation for development speed
- **Extensible**: Easy to add more features like HTML support, multiple attachments, etc.

## Contributing

When modifying this script:
1. Maintain backward compatibility with existing Docker setup
2. Test with various email formats
3. Update this README for any new features
4. Keep it simple for development use cases