package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type GoKeeperServerAdress struct {
	GoKeeperAdress string `env:"GOKEEPER_ADDRESS"`
}

type GrpcServerAdress struct {
	AddressGrpc string `env:"GRPC_ADDRESS"`
}

type Server struct {
	Address string `env:"SERVER_ADDRESS"`
}

type DataBase struct {
	DatabaseURL string `env:"DATABASE_MONGO"`
}

type Config struct {
	Server               Server
	DataBase             DataBase
	GrpcServerAdress     GrpcServerAdress
	GoKeeperServerAdress GoKeeperServerAdress
}

var (
	databaseURL          string
	addresServer         string
	addressGrpc          string
	goKeeperServerAdress string
)

func init() {
	flag.StringVar(&goKeeperServerAdress, "adrgokeeper", "localhost:8000", "location http server")
	flag.StringVar(&addressGrpc, "adrgrpc", "localhost:50051", "location http server")
	flag.StringVar(&addresServer, "adr", ":3000", "location http server")
	flag.StringVar(&databaseURL, "dbm", "mongodb://admin:admin@localhost:27017", "database url for conection mongo")
}

// Разбираем конфигурацию по структурам
func ParseConfig() (*Config, error) {
	flag.Parse()

	var config Config

	config.DataBase.DatabaseURL = databaseURL
	config.Server.Address = addresServer
	config.GrpcServerAdress.AddressGrpc = addressGrpc
	config.GoKeeperServerAdress.GoKeeperAdress = goKeeperServerAdress

	//Init by environment variables
	err := env.Parse(&config.DataBase)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.Server)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.GrpcServerAdress)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&config.GoKeeperServerAdress)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
