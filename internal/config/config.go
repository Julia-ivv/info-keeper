// Package config receives settings when the application starts.
package config

import (
	"flag"

	"github.com/caarlos0/env"
)

// Flags stores application launch settings.
type Flags struct {
	// GRPC (flag -g) - port for gRPC, e.g. :3200.
	GRPC string `env:"GRPC_PORT" json:"grpc"`
	// URL (flag -b) - the base address of the resulting shortened URL, e.g.  http://localhost:8080.
	DBDSN string `env:"DATABASE_DSN" json:"database_dsn"`
	// ConfigFileName (flag -c/-config) - the name of configuration file
}

// Default values for flags.
const (
	defGRPC string = ":3200"
)

// NewConfig creates an instance with settings from flags or environment variables.
func NewConfig() *Flags {
	c := &Flags{}

	flag.StringVar(&c.GRPC, "g", defGRPC, "gRPC port")
	flag.StringVar(&c.DBDSN, "d", "", "database connection address")
	flag.Parse()

	env.Parse(c)
	return c
}
