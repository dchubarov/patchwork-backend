package logging

import (
	"sync"
)

// Facade represents logging facade
type Facade interface {
	Trace(msg string, v ...any)
	Debug(msg string, v ...any)
	Request(msg string, v ...any)
	Info(msg string, v ...any)
	Warn(msg string, v ...any)
	Error(msg string, v ...any)
	Panic(msg string, v ...any)

	TraceFields(fields interface{}, msg string, v ...any)
	DebugFields(fields interface{}, msg string, v ...any)
	RequestFields(fields interface{}, msg string, v ...any)
	InfoFields(fields interface{}, msg string, v ...any)
	WarnFields(fields interface{}, msg string, v ...any)
	ErrorFields(fields interface{}, msg string, v ...any)
	PanicFields(fields interface{}, msg string, v ...any)

	Context(name string) Facade

	IsDebugEnabled() bool
}

var (
	root Facade
	once sync.Once
)

func Root() Facade {
	once.Do(func() {
		root = zeroLogInit()
	})
	return root
}

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
func Context(name string) Facade {
	return Root().Context(name)
}
