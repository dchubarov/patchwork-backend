package logging

import (
	"context"
	"github.com/rs/zerolog"
)

const (
	CorrelationRequestId = "request-id"
	CorrelationJobId     = "job-id"
)

// Facade represents logging facade
type Facade interface {
	// WithComponent creates a child logger with specified component name
	WithComponent(string) Facade
	// WithContext creates a child logger with correlation fields from context
	WithContext(context.Context) Facade
	// Trace creates a new logging event at trace level
	Trace() *zerolog.Event
	// Debug creates a new logging event at debug level
	Debug() *zerolog.Event
	// Info creates a new logging event at info level
	Info() *zerolog.Event
	// Warn creates a new logging event at warn level
	Warn() *zerolog.Event
	// Error creates a new logging event at error level
	Error() *zerolog.Event
	// Panic creates a new logging event at panic level
	Panic() *zerolog.Event
}
