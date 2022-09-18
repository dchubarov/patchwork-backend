package apiserver

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"twowls.org/patchwork/server/bootstrap/config"
	"twowls.org/patchwork/server/bootstrap/logging"
	"twowls.org/patchwork/server/bootstrap/shutdown"
)

const (
	shutdownTimeout = 10 * time.Second
)

var log = logging.WithComponent("apiserver")

// Start starts the API server.
func Start() {
	cfg := &config.Values().Apiserver
	log.Info().Msgf("Starting on '%s:%d'...", cfg.ListenAddr, cfg.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.ListenAddr, cfg.Port),
		Handler: Router(log),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("failed to start")
		}
	}()

	shutdown.Register("apiserver", shutdownTimeout, srv.Shutdown)
}
