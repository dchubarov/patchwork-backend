package logging

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
	"twowls.org/patchwork/commons/logging"
	"twowls.org/patchwork/server/bootstrap/config"
)

const (
	componentFieldName = "component"
	hostFieldName      = "host"
	pidFieldName       = "pid"
	rootComponent      = "main"
)

type defaultFacade struct {
	parent    *defaultFacade
	logger    *zerolog.Logger
	component string
	enricher  logging.CtxEnricher
}

var root *defaultFacade

// logging.Facade methods -> defaultFacade

func (f *defaultFacade) Logger() *zerolog.Logger {
	return f.logger
}

func (f *defaultFacade) LoggerCtx(ctx context.Context) *zerolog.Logger {
	if f.enricher != nil {
		l := f.logger.With().Fields(f.enricher(ctx)).Logger()
		return &l
	} else {
		return f.logger
	}
}

func (f *defaultFacade) Trace() *zerolog.Event {
	return f.logger.Trace()
}

func (f *defaultFacade) Debug() *zerolog.Event {
	return f.logger.Debug()
}

func (f *defaultFacade) Info() *zerolog.Event {
	return f.logger.Info()
}

func (f *defaultFacade) Warn() *zerolog.Event {
	return f.logger.Warn()
}

func (f *defaultFacade) Error() *zerolog.Event {
	return f.logger.Error()
}

func (f *defaultFacade) Panic() *zerolog.Event {
	return f.logger.Panic()
}

func (f *defaultFacade) TraceCtx(ctx context.Context) *zerolog.Event {
	return f.enrichEvent(ctx, f.logger.Trace())
}

func (f *defaultFacade) DebugCtx(ctx context.Context) *zerolog.Event {
	return f.enrichEvent(ctx, f.logger.Debug())
}

func (f *defaultFacade) InfoCtx(ctx context.Context) *zerolog.Event {
	return f.enrichEvent(ctx, f.logger.Info())
}

func (f *defaultFacade) WarnCtx(ctx context.Context) *zerolog.Event {
	return f.enrichEvent(ctx, f.logger.Warn())
}

func (f *defaultFacade) ErrorCtx(ctx context.Context) *zerolog.Event {
	return f.enrichEvent(ctx, f.logger.Error())
}

func (f *defaultFacade) PanicCtx(ctx context.Context) *zerolog.Event {
	return f.enrichEvent(ctx, f.logger.Panic())
}

func (f *defaultFacade) enrichEvent(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if f.enricher != nil {
		return event.Fields(f.enricher(ctx))
	} else {
		return event
	}
}

func (f *defaultFacade) WithComponent(component string) logging.Facade {
	label := component
	for p := f; p != nil; p = p.parent {
		if len(p.component) > 0 {
			label = p.component + "." + label
		}
	}

	subLogger := f.logger.With().
		Str(componentFieldName, prettyComponent(label)).
		Logger()

	return &defaultFacade{
		parent:    f,
		logger:    &subLogger,
		component: component,
		enricher:  f.enricher,
	}
}

func (f *defaultFacade) WithCtxEnricher(enricher logging.CtxEnricher) logging.Facade {
	return &defaultFacade{
		parent:    f,
		logger:    f.logger,
		component: f.component,
		enricher:  enricher,
	}
}

// global convenience functions

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

func TraceCtx(ctx context.Context) *zerolog.Event {
	return root.TraceCtx(ctx)
}

func DebugCtx(ctx context.Context) *zerolog.Event {
	return root.DebugCtx(ctx)
}

func InfoCtx(ctx context.Context) *zerolog.Event {
	return root.InfoCtx(ctx)
}

func WarnCtx(ctx context.Context) *zerolog.Event {
	return root.WarnCtx(ctx)
}

func ErrorCtx(ctx context.Context) *zerolog.Event {
	return root.ErrorCtx(ctx)
}

func PanicCtx(ctx context.Context) *zerolog.Event {
	return root.PanicCtx(ctx)
}

func WithComponent(component string) logging.Facade {
	return root.WithComponent(component)
}

func WithCtxEnricher(enricher logging.CtxEnricher) logging.Facade {
	return root.WithCtxEnricher(enricher)
}

// private

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

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

	logger := zerolog.New(console).With().
		Str(componentFieldName, prettyComponent(rootComponent)).
		Str(hostFieldName, hostname).
		Int(pidFieldName, os.Getpid()).
		Timestamp().
		//Caller().
		Logger()

	if level, err := zerolog.ParseLevel(config.Values().Logging.Level); err == nil {
		logger = logger.Level(level)
	}

	root = &defaultFacade{
		logger:   &logger,
		enricher: defaultCtxEnricher,
	}
}

func prettyComponent(component string) string {
	return fmt.Sprintf("[%12s]", component)
}

func defaultCtxEnricher(ctx context.Context) any {
	requestId, hasRequestId := ctx.Value(logging.RequestCorrelationId).(string)
	jobId, hasJobId := ctx.Value(logging.JobCorrelationId).(string)
	if hasRequestId || hasJobId {
		m := make(map[string]any)
		if hasRequestId {
			m[logging.RequestCorrelationId] = requestId
		}
		if hasJobId {
			m[logging.JobCorrelationId] = jobId
		}
		return m
	}
	return nil
}
