package logging

import (
	"context"
	"github.com/rs/zerolog"
)

type CtxEnricher func(context.Context) any

// Facade represents logging facade
type Facade interface {
	Logger() *zerolog.Logger
	LoggerCtx(ctx context.Context) *zerolog.Logger

	Trace() *zerolog.Event
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Panic() *zerolog.Event

	TraceCtx(ctx context.Context) *zerolog.Event
	DebugCtx(ctx context.Context) *zerolog.Event
	InfoCtx(ctx context.Context) *zerolog.Event
	WarnCtx(ctx context.Context) *zerolog.Event
	ErrorCtx(ctx context.Context) *zerolog.Event
	PanicCtx(ctx context.Context) *zerolog.Event

	WithComponent(component string) Facade
	WithCtxEnricher(enricher CtxEnricher) Facade
}
