package main

import (
	"net/http"

	globalGRPC "github.com/atomic-blend/backend/productivity/grpc/global"
	"github.com/atomic-blend/backend/productivity/repositories"
	"github.com/atomic-blend/backend/productivity/utils/db"

	"github.com/atomic-blend/backend/grpc/gen/productivity/productivityconnect"
	"github.com/rs/zerolog/log"
)

func startGRPCServer() {
	// Initialize repositories
	taskRepo := repositories.NewTaskRepository(db.Database)
	habitRepo := repositories.NewHabitRepository(db.Database)

	globalGRPCServer := globalGRPC.NewGrpcServer(taskRepo, habitRepo)

	// TODO: register gRPC services here
	globalPath, globalHandler := productivityconnect.NewProductivityServiceHandler(globalGRPCServer)

	// Créez un serveur HTTP mux
	mux := http.NewServeMux()
	mux.Handle(globalPath, globalHandler)

	// Démarrez le serveur HTTP
	log.Info().Msg("Starting Connect-RPC server on :50051")
	if err := http.ListenAndServe(":50051", mux); err != nil {
		log.Error().Err(err).Msg("Error serving Connect-RPC server")
	}
}
