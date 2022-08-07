package config

import (
	"fmt"
	"github.com/vrischmann/envconfig"
	"sync"
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

type Log struct {
	// Level specifies logging level
	Level string `envconfig:"default=info"`
	// NoColor specifies whether not to use ANSI colors in console
	NoColor bool `envconfig:"optional"`
}

// Root contains application configuration
type Root struct {
	// Apiserver api server configuration
	Apiserver Apiserver
	// Database database configuration
	Database Database
	// Log logging configuration
	Log Log
}

var (
	once sync.Once
	root *Root
)

func Values() *Root {
	once.Do(func() {
		root = new(Root)
		if err := envconfig.Init(root); err != nil {
			// cannot use logging here since it is not initialized yet
			panic(fmt.Sprintf("unable to load configuration: %v", err))
		}
	})
	return root
}
