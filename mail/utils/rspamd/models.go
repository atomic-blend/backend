package rspamd

import (
	"net/http"
	"time"
)

// Client represents an Rspamd HTTP client for spam checking
type Client struct {
	// httpClient is the underlying HTTP client used for making requests to Rspamd.
	// Configured with timeout and other HTTP-specific settings.
	httpClient *http.Client

	// baseURL is the base URL of the Rspamd server.
	// Used to construct full URLs for API endpoints.
	baseURL string

	// password is the authentication password for the Rspamd server.
	// Sent as the "Password" header in requests when configured.
	password string
}

// CheckRequest represents the request structure for Rspamd spam checking
type CheckRequest struct {
	// DeliverTo defines the actual delivery recipient of the message.
	// Can be used for personalized statistics and user-specific options.
	DeliverTo string `json:"deliver_to,omitempty"`

	// IP defines the IP address from which this message was received.
	// Used for IP-based filtering and reputation checks.
	IP string `json:"ip,omitempty"`

	// Helo defines the SMTP HELO command value.
	// Used for HELO-based filtering and validation.
	Helo string `json:"helo,omitempty"`

	// Hostname defines the resolved hostname of the sending server.
	// Used for hostname-based filtering and reverse DNS checks.
	Hostname string `json:"hostname,omitempty"`

	// Flags defines output flags as a comma-separated list.
	// Supported flags: pass_all, groups, zstd, no_log, milter, profile, body_block, ext_urls, skip, skip_process
	Flags []string `json:"flags,omitempty"`

	// From defines the SMTP MAIL FROM command data (sender email address).
	// Used for sender-based filtering and reputation checks.
	From string `json:"from,omitempty"`

	// QueueID defines the SMTP queue ID for the message.
	// Can be used instead of message ID in logging for better traceability.
	QueueID string `json:"queue_id,omitempty"`

	// Raw indicates whether the content should be treated as raw data instead of MIME.
	// If set to true, Rspamd assumes the content is not MIME and treats it as raw data.
	Raw bool `json:"raw,omitempty"`

	// Rcpt defines the SMTP recipient(s) of the message.
	// Multiple recipients can be specified for multi-recipient messages.
	Rcpt []string `json:"rcpt,omitempty"`

	// Pass controls which filters should be checked for this message.
	// If set to "all", all filters will be checked regardless of settings.
	Pass string `json:"pass,omitempty"`

	// Subject defines the subject of the message.
	// Used for non-MIME messages or when subject needs to be explicitly specified.
	Subject string `json:"subject,omitempty"`

	// User defines the username for authenticated SMTP client.
	// Used for user-based filtering and authentication checks.
	User string `json:"user,omitempty"`

	// Message contains the actual email message content in RFC 2822 format.
	// This is the raw email data that will be analyzed by Rspamd.
	Message []byte `json:"message"`
}

// CheckResponse represents the response from Rspamd spam checking
type CheckResponse struct {
	// Action indicates the action taken by Rspamd based on the spam score.
	// Possible values: "reject", "soft reject", "add header", "rewrite subject", "greylist", "no action"
	Action string `json:"action"`

	// Score is the calculated spam score for the message.
	// Higher scores indicate higher likelihood of spam.
	Score float64 `json:"score"`

	// RequiredScore is the threshold score required to trigger the action.
	// This is the configured threshold for the action that was taken.
	RequiredScore float64 `json:"required_score"`

	// Symbols contains the triggered spam detection symbols and their details.
	// Each symbol represents a specific spam detection rule that was matched.
	Symbols map[string]interface{} `json:"symbols"`

	// Messages contains optional messages added by Rspamd filters.
	// The "smtp_message" key contains text intended to be returned as SMTP response.
	Messages map[string]string `json:"messages"`

	// Subject contains the modified subject if the action was "rewrite subject".
	// This is the new subject that should be used for the message.
	Subject string `json:"subject,omitempty"`

	// URLs contains a list of URLs found in the message (only hostnames).
	// Useful for URL-based filtering and analysis.
	URLs []string `json:"urls,omitempty"`

	// Emails contains a list of email addresses found in the message.
	// Useful for email address extraction and analysis.
	Emails []string `json:"emails,omitempty"`

	// MessageID is the ID of the message, useful for logging and tracking.
	// This helps correlate Rspamd results with message processing.
	MessageID string `json:"message-id,omitempty"`

	// Milter contains milter-specific response data for header manipulation.
	// Used when integrating with MTA milter interfaces.
	Milter *MilterResponse `json:"milter,omitempty"`
}

// MilterResponse represents milter-specific response data for header manipulation
type MilterResponse struct {
	// AddHeaders specifies headers to add to the message.
	// The key is the header name, and the value contains the header value and insertion order.
	AddHeaders map[string]HeaderValue `json:"add_headers,omitempty"`

	// RemoveHeaders specifies headers to remove from the message.
	// The key is the header name, and the value is the order of the header to remove.
	// Special values: 0 = remove all headers with this name, negative = remove Nth header from end.
	RemoveHeaders map[string]int `json:"remove_headers,omitempty"`

	// ChangeFrom specifies a new SMTP MAIL FROM value to replace the original.
	// Used to modify the sender address during message processing.
	ChangeFrom string `json:"change_from,omitempty"`

	// Reject specifies a custom rejection message.
	// Values like "discard" or "quarantine" can be used for custom rejection handling.
	Reject string `json:"reject,omitempty"`

	// SpamHeader specifies a custom spam header name to add.
	// The header will contain spam detection information.
	SpamHeader string `json:"spam_header,omitempty"`

	// NoAction indicates that instead of performing any action, just add an X-Rspamd-Action header.
	// The message will be accepted with the action information in the header.
	NoAction bool `json:"no_action,omitempty"`

	// AddRcpt specifies new recipients to add to the message.
	// Used to add additional recipients during message processing.
	AddRcpt []string `json:"add_rcpt,omitempty"`

	// DelRcpt specifies recipients to remove from the message.
	// Used to remove recipients during message processing.
	DelRcpt []string `json:"del_rcpt,omitempty"`
}

// HeaderValue represents a header value with insertion order
type HeaderValue struct {
	// Value is the actual header value content.
	// This is the text that will be set as the header value.
	Value string `json:"value"`

	// Order specifies the insertion order of the header.
	// Lower numbers are inserted first (e.g., 0 = first header).
	Order int `json:"order"`
}

// Config holds the configuration for the Rspamd client
type Config struct {
	// BaseURL is the base URL of the Rspamd server.
	// Should include protocol, hostname, and port (e.g., "http://localhost:11333").
	BaseURL string

	// Password is the authentication password for the Rspamd server.
	// Required if Rspamd is configured with password authentication.
	Password string

	// Timeout is the HTTP timeout for requests to the Rspamd server.
	// Should be set appropriately for your network conditions and message sizes.
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts for failed requests.
	// Used to handle temporary network issues or server unavailability.
	MaxRetries int
}
