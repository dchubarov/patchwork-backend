package shutdown

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"time"
	"twowls.org/patchwork/commons/utils/singleton"
	"twowls.org/patchwork/server/bootstrap/logging"
)

type hookFunc func(ctx context.Context) error

type hook struct {
	tag     string
	handler hookFunc
	timeout time.Duration
	elapsed time.Duration
	error   error
	done    bool
}

type Shutdown struct {
	hooks      []*hook
	InProgress bool
}

var (
	s = singleton.NewLazy(func() *Shutdown {
		return &Shutdown{}
	})
)

func (s *Shutdown) Register(tag string, timeout time.Duration, handler hookFunc) {
	if s.InProgress {
		logging.Panic("cannot register while shutdown is in progress")
	}

	logging.Debug("Registered shutdown hook for %s", tag)
	s.hooks = slices.Insert(s.hooks, 0, &hook{
		tag:     tag,
		handler: handler,
		timeout: timeout,
	})
}

func (s *Shutdown) ShutdownAll() {
	log := logging.Context("shutdown")
	hookCount := len(s.hooks)
	if hookCount < 1 {
		log.Debug("No shutdown hooks were registered")
		return
	}

	log.Info("Shutting down...")

	s.InProgress = true
	defer func() {
		s.InProgress = false
	}()

	ctx := context.TODO()
	shutdownOne := func(h *hook) {
		c, cancel := context.WithTimeout(ctx, h.timeout)
		start := time.Now()
		defer func() {
			if r := recover(); r != nil {
				h.error = errors.New(fmt.Sprintf("shutdown handler paniced with: %v", r))
			}
			h.elapsed = time.Since(start)
			h.done = true
			cancel()
		}()

		if err := h.handler(c); err != nil {
			h.error = err
		}
	}

	for i, h := range s.hooks {
		shutdownOne(h)
		elapsedMillis := h.elapsed.Round(time.Microsecond).String()
		if h.error != nil {
			log.Warn("[%d of %d] %s: failed: %v (%s)",
				i+1, hookCount, h.tag, h.error, elapsedMillis)
		} else {
			log.Info("%d/%d %s: done (%s)",
				i+1, hookCount, h.tag, elapsedMillis)
		}
	}
}

// All is a shortcut for Shutdown.ShutdownAll() on default instance
func All() {
	s.Instance().ShutdownAll()
}

// Register is a shortcut for Shutdown.Register() on default instance
func Register(tag string, timeout time.Duration, handler hookFunc) {
	s.Instance().Register(tag, timeout, handler)
}
