// Package config receives settings when the application starts.
package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env"

	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
)

// Flags stores application launch settings.
type Flags struct {
	GRPC           string `env:"GRPC_PORT" json:"grpc"`
	DBURI          string `env:"DATABASE_NAME" json:"database_uri"`
	ConfigFileName string `env:"CONFIG"`
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

	if c.DBURI == "" {
		c.DBURI = conf.DBURI
	}
	if c.GRPC == "" {
		c.GRPC = conf.GRPC
	}

	return nil
}

// NewConfig creates an instance with settings from flags or environment variables.
func NewConfig() *Flags {
	c := &Flags{}

	flag.StringVar(&c.GRPC, "g", defGRPC, "gRPC port")
	flag.StringVar(&c.DBURI, "d", "", "path to the database file")
	flag.StringVar(&c.ConfigFileName, "c", "./config.json", "the name of configuration file")
	flag.StringVar(&c.ConfigFileName, "config", "./config.json", "the name of configuration file")
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
