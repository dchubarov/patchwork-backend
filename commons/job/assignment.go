package job

import (
	"context"
	"time"
)

type WorkerFunc func(context.Context) error

// Assignment represents a job assignment
type Assignment struct {
	JobTitle         string
	JobWorker        WorkerFunc
	IntervalSchedule string
	CronSchedule     string
	StartImmediately bool
	Timeout          time.Duration
}

func (a Assignment) String() string {
	return a.JobTitle
}
