package logging

import (
	"errors"
	"strings"
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/commons/logging"
	"twowls.org/patchwork/commons/singleton"
	"twowls.org/patchwork/server/bootstrap/config"
	"twowls.org/patchwork/server/bootstrap/plugins"
)

const loggingPluginPrefix = "logging-"

var rootLogger = singleton.Lazy(func() logging.Facade {
	var err error
	if pluginName := config.Values().Logging.Plugin; pluginName != "" {
		var info extension.PluginInfo
		if info, err = plugins.Load(loggingPluginPrefix + strings.ToLower(pluginName)); err == nil {
			if ext := info.DefaultExtension(); ext != nil {
				if logger, ok := ext.(logging.Facade); ok {
					options := extension.EmptyOptions().
						PutConfig("level", config.Values().Logging.Level).
						PutConfig("noColor", config.Values().Logging.NoColor)

					if err = ext.Configure(options); err == nil {
						logger.Info("Logging is provided via plugin: %q (%s)", pluginName, info.Description())
						return logger
					}
				}
			}

			if err == nil {
				err = errors.New("invalid extension")
			}
		}
	}

	logger := &fallbackLogger{}
	if err != nil {
		logger.Error("Logging plugin failed to load: %v", err)
	} else {
		logger.Warn("Logging plugin is not configured")
	}

	return logger
})

func Root() logging.Facade {
	return rootLogger.Instance()
}

// Package-level logging facade (convenience shortcuts)

// Trace is a shortcut for Facade.Trace() on root logger
func Trace(format string, v ...any) {
	Root().Debug(format, v...)
}

// Debug is a shortcut for Facade.Debug() on root logger
func Debug(format string, v ...any) {
	Root().Debug(format, v...)
}

// Request is a shortcut for Facade.Request() on root logger
func Request(format string, v ...any) {
	Root().Request(format, v...)
}

// Info is a shortcut for Facade.Info() on root logger
func Info(format string, v ...any) {
	Root().Info(format, v...)
}

// Warn is a shortcut for Facade.Warn() on root logger
func Warn(format string, v ...any) {
	Root().Warn(format, v...)
}

// Error is a shortcut for Facade.Error() on root logger
func Error(format string, v ...any) {
	Root().Error(format, v...)
}

// Panic is a shortcut for Facade.Panic() on root logger
func Panic(format string, v ...any) {
	Root().Panic(format, v...)
}

// Context is a shortcut for Facade.Context() on root logger
func Context(name string) logging.Facade {
	return Root().Context(name)
}
