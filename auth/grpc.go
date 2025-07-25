package main

import (
	"net/http"

	userGrpc "github.com/atomic-blend/backend/auth/grpc/server/user"
	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/auth/utils/db"
	authconnect "github.com/atomic-blend/backend/grpc/gen/auth/authconnect"

	"github.com/rs/zerolog/log"
)

func startGRPCServer() {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.Database)

	UserGrpcServer := userGrpc.NewUserGrpcServer(userRepo)

	// TODO: register gRPC services here
	globalPath, globalHandler := authconnect.NewUserServiceHandler(UserGrpcServer)

	// Créez un serveur HTTP mux
	mux := http.NewServeMux()
	mux.Handle(globalPath, globalHandler)

	// Démarrez le serveur HTTP
	log.Info().Msg("Starting Connect-RPC server on :50051")
	if err := http.ListenAndServe(":50051", mux); err != nil {
		log.Error().Err(err).Msg("Error serving Connect-RPC server")
	}
}
