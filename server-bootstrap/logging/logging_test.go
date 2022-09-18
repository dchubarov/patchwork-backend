package logging

import (
	"context"
	"testing"
)

func TestLogging(t *testing.T) {
	Trace().Msgf("Trace message with variables: %d", 1)
	Debug().Msgf("Debug message with variables: %d", 1)
	Info().Msgf("Info message with variables: %d", 1)
	Warn().Msgf("Warning message with variables: %d", 1)
	Error().Msgf("Error message with variables: %d", 1)
}

func TestLogging_CtxFields(t *testing.T) {
	ctx := context.WithValue(context.Background(), "request-id", "22eaaa6d-ec3e-4dc6-ab74-b319673f928f")

	TraceCtx(ctx).Msgf("Trace with context")
	DebugCtx(ctx).Msgf("Debug with context")
	InfoCtx(ctx).Msgf("Info with context")
	WarnCtx(ctx).Msgf("Warn with context")
	ErrorCtx(ctx).Msgf("Error with context")

	l := WithCtxEnricher(func(ctx context.Context) any {
		if v, ok := ctx.Value("request-id").(string); ok {
			m := make(map[string]any)
			m["correlation-id"] = v
			return m
		}
		return nil
	}).LoggerCtx(ctx)

	l.Info().Msg("Test")
}

func TestLogging_Panic(t *testing.T) {
	defer func() {
		p := recover()
		if p != nil {
			t.Fail()
		}
	}()

	PanicCtx(context.TODO())
}
