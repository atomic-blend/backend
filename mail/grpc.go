package main

import (
	"net/http"

	"github.com/rs/zerolog/log"

	mailconnect "github.com/atomic-blend/backend/grpc/gen/mail/v1/mailv1connect"
	mailGrpcServer "github.com/atomic-blend/backend/mail/grpc/server/global"
)

func startGRPCServer() {
	mailGrpcServer := mailGrpcServer.NewGrpcServer()

	globalPath, globalHandler := mailconnect.NewMailServiceHandler(mailGrpcServer)

	mux := http.NewServeMux()
	mux.Handle(globalPath, globalHandler)

	// Démarrez le serveur HTTP
	log.Info().Msg("Starting Connect-RPC server on :50051")
	if err := http.ListenAndServe(":50051", mux); err != nil {
		log.Error().Err(err).Msg("Error serving Connect-RPC server")
	}
}
