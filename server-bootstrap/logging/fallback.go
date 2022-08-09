package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"twowls.org/patchwork/commons/logging"
)

const (
	fallbackPrefix       = "!noLogger"
	fallbackTraceLevel   = "TRACE"
	fallbackDebugLevel   = "DEBUG"
	fallbackRequestLevel = "REQST"
	fallbackInfoLevel    = "INFO "
	fallbackWarningLevel = "WARN "
	fallbackErrorLevel   = "ERROR"
	fallbackPanicLevel   = "PANIC"
)

type fallbackLogger struct{}

func (l *fallbackLogger) Trace(format string, v ...any) {
	logDefault(false, fallbackTraceLevel, nil, format, v...)
}

func (l *fallbackLogger) Debug(format string, v ...any) {
	logDefault(false, fallbackDebugLevel, nil, format, v...)
}

func (l *fallbackLogger) Request(format string, v ...any) {
	logDefault(false, fallbackRequestLevel, nil, format, v...)
}

func (l *fallbackLogger) Info(format string, v ...any) {
	logDefault(false, fallbackInfoLevel, nil, format, v...)
}

func (l *fallbackLogger) Warn(format string, v ...any) {
	logDefault(false, fallbackWarningLevel, nil, format, v...)
}

func (l *fallbackLogger) Error(format string, v ...any) {
	logDefault(false, fallbackErrorLevel, nil, format, v...)
}

func (l *fallbackLogger) Panic(format string, v ...any) {
	logDefault(true, fallbackPanicLevel, nil, format, v...)
}

func (l *fallbackLogger) TraceFields(fields any, format string, v ...any) {
	logDefault(false, fallbackTraceLevel, fields, format, v...)
}

func (l *fallbackLogger) DebugFields(fields any, format string, v ...any) {
	logDefault(false, fallbackDebugLevel, fields, format, v...)
}

func (l *fallbackLogger) RequestFields(fields any, format string, v ...any) {
	logDefault(false, fallbackRequestLevel, fields, format, v...)
}

func (l *fallbackLogger) InfoFields(fields any, format string, v ...any) {
	logDefault(false, fallbackInfoLevel, fields, format, v...)
}

func (l *fallbackLogger) WarnFields(fields any, format string, v ...any) {
	logDefault(false, fallbackWarningLevel, fields, format, v...)
}

func (l *fallbackLogger) ErrorFields(fields any, format string, v ...any) {
	logDefault(false, fallbackErrorLevel, fields, format, v...)
}

func (l *fallbackLogger) PanicFields(fields any, format string, v ...any) {
	logDefault(true, fallbackPanicLevel, fields, format, v...)
}

func (l *fallbackLogger) IsDebugEnabled() bool {
	return true
}

func (l *fallbackLogger) Context(string) logging.Facade {
	return l
}

func logDefault(isPanic bool, level string, fields any, msg string, v ...any) {
	fieldsOut := ""
	if fields != nil {
		if fieldsJson, err := json.Marshal(fields); err == nil {
			fieldsOut = string(fieldsJson)
		}
	}

	out := fmt.Sprintf(msg, v...)
	_, _ = fmt.Fprintln(os.Stderr, fallbackPrefix, level, out, fieldsOut)
	if isPanic {
		panic(out)
	}
}
