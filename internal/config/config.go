// Package config receives settings when the application starts.
package config

import (
	"flag"

	"github.com/caarlos0/env"
)

// Flags stores application launch settings.
type Flags struct {
	// Host (flag -a) - HTTP server launch address,  e.g. localhost:8080.
	Host string `env:"SERVER_ADDRESS" json:"server_address"`
	// URL (flag -b) - the base address of the resulting shortened URL, e.g.  http://localhost:8080.
	DBDSN string `env:"DATABASE_DSN" json:"database_dsn"`
	// ConfigFileName (flag -c/-config) - the name of configuration file
}

// Default values for flags.
const (
	defHost string = ":8080"
)

// NewConfig creates an instance with settings from flags or environment variables.
func NewConfig() *Flags {
	c := &Flags{}

	flag.StringVar(&c.Host, "a", defHost, "HTTP server start address")
	flag.StringVar(&c.DBDSN, "d", "", "database connection address")
	flag.Parse()

	env.Parse(c)
	return c
}
