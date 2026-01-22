package main

import (
	"net/http"

	"github.com/rs/zerolog/log"

	calendarGrpcServer "github.com/atomic-blend/backend/calendar/grpc/server"
	calendarconnect "github.com/atomic-blend/backend/grpc/gen/calendar/v1/calendarv1connect"
)

func startGRPCServer() {
	calendarGrpcServer := calendarGrpcServer.NewGrpcServer()

	globalPath, globalHandler := calendarconnect.NewCalendarServiceHandler(calendarGrpcServer)

	mux := http.NewServeMux()
	mux.Handle(globalPath, globalHandler)

	log.Info().Msg("Starting Connect-RPC server on :50051")
	if err := http.ListenAndServe(":50051", mux); err != nil {
		log.Error().Err(err).Msg("Error serving Connect-RPC server")
	}
}
