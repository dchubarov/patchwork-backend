package database

import (
	"context"
	"strings"
	"sync"
	"time"
	"twowls.org/patchwork/server/bootstrap/config"
	"twowls.org/patchwork/server/bootstrap/database/mongo"
	"twowls.org/patchwork/server/bootstrap/logging"
	"twowls.org/patchwork/server/bootstrap/shutdown"
)

type ClientMethods interface {
	// Connect establishes database connection
	Connect(ctx context.Context)

	// Disconnect shuts down database connection
	Disconnect(ctx context.Context) error
}

var (
	log    = logging.Context("database")
	client ClientMethods
	once   sync.Once
)

func Client() ClientMethods {
	once.Do(func() {
		cfg := config.Values().Database
		if strings.HasPrefix(cfg.Url, "mongodb://") {
			client = mongo.New(cfg)
		} else {
			log.Panic("connection url specifies unknown database type: %s", cfg.Url)
		}

		if client != nil {
			shutdown.Register("database", 3*time.Second, client.Disconnect)
		}
	})
	return client
}