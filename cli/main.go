/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	rootcmd "github.com/atomic-blend/backend/cli/root_cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	debug := os.Getenv("LOG_LEVEL")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	rootcmd.Execute()
}
