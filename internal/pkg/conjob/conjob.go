package conjob

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type conjobRun interface {
	Run() error
}

func NewJob(scheduler gocron.Scheduler, name string, timeDuration time.Duration, newFunc conjobRun) {
	err := newFunc.Run()
	if err != nil {
		return
	}
	scheduler.NewJob(
		gocron.DurationJob(
			timeDuration,
		),
		gocron.NewTask(
			func() {
				newFunc.Run()
			},
		),
		gocron.WithName(name),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(
				func(jobID uuid.UUID, jobName string) {
					log.Printf("Job starting: %s, %s \n", jobID.String(), jobName)
				},
			),
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					log.Printf("Job completed: %s, %s \n", jobID.String(), jobName)
				},
			),
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					log.Printf("Job had an error: %s, %s %v\n", jobID.String(), jobName, err)
					log.Fatal(err)
				},
			),
		),
	)

}
