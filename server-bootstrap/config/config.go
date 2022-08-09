package config

import (
	"fmt"
	"github.com/vrischmann/envconfig"
	"twowls.org/patchwork/commons/utils/singleton"
)

// Apiserver contains configuration of API server
type Apiserver struct {
	// ListenAddr server listen address
	ListenAddr string `envconfig:"optional"`
	// Port server port
	Port uint16 `envconfig:"default=8080"`
}

// Database contains database configuration values
type Database struct {
	// Url database connection URL
	Url string
	// Username database username (empty if supplied as part of Url)
	Username string `envconfig:"optional"`
	// Password database password (empty if supplied as part of Url)
	Password string `envconfig:"optional"`
}

// Logging contains logging configuration
type Logging struct {
	// Level specifies logging level
	Level string `envconfig:"default=info"`
	// NoColor specifies whether not to use ANSI colors in console
	NoColor bool `envconfig:"optional"`
	// Plugin specifies logging plugin name
	Plugin string `envconfig:"optional"`
}

// Root contains application configuration
type Root struct {
	// Apiserver api server configuration
	Apiserver Apiserver
	// Database database configuration
	Database Database
	// Logging logging configuration
	Logging Logging
	// PluginsDir specifies plugins directory
	PluginsDir string
}

var values = singleton.NewLazy(load)

func load() *Root {
	c := new(Root)
	if err := envconfig.Init(c); err != nil {
		// cannot use logging here since it might not be initialised yet
		panic(fmt.Sprintf("unable to load configuration: %v", err))
	}
	return c
}

func Values() *Root {
	return values.Instance()
}
