package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
	"twowls.org/patchwork/commons/logging"
	"twowls.org/patchwork/server/bootstrap/config"
)

const (
	componentFieldName = "component"
	rootComponentName  = "main"
)

type defaultFacade struct {
	parent    *defaultFacade
	logger    *zerolog.Logger
	component string
}

var root *defaultFacade

// logging.Facade methods -> defaultFacade

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

func (f *defaultFacade) Tracef(format string, args ...any) {
	f.Trace().Msgf(format, args...)
}

func (f *defaultFacade) Debugf(format string, args ...any) {
	f.Debug().Msgf(format, args...)
}

func (f *defaultFacade) Infof(format string, args ...any) {
	f.Info().Msgf(format, args...)
}

func (f *defaultFacade) Warnf(format string, args ...any) {
	f.Warn().Msgf(format, args...)
}

func (f *defaultFacade) Errorf(format string, args ...any) {
	f.Error().Msgf(format, args...)
}

func (f *defaultFacade) Panicf(format string, args ...any) {
	f.Panic().Msgf(format, args...)
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

func Tracef(format string, args ...any) {
	root.Tracef(format, args...)
}

func Debugf(format string, args ...any) {
	root.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	root.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	root.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	root.Errorf(format, args...)
}

func Panicf(format string, args ...any) {
	root.Panicf(format, args...)
}

func WithComponent(component string) logging.Facade {
	return root.WithComponent(component)
}

// initialization

func init() {
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
		},
	}

	logger := zerolog.New(console).With().
		Str(componentFieldName, prettyComponent(rootComponentName)).
		Timestamp().
		Logger()

	if level, err := zerolog.ParseLevel(config.Values().Logging.Level); err == nil {
		logger = logger.Level(level)
	}

	root = &defaultFacade{
		logger: &logger,
	}
}

// private

func prettyComponent(component string) string {
	return fmt.Sprintf("[%12s]", component)
}
