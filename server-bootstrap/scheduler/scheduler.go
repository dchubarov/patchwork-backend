package scheduler

import (
	"context"
	"errors"
	"github.com/go-co-op/gocron"
	"time"
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
}
