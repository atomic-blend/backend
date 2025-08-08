package main

import (
	"github.com/atomic-blend/backend/mail-server/cron-worker/mail"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

func startCronWorker() {
	log.Info().Msg("Starting cron worker...")
	// Initialize the cron scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing cron scheduler")
		return
	}

	// Get the mail cron jobs
	mailCronJobs := mail.GetMailCronJobs()

	// Add each job to the scheduler
	for _, job := range mailCronJobs {
		log.Info().Msgf("Adding job: %s with duration: %v", job.Name, job.Duration)
		_, err := scheduler.NewJob(job.Duration, job.Task)
		if err != nil {
			log.Error().Err(err).Msg("Error adding job to scheduler")
		}
	}

	// Start the scheduler
	scheduler.Start()

	select {} // wait forever
}
