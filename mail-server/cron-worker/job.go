package cronworker

import "github.com/go-co-op/gocron/v2"

type CronJob struct {
	Name     string
	Duration gocron.JobDefinition
	Task gocron.Task
}