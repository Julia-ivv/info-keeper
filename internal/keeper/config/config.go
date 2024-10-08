// Пакет config получает настройки при запуске приложения.
package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env"

	"github.com/Julia-ivv/info-keeper.git/pkg/logger"
)

// Flags хранит настройки приложения.
type Flags struct {
	// GRPC (флаг -g) - порт для gRPC, например :3200.
	GRPC string `env:"GRPC_PORT" json:"grpc"`
	// DBDSN (флаг -d) - имя для доступа к БД.
	DBDSN string `env:"DATABASE_DSN" json:"database_dsn"`
	// ConfigFileName (флаг -c/-config) - имя файла конфигурации.
	ConfigFileName string `env:"CONFIG"`
	// SecretKey ключ для создания токена.
	SecretKey string `env:"SKEY" json:"key"`
}

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

// NewConfig создает объект с настройками приложения из флагов или переменных окружения.
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
