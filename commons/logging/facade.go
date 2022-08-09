package logging

// Facade represents logging facade
type Facade interface {
	Trace(msg string, v ...any)
	Debug(msg string, v ...any)
	Request(msg string, v ...any)
	Info(msg string, v ...any)
	Warn(msg string, v ...any)
	Error(msg string, v ...any)
	Panic(msg string, v ...any)

	TraceFields(fields any, msg string, v ...any)
	DebugFields(fields any, msg string, v ...any)
	RequestFields(fields any, msg string, v ...any)
	InfoFields(fields any, msg string, v ...any)
	WarnFields(fields any, msg string, v ...any)
	ErrorFields(fields any, msg string, v ...any)
	PanicFields(fields any, msg string, v ...any)

	Context(name string) Facade

	IsDebugEnabled() bool
}
