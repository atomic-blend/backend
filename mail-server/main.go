package main

import (
	"os"
	"time"

	smtpserver "github.com/atomic-blend/backend/mail-server/smtp-server"
	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	"github.com/emersion/go-smtp"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	_ = godotenv.Load()

	// Define port
	port := os.Getenv("MAIL_PORT")
	if port == "" {
		port = "1025"
	}

	host := os.Getenv("MAIL_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	amqp.InitConsumerAMQP()
	amqp.InitProducerAMQP()

	// launch the AMQP consumer in a goroutine
	go processMessages()

	// start the grpc server
	go startGRPCServer()

	// instanciate the smtp backend
	be := &smtpserver.Backend{}

	// create the smtp server
	s := smtp.NewServer(be)

	s.Addr = host + ":" + port
	s.Domain = host
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Info().Msgf("Starting server at %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
	}
}
