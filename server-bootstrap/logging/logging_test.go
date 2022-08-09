package logging

import "testing"

func TestLogging(t *testing.T) {
	Trace("Trace message with variables: %d", 1)
	Debug("Debug message with variables: %d", 1)
	Request("Request message with variables: %d", 1)
	Info("Info message with variables: %d", 1)
	Warn("Warning message with variables: %d", 1)
	Error("Error message with variables: %d", 1)
}
