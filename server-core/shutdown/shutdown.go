package shutdown

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"sync"
	"time"
	"twowls.org/patchwork/backend/logging"
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
	shutdown Shutdown
	once     sync.Once
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

// Instance returns default instance of Shutdown
func Instance() *Shutdown {
	once.Do(func() {
		shutdown = Shutdown{}
	})
	return &shutdown
}

// Register is a shortcut for Shutdown.Register() on default instance
func Register(tag string, timeout time.Duration, handler hookFunc) {
	Instance().Register(tag, timeout, handler)
}

// All is a shortcut for Shutdown.ShutdownAll() on default instance
func All() {
	Instance().ShutdownAll()
}
