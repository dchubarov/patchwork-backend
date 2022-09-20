package logging

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
	"twowls.org/patchwork/commons/logging"
	"twowls.org/patchwork/commons/util/singleton"
	"twowls.org/patchwork/server/bootstrap/config"
)

const (
	componentFieldName = "component"
	hostFieldName      = "hostname"
	pidFieldName       = "pid"

	rootComponent = "main"
	unknownHost   = "unknown"
)

type defaultFacade struct {
	logger    singleton.S[*zerolog.Logger]
	parent    *defaultFacade
	enricher  func(context.Context, zerolog.Context) zerolog.Context
	component string
}

var root = &defaultFacade{
	logger:   singleton.Eager(newRootLogger),
	enricher: defaultContextEnricher,
}

// logging.Facade methods

func (f *defaultFacade) WithComponent(component string) logging.Facade {
	if len(component) < 1 {
		return f
	}

	child := f.newChild(func() *zerolog.Logger {
		componentPath := component
		for p := f; p != nil; p = p.parent {
			if len(p.component) > 0 {
				componentPath = p.component + "." + componentPath
			}
		}

		l := f.logger.Instance().With().
			Str("component", formatComponentName(componentPath)).
			Logger()

		return &l
	})

	child.component = component
	return child
}

func (f *defaultFacade) WithContext(ctx context.Context) logging.Facade {
	return f.newChild(func() *zerolog.Logger {
		lc := f.logger.Instance().With()
		for p := f; p != nil; p = p.parent {
			if p.enricher != nil {
				lc = p.enricher(ctx, lc)
			}
		}
		l := lc.Logger()
		return &l
	})
}

func (f *defaultFacade) Trace() *zerolog.Event {
	return f.logger.Instance().Trace()
}

func (f *defaultFacade) Debug() *zerolog.Event {
	return f.logger.Instance().Debug()
}

func (f *defaultFacade) Info() *zerolog.Event {
	return f.logger.Instance().Info()
}

func (f *defaultFacade) Warn() *zerolog.Event {
	return f.logger.Instance().Warn()
}

func (f *defaultFacade) Error() *zerolog.Event {
	return f.logger.Instance().Error()
}

func (f *defaultFacade) Panic() *zerolog.Event {
	return f.logger.Instance().Panic()
}

// root logger interface

func WithComponent(component string) logging.Facade {
	return root.WithComponent(component)
}

func WithContext(ctx context.Context) logging.Facade {
	return root.WithContext(ctx)
}

func Trace() *zerolog.Event {
	return root.Trace()
}

func Debug() *zerolog.Event {
	return root.Debug()
}

func Info() *zerolog.Event {
	return root.Info()
}

func Warn() *zerolog.Event {
	return root.Warn()
}

func Error() *zerolog.Event {
	return root.Error()
}

func Panic() *zerolog.Event {
	return root.Panic()
}

// private

func (f *defaultFacade) newChild(loggerFactory func() *zerolog.Logger) *defaultFacade {
	return &defaultFacade{
		logger: singleton.Lazy(loggerFactory),
		parent: f,
	}
}

func newRootLogger() *zerolog.Logger {
	pid := os.Getpid()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = unknownHost
	}

	var loggerInit zerolog.Context
	if config.Values().Logging.Pretty {
		console := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    config.Values().Logging.NoColor,
			TimeFormat: time.Stamp,
			PartsOrder: []string{
				zerolog.TimestampFieldName,
				zerolog.LevelFieldName,
				componentFieldName,
				zerolog.CallerFieldName,
				zerolog.MessageFieldName,
			},
			FieldsExclude: []string{
				componentFieldName,
				hostFieldName,
				pidFieldName,
			},
		}
		loggerInit = zerolog.New(console).With().Caller()
	} else {
		loggerInit = zerolog.New(zerolog.MultiLevelWriter(os.Stdout)).With()
	}

	loggerInit = loggerInit.
		Str(componentFieldName, formatComponentName(rootComponent)).
		Str(hostFieldName, hostname).
		Int(pidFieldName, pid).
		Timestamp()

	logger := loggerInit.Logger()
	if level, err := zerolog.ParseLevel(config.Values().Logging.Level); err == nil {
		logger = logger.Level(level)
	}

	return &logger
}

func defaultContextEnricher(ctx context.Context, lc zerolog.Context) zerolog.Context {
	// add correlation request id as logger field
	if requestId, ok := ctx.Value(logging.CorrelationRequestId).(string); ok {
		lc = lc.Str(logging.CorrelationRequestId, requestId)
	}

	// add correlation job id as logger field
	if jobId, ok := ctx.Value(logging.CorrelationJobId).(string); ok {
		lc = lc.Str(logging.CorrelationJobId, jobId)
	}

	return lc
}

func formatComponentName(component string) string {
	if config.Values().Logging.Pretty {
		return fmt.Sprintf("[%16s]", component)
	} else {
		return component
	}
}
