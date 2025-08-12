// Package rspamdclient contains the Rspamd client models
package rspamdclient

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
	// This is the new subject that should be used when sending the email.
	Subject string `json:"subject,omitempty"`

	// MessageID is the unique identifier for the message being processed.
	// This helps correlate Rspamd results with message processing.
	MessageID string `json:"message_id,omitempty"`

	// TimeReal is the real processing time for the message.
	// Measured in seconds, indicates how long Rspamd took to process the message.
	TimeReal float64 `json:"time_real,omitempty"`

	// TimeVirtual is the virtual processing time for the message.
	// This is the CPU time used for processing, not wall clock time.
	TimeVirtual float64 `json:"time_virtual,omitempty"`

	// MilterConfig is the milter configuration used for this message.
	// Contains milter-specific settings and configuration.
	MilterConfig map[string]interface{} `json:"milter_config,omitempty"`

	// URLReputation contains URL reputation information if available.
	// Provides reputation scores for URLs found in the message.
	URLReputation map[string]interface{} `json:"url_reputation,omitempty"`

	// Emails contains email reputation information if available.
	// Provides reputation scores for email addresses found in the message.
	Emails map[string]interface{} `json:"emails,omitempty"`

	// IPReputation contains IP reputation information if available.
	// Provides reputation scores for IP addresses found in the message.
	IPReputation map[string]interface{} `json:"ip_reputation,omitempty"`

	// ASN contains Autonomous System Number information if available.
	// Provides information about the network provider of the sending IP.
	ASN map[string]interface{} `json:"asn,omitempty"`

	// Country contains country information if available.
	// Provides the country of origin for the sending IP address.
	Country map[string]interface{} `json:"country,omitempty"`

	// City contains city information if available.
	// Provides the city of origin for the sending IP address.
	City map[string]interface{} `json:"city,omitempty"`

	// Coordinates contains geographical coordinates if available.
	// Provides latitude and longitude for the sending IP address.
	Coordinates map[string]interface{} `json:"coordinates,omitempty"`

	// Network contains network information if available.
	// Provides network details for the sending IP address.
	Network map[string]interface{} `json:"network,omitempty"`

	// TLS contains TLS information if available.
	// Provides details about TLS encryption used in the connection.
	TLS map[string]interface{} `json:"tls,omitempty"`

	// Authenticated contains authentication information if available.
	// Provides details about SMTP authentication if used.
	Authenticated map[string]interface{} `json:"authenticated,omitempty"`

	// DKIM contains DKIM signature information if available.
	// Provides details about DKIM signatures found in the message.
	DKIM map[string]interface{} `json:"dkim,omitempty"`

	// SPF contains SPF record information if available.
	// Provides details about SPF records for the sending domain.
	SPF map[string]interface{} `json:"spf,omitempty"`

	// DMARC contains DMARC policy information if available.
	// Provides details about DMARC policies for the sending domain.
	DMARC map[string]interface{} `json:"dmarc,omitempty"`

	// ARC contains ARC signature information if available.
	// Provides details about ARC signatures found in the message.
	ARC map[string]interface{} `json:"arc,omitempty"`

	// RspamdVersion contains the version of Rspamd that processed the message.
	// Useful for debugging and compatibility checking.
	RspamdVersion string `json:"rspamd_version,omitempty"`

	// ScanTime contains the timestamp when the message was scanned.
	// Useful for logging and debugging purposes.
	ScanTime string `json:"scan_time,omitempty"`

	// QueueID contains the queue ID if provided in the request.
	// This helps correlate Rspamd results with message processing.
	QueueID string `json:"queue_id,omitempty"`

	// NoAction indicates that instead of performing any action, just add an X-Rspamd-Action header.
	// This is useful when you want to log the action but not enforce it.
	NoAction bool `json:"no_action,omitempty"`

	// Skip indicates that the message should be skipped by certain filters.
	// This can be used to bypass specific checks for testing or special cases.
	Skip bool `json:"skip,omitempty"`

	// SkipProcess indicates that the message should not be processed by Rspamd.
	// This is useful for bypassing all Rspamd processing.
	SkipProcess bool `json:"skip_process,omitempty"`

	// Profile indicates which profile should be used for processing.
	// Different profiles can have different rules and thresholds.
	Profile string `json:"profile,omitempty"`

	// Groups indicates which groups of rules should be applied.
	// This can be used to enable or disable specific rule groups.
	Groups []string `json:"groups,omitempty"`

	// Zstd indicates whether to use zstd compression for the message.
	// This can reduce bandwidth usage for large messages.
	Zstd bool `json:"zstd,omitempty"`

	// NoLog indicates whether to suppress logging for this message.
	// Useful for reducing log noise during testing or debugging.
	NoLog bool `json:"no_log,omitempty"`

	// Milter indicates whether to use milter mode for processing.
	// Milter mode integrates with mail servers for real-time processing.
	Milter bool `json:"milter,omitempty"`

	// BodyBlock indicates whether to block the message body if it contains certain content.
	// This can be used to block messages with specific content patterns.
	BodyBlock bool `json:"body_block,omitempty"`

	// ExtUrls indicates whether to extract and analyze URLs in the message.
	// This can help detect phishing and malicious links.
	ExtUrls bool `json:"ext_urls,omitempty"`
}

// Config holds the configuration for the Rspamd client
type Config struct {
	// BaseURL is the base URL of the Rspamd server.
	// Should include protocol (http/https) and port if not default.
	BaseURL string

	// Password is the authentication password for the Rspamd server.
	// Required if Rspamd is configured with password authentication.
	Password string

	// Timeout is the HTTP timeout for requests to the Rspamd server.
	// Should be set to a reasonable value to avoid hanging requests.
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts for failed requests.
	// Useful for handling temporary network issues or server overload.
	MaxRetries int
}
