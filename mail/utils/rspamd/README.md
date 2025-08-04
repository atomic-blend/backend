# Rspamd HTTP Client

A Go HTTP client for interacting with the Rspamd spam filtering service. This client supports the `/checkv2` endpoint for spam checking and is configurable via environment variables.

## Features

- HTTP client for Rspamd spam checking
- Configurable via environment variables
- Support for all Rspamd HTTP headers
- Comprehensive response parsing
- Connection health checking via ping
- Helper methods for spam detection

## Installation

The client is part of the mail utilities package and can be imported as:

```go
import "your-project/mail/utils/rspamd"
```

## Configuration

The client can be configured using environment variables or programmatically.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RSPAMD_BASE_URL` | Rspamd server URL | `http://localhost:11333` |
| `RSPAMD_PASSWORD` | Password for authentication | (empty) |
| `RSPAMD_TIMEOUT_SECONDS` | HTTP timeout in seconds | `30` |
| `RSPAMD_MAX_RETRIES` | Maximum retry attempts | `3` |

### Programmatic Configuration

```go
config := &rspamd.Config{
    BaseURL:    "http://rspamd.internal:11333",
    Password:   "mysecretpassword",
    Timeout:    60 * time.Second,
    MaxRetries: 5,
}
client := rspamd.NewClient(config)
```

## Usage

### Basic Usage

```go
// Create client with default configuration
client := rspamd.NewClient(nil)

// Create check request
req := &rspamd.CheckRequest{
    From:    "sender@example.com",
    Rcpt:    []string{"recipient@example.com"},
    IP:      "192.168.1.1",
    Helo:    "mail.example.com",
    Hostname: "mail.example.com",
    Message: emailMessage, // []byte containing RFC 2822 email
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
```

### Advanced Usage with All Headers

```go
req := &rspamd.CheckRequest{
    DeliverTo: "user@company.com",
    IP:        "1.2.3.4",
    Helo:      "mail.example.com",
    Hostname:  "mail.example.com",
    Flags:     []string{"pass_all", "groups"},
    From:      "sender@example.com",
    QueueID:   "queue123",
    Raw:       false,
    Rcpt:      []string{"user@company.com", "admin@company.com"},
    Pass:      "all",
    Subject:   "Important Message",
    User:      "authenticated_user",
    Message:   emailMessage,
}
```

### Health Checking

```go
// Test connection to Rspamd
if err := client.Ping(); err != nil {
    log.Printf("Rspamd is not available: %v", err)
} else {
    log.Println("Rspamd is available")
}
```

## Response Structure

The `CheckResponse` contains the following fields:

- `Action`: The action taken by Rspamd (e.g., "reject", "soft reject", "add header", "no action")
- `Score`: The spam score
- `RequiredScore`: The threshold score for actions
- `Symbols`: Map of triggered symbols and their details
- `Messages`: Optional messages from Rspamd filters
- `Subject`: Modified subject (if action is "rewrite subject")
- `URLs`: List of URLs found in the message
- `Emails`: List of email addresses found in the message
- `MessageID`: Message ID for logging
- `Milter`: Milter-specific response data

### Helper Methods

```go
// Check if message is classified as spam
if resp.IsSpam() {
    // Handle spam
}

// Get spam score
score := resp.GetScore()

// Get action taken
action := resp.GetAction()
```

## Supported Rspamd Headers

The client supports all standard Rspamd HTTP headers:

- `Deliver-To`: Actual delivery recipient
- `IP`: Source IP address
- `Helo`: SMTP HELO command
- `Hostname`: Resolved hostname
- `Flags`: Output flags (comma-separated)
- `From`: SMTP MAIL FROM command
- `Queue-Id`: SMTP queue ID
- `Raw`: Treat content as raw data
- `Rcpt`: SMTP recipient(s)
- `Pass`: Pass all filters
- `Subject`: Message subject
- `User`: Authenticated SMTP client username
- `Password`: Authentication password

## Error Handling

The client returns descriptive errors for various failure scenarios:

- Network connectivity issues
- Invalid HTTP responses
- JSON parsing errors
- Configuration errors

## Testing

Run the tests with:

```bash
go test ./mail/utils/rspamd
```

## Examples

See `example.go` for complete usage examples including:

- Basic spam checking
- Custom configuration
- Environment variable usage
- Response processing

## References

- [Rspamd Protocol Documentation](https://docs.rspamd.com/developers/protocol)
- [Rspamd HTTP API](https://docs.rspamd.com/developers/protocol#rspamd-http-request) 