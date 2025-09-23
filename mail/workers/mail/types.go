package mail

// ReceivedMailPayload represents the complete AMQP payload for received mail messages
type ReceivedMailPayload struct {
	Content    string   `json:"content"`     // MIME content of the email
	IP         string   `json:"ip"`          // Client IP address
	Hostname   string   `json:"hostname"`    // Server hostname
	From       string   `json:"from"`        // Sender email address
	Rcpt       []string `json:"rcpt"`        // Recipient email addresses
	QueueID    string   `json:"queue_id"`    // Queue identifier
	User       string   `json:"user"`        // Authenticated user (if any)
	DeliverTo  string   `json:"deliver_to"`  // Primary delivery address
	ReceivedAt string   `json:"received_at"` // Date and time when the email was received
}
