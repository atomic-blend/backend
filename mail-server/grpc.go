package main

import (
	"net/http"

	"github.com/rs/zerolog/log"

	mailconnect "github.com/atomic-blend/backend/grpc/gen/mailserver/v1/mailserverv1connect"
	mailserverGrpcServer "github.com/atomic-blend/backend/mail-server/grpc/server"
)

func startGRPCServer() {
	mailGrpcServer := mailserverGrpcServer.NewGrpcServer()

	globalPath, globalHandler := mailconnect.NewMailServerServiceHandler(mailGrpcServer)

	mux := http.NewServeMux()
	mux.Handle(globalPath, globalHandler)

	// Démarrez le serveur HTTP
	log.Info().Msg("Starting Connect-RPC server on :50051")
	if err := http.ListenAndServe(":50051", mux); err != nil {
		log.Error().Err(err).Msg("Error serving Connect-RPC server")
	}
}
