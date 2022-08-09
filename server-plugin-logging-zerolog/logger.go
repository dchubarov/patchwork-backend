package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/commons/logging"
)

const (
	componentFieldName = "component"
	rootComponentName  = "main"
)

type zerologExtension struct {
	logger    *zerolog.Logger
	component string
}

// Extension methods

func (ext *zerologExtension) Configure(options *extension.Options) error {
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    options.BoolConfigDefault("noColor", false),
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

	logger := zerolog.New(console)
	if levelStr, found := options.StrConfig("level"); found {
		if level, err := zerolog.ParseLevel(levelStr); err == nil {
			logger = logger.Level(level)
		}
	}

	logger = logger.With().
		Timestamp().
		Logger()

	ext.component = rootComponentName
	ext.logger = &logger
	return nil
}

// logging.Facade methods

func (ext *zerologExtension) Trace(msg string, v ...any) {
	ext.sendEvent(ext.logger.Trace(), nil, msg, v...)
}

func (ext *zerologExtension) Debug(msg string, v ...any) {
	ext.sendEvent(ext.logger.Debug(), nil, msg, v...)
}

func (ext *zerologExtension) Request(msg string, v ...any) {
	ext.sendEvent(ext.logger.Info(), nil, msg, v...)
}

func (ext *zerologExtension) Info(msg string, v ...any) {
	ext.sendEvent(ext.logger.Info(), nil, msg, v...)
}

func (ext *zerologExtension) Warn(msg string, v ...any) {
	ext.sendEvent(ext.logger.Warn(), nil, msg, v...)
}

func (ext *zerologExtension) Error(msg string, v ...any) {
	ext.sendEvent(ext.logger.Error(), nil, msg, v...)
}

func (ext *zerologExtension) Panic(msg string, v ...any) {
	ext.sendEvent(ext.logger.Panic(), nil, msg, v...)
}

func (ext *zerologExtension) TraceFields(fields any, msg string, v ...any) {
	ext.sendEvent(ext.logger.Trace(), fields, msg, v...)
}

func (ext *zerologExtension) DebugFields(fields any, msg string, v ...any) {
	ext.sendEvent(ext.logger.Debug(), fields, msg, v...)
}

func (ext *zerologExtension) RequestFields(fields any, msg string, v ...any) {
	ext.sendEvent(ext.logger.Info(), fields, msg, v...)
}

func (ext *zerologExtension) InfoFields(fields any, msg string, v ...any) {
	ext.sendEvent(ext.logger.Info(), fields, msg, v...)
}

func (ext *zerologExtension) WarnFields(fields any, msg string, v ...any) {
	ext.sendEvent(ext.logger.Warn(), fields, msg, v...)
}

func (ext *zerologExtension) ErrorFields(fields any, msg string, v ...any) {
	ext.sendEvent(ext.logger.Error(), fields, msg, v...)
}

func (ext *zerologExtension) PanicFields(fields any, msg string, v ...any) {
	ext.sendEvent(ext.logger.Panic(), fields, msg, v...)
}

func (ext *zerologExtension) Context(name string) logging.Facade {
	return &zerologExtension{
		logger:    ext.logger,
		component: name,
	}
}

func (ext *zerologExtension) IsDebugEnabled() bool {
	return ext.logger.GetLevel() <= zerolog.DebugLevel
}

func (ext *zerologExtension) sendEvent(initial *zerolog.Event, fields any, msg string, v ...any) {
	e := initial.Str(componentFieldName, fmt.Sprintf("[%12s]", ext.component))
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
