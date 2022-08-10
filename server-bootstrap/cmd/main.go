package main

import (
	"os"
	"os/signal"
	"syscall"
	"twowls.org/patchwork/server/bootstrap/apiserver"
	"twowls.org/patchwork/server/bootstrap/database"
	"twowls.org/patchwork/server/bootstrap/logging"
	"twowls.org/patchwork/server/bootstrap/shutdown"
)

func main() {
	defer shutdown.All()
	database.MustConnect()
	apiserver.Start()
	awaitTermination()
}

func awaitTermination() os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	logging.Info("Received interrupt signal: %v", s)
	return s
}
