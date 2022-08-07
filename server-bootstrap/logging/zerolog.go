package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
	"twowls.org/patchwork/server/bootstrap/config"
)

const (
	componentFieldName = "component"
	rootComponentName  = "main"
)

type zeroLogFacade struct {
	logger    *zerolog.Logger
	component string
}

func (f *zeroLogFacade) Trace(msg string, v ...any) {
	f.sendEvent(f.logger.Trace(), nil, msg, v...)
}

func (f *zeroLogFacade) Debug(msg string, v ...any) {
	f.sendEvent(f.logger.Debug(), nil, msg, v...)
}

func (f *zeroLogFacade) Request(msg string, v ...any) {
	f.sendEvent(f.logger.Info(), nil, msg, v...)
}

func (f *zeroLogFacade) Info(msg string, v ...any) {
	f.sendEvent(f.logger.Info(), nil, msg, v...)
}

func (f *zeroLogFacade) Warn(msg string, v ...any) {
	f.sendEvent(f.logger.Warn(), nil, msg, v...)
}

func (f *zeroLogFacade) Error(msg string, v ...any) {
	f.sendEvent(f.logger.Error(), nil, msg, v...)
}

func (f *zeroLogFacade) Panic(msg string, v ...any) {
	f.sendEvent(f.logger.Panic(), nil, msg, v...)
}

func (f *zeroLogFacade) TraceFields(fields interface{}, msg string, v ...any) {
	f.sendEvent(f.logger.Trace(), fields, msg, v...)
}

func (f *zeroLogFacade) DebugFields(fields interface{}, msg string, v ...any) {
	f.sendEvent(f.logger.Debug(), fields, msg, v...)
}

func (f *zeroLogFacade) RequestFields(fields interface{}, msg string, v ...any) {
	f.sendEvent(f.logger.Info(), fields, msg, v...)
}

func (f *zeroLogFacade) InfoFields(fields interface{}, msg string, v ...any) {
	f.sendEvent(f.logger.Info(), fields, msg, v...)
}

func (f *zeroLogFacade) WarnFields(fields interface{}, msg string, v ...any) {
	f.sendEvent(f.logger.Warn(), fields, msg, v...)
}

func (f *zeroLogFacade) ErrorFields(fields interface{}, msg string, v ...any) {
	f.sendEvent(f.logger.Error(), fields, msg, v...)
}

func (f *zeroLogFacade) PanicFields(fields interface{}, msg string, v ...any) {
	f.sendEvent(f.logger.Panic(), fields, msg, v...)
}

func (f *zeroLogFacade) Context(name string) Facade {
	return &zeroLogFacade{
		logger:    f.logger,
		component: name,
	}
}

func (f *zeroLogFacade) IsDebugEnabled() bool {
	return f.logger.GetLevel() <= zerolog.DebugLevel
}

func (f *zeroLogFacade) sendEvent(initial *zerolog.Event, fields interface{}, msg string, v ...any) {
	e := initial.Str(componentFieldName, fmt.Sprintf("[%12s]", f.component))
	if fields != nil {
		e = e.Fields(fields)
	}

	if msg != "" {
		if len(v) > 0 {
			e.Msgf(msg, v...)
		} else {
			e.Msg(msg)
		}
	} else {
		e.Send()
	}
}

func zeroLogInit() Facade {
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    config.Values().Log.NoColor,
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

	level, err := zerolog.ParseLevel(config.Values().Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	logger := zerolog.New(console).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &zeroLogFacade{
		component: rootComponentName,
		logger:    &logger,
	}
}
