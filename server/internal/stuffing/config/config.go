package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Server struct {
	Address string `env:"SERVER_ADDRESS"`
}

type DataBase struct {
	DatabaseURL string `env:"DATABASE_MONGO"`
}

type Config struct {
	Server   Server
	DataBase DataBase
}

var (
	databaseURL  string
	addresServer string
)

func init() {
	flag.StringVar(&addresServer, "adr", "localhost:8000", "location http server")
	flag.StringVar(&databaseURL, "dbm", "mongodb://admin:admin@localhost:27017", "database url for conection mongo")
}

// Разбираем конфигурацию по структурам
func ParseConfig() (*Config, error) {
	flag.Parse()

	var config Config

	config.DataBase.DatabaseURL = databaseURL
	config.Server.Address = addresServer

	//Init by environment variables
	err := env.Parse(&config.DataBase)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.Server)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
