package scheduler

import (
	"github.com/go-co-op/gocron"
	"time"
)

type (
	covid struct {
		cron *gocron.Scheduler
		stop chan bool
	}

	IScheduler interface {
		ExecuteJob(jobs ...CovidJob) error
		Stop()
	}

	CovidJob struct {
		Cron       string
		Task       func() error
		TaskParams []interface{}
	}
)

const (
	defaultLocation = "Asia/Tokyo"
)

var (
	jstLocation *time.Location
)

func init() {
	jstLocation, _ = time.LoadLocation(defaultLocation)
}

func New() IScheduler {
	return &covid{
		cron: gocron.NewScheduler(jstLocation),
	}
}

func (sc *covid) ExecuteJob(jobs ...CovidJob) error {
	for _, job := range jobs {
		if _, err := sc.cron.Cron(job.Cron).Do(job.Task, job.TaskParams...); err != nil {
			return err
		}
		sc.cron.StartAsync()
	}
	return nil
}

func (sc *covid) Stop() {
	sc.cron.Stop()
}
