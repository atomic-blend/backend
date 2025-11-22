package sendemail

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type RandomOptions struct {
	Domain string
	Locale string
}

func GenerateRandomConfig(opts RandomOptions) (*EmailConfig, error) {
	if opts.Locale != "" {
		gofakeit.SetGlobalFaker(gofakeit.New(0))
	}
	gofakeit.Seed(time.Now().UnixNano())

	from := gofakeit.Email()
	to := gofakeit.Email()
	subj := gofakeit.Sentence(6)
	body := strings.Join([]string{gofakeit.Paragraph(1, 3, 12, " ")}, "\n\n")

	cfg := &EmailConfig{
		Sender:     from,
		Recipients: []string{to},
		Subject:    subj,
		Body:       body,
		SMTPServer: "localhost:1025",
	}
	if opts.Domain != "" {
		// ensure message id domain uses provided domain
		cfg.MessageID = fmt.Sprintf("%s@%s", gofakeit.UUID(), opts.Domain)
	}
	return cfg, nil
}
