package logging

import (
	"context"
	"testing"
	"twowls.org/patchwork/commons/logging"
)

func TestLogging_RootFunc(t *testing.T) {
	Trace().Msg("Trace message")
	Debug().Msg("Debug message")
	Info().Msg("Info message")
	Warn().Msg("Warning message")
	Error().Msg("Error message")
}
func TestLogging_ChildComponent(t *testing.T) {
	log := WithComponent("testing")
	log.Info().Msg("Logging on sub-component")
}

func TestLogging_Child(t *testing.T) {
	ctx := context.WithValue(context.Background(), logging.CorrelationJobId, "a9cf16d")
	log := WithContext(ctx)
	log.Info().Msg("Logging with context fields")
}
