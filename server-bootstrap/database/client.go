package database

import (
	"context"
	"strings"
	"time"
	"twowls.org/patchwork/commons/utils/singleton"
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
	client = singleton.NewLazy(func() ClientMethods {
		var c ClientMethods
		cfg := config.Values().Database
		if strings.HasPrefix(cfg.Url, "mongodb://") {
			c = mongo.New(cfg)
		} else {
			log.Panic("connection url specifies unknown database type: %s", cfg.Url)
		}

		if c != nil {
			shutdown.Register("database", 3*time.Second, c.Disconnect)
		}

		return c
	})
)

func Client() ClientMethods {
	return client.Instance()
}
