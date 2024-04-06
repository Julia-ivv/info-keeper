// Package config receives settings when the application starts.
package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"

	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
	"github.com/caarlos0/env"
)

// Flags stores application launch settings.
type Flags struct {
	// GRPC (flag -g) - port for gRPC, e.g. :3200.
	GRPC string `env:"GRPC_PORT" json:"grpc"`
	// URL (flag -b) - the base address of the resulting shortened URL, e.g.  http://localhost:8080.
	DBDSN string `env:"DATABASE_DSN" json:"database_dsn"`
	// ConfigFileName (flag -c/-config) - the name of configuration file
	ConfigFileName string `env:"CONFIG"`
	SecretKey      string `env:"SKEY" json:"key"`
}

// Default values for flags.
const (
	defGRPC string = ":3200"
)

func readFromConf(c *Flags) error {
	f, err := os.Open(c.ConfigFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	allData := []byte{}
	for scan.Scan() {
		allData = append(allData, scan.Bytes()...)
	}
	if err = scan.Err(); err != nil {
		return err
	}

	var conf Flags
	err = json.Unmarshal(allData, &conf)
	if err != nil {
		return err
	}

	if c.DBDSN == "" {
		c.DBDSN = conf.DBDSN
	}
	if c.GRPC == "" {
		c.GRPC = conf.GRPC
	}
	if c.SecretKey == "" {
		c.SecretKey = conf.SecretKey
	}

	return nil
}

// NewConfig creates an instance with settings from flags or environment variables.
func NewConfig() *Flags {
	c := &Flags{}

	flag.StringVar(&c.GRPC, "g", defGRPC, "gRPC port")
	flag.StringVar(&c.DBDSN, "d", "", "database connection address")
	flag.StringVar(&c.ConfigFileName, "c", "", "the name of configuration file")
	flag.StringVar(&c.ConfigFileName, "config", "", "the name of configuration file")
	flag.StringVar(&c.SecretKey, "k", "", "secret key for token")
	flag.Parse()

	env.Parse(c)

	if c.ConfigFileName != "" {
		err := readFromConf(c)
		if err != nil {
			logger.ZapSugar.Infow("reading configuration file", err)
		}
	}

	return c
}
