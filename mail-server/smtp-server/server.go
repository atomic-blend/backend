package smtpserver

import (
	"bytes"
	"errors"
	"io"

	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	"github.com/emersion/go-message/mail"
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

    // Save a copy of full raw message for MongoDB
    fullMessage := buf.Bytes()

    // Parse message with go-message
    mr, err := mail.CreateReader(bytes.NewReader(fullMessage))
    if err != nil {
        return err
    }

    // Extract headers (optional)
    header := mr.Header
    subject, _ := header.Subject()
    from, _ := header.AddressList("From")
    to, _ := header.AddressList("To")

	// attachements is list of objects representing the attachments
	attachments := []interface{}{}


    // Iterate over parts (text + attachments)
    for {
        p, err := mr.NextPart()
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        switch h := p.Header.(type) {
        case *mail.InlineHeader:
            // Body part
            body, _ := io.ReadAll(p.Body)
            log.Printf("Body: %s", string(body))

        case *mail.AttachmentHeader:
            filename, _ := h.Filename()
            contentType, _, _ := h.ContentType()

            // Read attachment content
            content, _ := io.ReadAll(p.Body)

			attachments = append(attachments, map[string]interface{}{
				"filename": filename,
				"contentType": contentType,
				"content": content,
			})

            log.Printf("Uploaded attachment: %s", filename)
        }
    }

	log.Debug().Msgf("Attachments: %s", attachments)

	amqp.PublishMessage("mail", "received", map[string]interface{}{
		"from": from,
		"to": to,
		"subject": subject,
		"attachments": attachments,
	})

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
