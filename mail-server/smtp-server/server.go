package smtpserver

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net"

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
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators.
func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if username != "username" || password != "password" {
			return errors.New("invalid_credentials")
		}
		s.auth = true
		s.user = username
		return nil
	}), nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	s.from = from
	log.Info().Msgf("Mail from: %s", from)
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	s.rcpts = append(s.rcpts, to)
	log.Info().Msgf("Rcpt to: %s", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}

	// Prepare message data for Rspamd analysis
	messageData := map[string]interface{}{
		"content":    buf.String(),
		"ip":         s.clientIP,
		"hostname":   s.hostname,
		"from":       s.from,
		"rcpt":       s.rcpts,
		"queue_id":   s.queueID,
		"user":       s.user,
		"deliver_to": s.rcpts[0], // Use first recipient as deliver_to
	}

	amqp.PublishMessage("mail", "received", messageData)

	return nil
}

func (s *Session) Reset() {
	// Reset session data for new message
	s.from = ""
	s.rcpts = nil
}

func (s *Session) Logout() error {
	return nil
}

// generateQueueID creates a unique queue ID for message tracking
func generateQueueID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
