package smtpserver

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net"
	"time"

	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

// The Backend implements SMTP server methods.
type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	// Get client IP address
	var clientIP string
	if addr, ok := c.Conn().RemoteAddr().(*net.TCPAddr); ok {
		clientIP = addr.IP.String()
	}

	// Get HELO / hostname from connection
	hostname := c.Hostname()

	// Generate a unique queue ID
	queueID := generateQueueID()

	return &Session{
		clientIP: clientIP,
		hostname: hostname,
		queueID:  queueID,
	}, nil
}

// A Session is returned after successful login.
type Session struct {
	auth     bool
	clientIP string
	// HELO / hostname
	hostname string
	queueID  string
	from     string
	rcpts    []string
	user     string
}

// AuthMechanisms returns a slice of available auth mechanisms; only PLAIN is
// supported in this example.
func (s *Session) AuthMechanisms() []string {
	return []string{}
}

// Auth is the handler for supported authenticators.
func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewAnonymousServer(func(trace string) error {
		s.auth = true
		s.user = "anonymous"
		return nil
	}), nil
}

// Mail is the handler for the MAIL command.
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	// if !s.auth {
	// 	return smtp.ErrAuthRequired
	// }
	s.from = from
	log.Info().Msgf("Mail from: %s", from)
	return nil
}

// Rcpt is the handler for the RCPT command.
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// if !s.auth {
	// 	return smtp.ErrAuthRequired
	// }
	s.rcpts = append(s.rcpts, to)
	log.Info().Msgf("Rcpt to: %s", to)
	return nil
}

// Data is the handler for the DATA command.
func (s *Session) Data(r io.Reader) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}

	// Prepare message data for Rspamd analysis
	messageData := map[string]interface{}{
		"content":     buf.String(),
		"ip":          s.clientIP,
		"hostname":    s.hostname,
		"from":        s.from,
		"rcpt":        s.rcpts,
		"queue_id":    s.queueID,
		"user":        s.user,
		"deliver_to":  s.rcpts[0], // Use first recipient as deliver_to
		"received_at": time.Now().Format(time.RFC3339),
	}

	amqp.PublishMessage("mail", "received", messageData)

	return nil
}

// Reset is the handler for the RESET command.
func (s *Session) Reset() {
	// Reset session data for new message
	s.from = ""
	s.rcpts = nil
}

// Logout is the handler for the LOGOUT command.
func (s *Session) Logout() error {
	return nil
}

// generateQueueID creates a unique queue ID for message tracking
func generateQueueID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
