package main

import (
	"context"
	"github.com/rs/zerolog"
	"runtime"
	"strings"
)

func (ext *ClientExtension) contextLogger(ctx context.Context) *zerolog.Logger {
	l := ext.log.LoggerCtx(ctx)
	if pc, _, _, ok := runtime.Caller(1); ok {
		f := runtime.FuncForPC(pc)
		if f != nil {
			name := f.Name()
			cl := l.With().Str("func", name[strings.LastIndex(name, ".")+1:]).Logger()
			return &cl
		}
	}
	return l
}
