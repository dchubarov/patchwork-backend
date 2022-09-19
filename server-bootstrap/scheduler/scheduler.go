package scheduler

import (
	"context"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"time"
	"twowls.org/patchwork/commons/job"
	"twowls.org/patchwork/commons/util/singleton"
	"twowls.org/patchwork/server/bootstrap/logging"
	"twowls.org/patchwork/server/bootstrap/shutdown"
)

var (
	log = logging.WithComponent("scheduler")
	goc = singleton.Lazy(
		func() *gocron.Scheduler {
			s := gocron.NewScheduler(time.UTC)
			s.WaitForScheduleAll()
			s.SingletonModeAll()
			return s
		})
)

func Start() {
	log.Info().Msg("Starting...")
	goc.Instance().StartAsync()

	shutdown.Register("scheduler", 10*time.Second, func(ctx context.Context) error {
		c, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			goc.Instance().Stop()
			cancel()
		}()

		select {
		case <-c.Done():
			if !errors.Is(c.Err(), context.Canceled) && goc.Instance().IsRunning() {
				return c.Err()
			} else {
				return nil
			}
		}
	})

	// TODO temp
	CreateJob(job.Assignment{
		JobTitle: "dummy job",
		JobWorker: func(context.Context) error {
			time.Sleep(time.Second)
			return nil
		},
		IntervalSchedule: "120s",
		StartImmediately: true,
	})
	// TODO end temp
}

func CreateJob(a job.Assignment) {
	s := goc.Instance()

	if a.JobWorker == nil {
		log.Error().Msgf("No worker defined for job %q", a)
		return
	}

	if a.CronSchedule != "" {
		log.Debug().Msgf("schedule job %q using cron expression %q", a, a.CronSchedule)
		s = s.Cron(a.CronSchedule)
	} else if a.IntervalSchedule != "" {
		log.Debug().Msgf("schedule job %q to run every %q", a, a.IntervalSchedule)
		s = s.Every(a.IntervalSchedule)
	} else {
		log.Error().Msgf("No schedule specified for job %q", a)
		return
	}

	if a.StartImmediately {
		s = s.StartImmediately()
	}

	_, err := s.DoWithJobDetails(func(job gocron.Job) {
		log.Info().Msgf("Starting job %q #%d", a, job.RunCount())
		ctx := context.TODO()

		start := time.Now()
		err := a.JobWorker(ctx)
		elapsed := time.Since(start)

		var e *zerolog.Event
		if err != nil {
			e = log.Error().Err(err)
		} else {
			e = log.Info()
		}

		e.Dur("duration", elapsed).
			Msgf("Finished job %q [%s]", a, elapsed.Round(time.Microsecond).String())
	})

	if err != nil {
		log.Error().Err(err).Msgf("CreateJob() failed")
	}
}
