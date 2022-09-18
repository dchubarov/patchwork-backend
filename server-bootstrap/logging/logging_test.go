package logging

import "testing"

func TestLogging(t *testing.T) {
	Trace().Msgf("Trace message with variables: %d", 1)
	Debug().Msgf("Debug message with variables: %d", 1)
	Info().Msgf("Info message with variables: %d", 1)
	Warn().Msgf("Warning message with variables: %d", 1)
	Error().Msgf("Error message with variables: %d", 1)
}
