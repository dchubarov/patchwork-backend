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

var log = logging.Context("apiserver")

// Start starts the API server.
func Start() {
	cfg := &config.Values().Apiserver
	log.Info("Starting on '%s:%d'...", cfg.ListenAddr, cfg.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.ListenAddr, cfg.Port),
		Handler: Router(log),
	}

	go func() {
		shutdown.Register("apiserver", shutdownTimeout, srv.Shutdown)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start: %v", err)
		}
	}()
}
