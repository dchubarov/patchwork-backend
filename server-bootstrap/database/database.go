package database

import (
	"context"
	"regexp"
	"time"
	"twowls.org/patchwork/commons/database"
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/commons/util/singleton"
	"twowls.org/patchwork/server/bootstrap/config"
	"twowls.org/patchwork/server/bootstrap/logging"
	"twowls.org/patchwork/server/bootstrap/plugins"
	"twowls.org/patchwork/server/bootstrap/shutdown"
)

const databasePluginPrefix = "database-"

var (
	log      = logging.WithComponent("database")
	dbClient = singleton.Lazy(func() database.Client {
		cfg := config.Values().Database

		schemeRegexp := regexp.MustCompile("^(\\w+)://")
		scheme := schemeRegexp.FindStringSubmatch(cfg.Url)
		if scheme == nil || len(scheme) < 2 {
			log.Panic().Msg("Cannot determine database connection schema from URI")
		}

		if p, err := plugins.Load(databasePluginPrefix + scheme[1]); err == nil {
			if clientExt := p.DefaultExtension(); clientExt != nil {
				if client, ok := clientExt.(database.Client); ok {
					opts := extension.BasicOptions(false, log.WithComponent(scheme[1])).
						PutConfig("uri", cfg.Url).
						PutConfig("username", cfg.Username).
						PutConfig("password", cfg.Password)

					if err := clientExt.Configure(opts); err != nil {
						log.Panic().Err(err).Msg("Could not configure database")
					}

					return client
				}
			}
		}

		log.Panic().Msgf("Unable to initialize database plugin for scheme %q", scheme[1])
		return nil
	})
)

func Client() database.Client {
	return dbClient.Instance()
}

func MustConnect() {
	if err := Client().Connect(context.TODO()); err != nil {
		log.Panic().Err(err).Msg("Database connection failed")
	}

	shutdown.Register("database", 10*time.Second, Client().Disconnect)
}
