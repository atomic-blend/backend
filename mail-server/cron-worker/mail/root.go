package mail

import (
	"time"

	cronworker "github.com/atomic-blend/backend/mail-server/cron-worker"
	"github.com/go-co-op/gocron/v2"
)

func GetMailCronJobs() []cronworker.CronJob {
	jobList := []cronworker.CronJob{}

	// send mail job
	jobList = append(jobList, cronworker.CronJob{
		Name:     "SendMail",
		Duration: gocron.DurationJob(500 * time.Millisecond),
		Task:     gocron.NewTask(func() { sendOrRetryPendingSentEmails() }),
	})

	return jobList
}
