package logging

import "github.com/rs/zerolog"

// Facade represents logging facade
type Facade interface {
	Trace() *zerolog.Event
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Panic() *zerolog.Event

	WithComponent(component string) Facade
}
