package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func startGRPCServer() {
	// Créez un serveur HTTP mux
	mux := http.NewServeMux()

	// Démarrez le serveur HTTP
	log.Info().Msg("Starting Connect-RPC server on :50051")
	if err := http.ListenAndServe(":50051", mux); err != nil {
		log.Error().Err(err).Msg("Error serving Connect-RPC server")
	}
}
