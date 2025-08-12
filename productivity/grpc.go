package main

import (
	"net/http"

	globalGRPC "github.com/atomic-blend/backend/productivity/grpc/server/global"
	"github.com/atomic-blend/backend/productivity/repositories"
	"github.com/atomic-blend/backend/shared/utils/db"

	"github.com/atomic-blend/backend/grpc/gen/productivity/v1/productivityv1connect"
	"github.com/rs/zerolog/log"
)

func startGRPCServer() {
	// Initialize repositories
	taskRepo := repositories.NewTaskRepository(db.Database)
	habitRepo := repositories.NewHabitRepository(db.Database)
	noteRepo := repositories.NewNoteRepository(db.Database)
	tagRepo := repositories.NewTagRepository(db.Database)
	folderRepo := repositories.NewFolderRepository(db.Database)
	timeEntryRepo := repositories.NewTimeEntryRepository(db.Database)

	globalGRPCServer := globalGRPC.NewGrpcServer(taskRepo, habitRepo, noteRepo, tagRepo, folderRepo, timeEntryRepo)

	// TODO: register gRPC services here
	globalPath, globalHandler := productivityv1connect.NewProductivityServiceHandler(globalGRPCServer)

	// Créez un serveur HTTP mux
	mux := http.NewServeMux()
	mux.Handle(globalPath, globalHandler)

	// Démarrez le serveur HTTP
	log.Info().Msg("Starting Connect-RPC server on :50051")
	if err := http.ListenAndServe(":50051", mux); err != nil {
		log.Error().Err(err).Msg("Error serving Connect-RPC server")
	}
}
