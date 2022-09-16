package logging

import "testing"

func TestLogging(t *testing.T) {
	Tracef("Trace message with variables: %d", 1)
	Debugf("Debug message with variables: %d", 1)
	Infof("Info message with variables: %d", 1)
	Warnf("Warning message with variables: %d", 1)
	Errorf("Error message with variables: %d", 1)
}
