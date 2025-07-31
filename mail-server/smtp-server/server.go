package smtpserver

import (
	"bytes"
	"errors"
	"io"

	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

// The Backend implements SMTP server methods.
type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after successful login.
type Session struct {
	auth bool
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
		return nil
	}), nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	log.Info().Msgf("Mail from: %s", from)
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	log.Info().Msgf("Rcpt to: %s", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}

	amqp.PublishMessage("mail", "received", map[string]interface{}{
		"content": buf.String(),
	})

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
