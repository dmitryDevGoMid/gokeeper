package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type GrpcServerPort struct {
	Port string `env:"GRPC_PORT"`
}

type Server struct {
	Address string `env:"SERVER_ADDRESS"`
}

type DataBaseRedis struct {
	DatabaseURL string `env:"DATABASE_REDIS"`
}

type DataBaseMongo struct {
	DatabaseURL string `env:"DATABASE_MONGO"`
}

type Config struct {
	Server         Server
	DataBaseMongo  DataBaseMongo
	DataBaseRedis  DataBaseRedis
	GrpcServerPort GrpcServerPort
}

var (
	databaseMongoURL string
	databaseRedisURL string
	addresServer     string
	addressGrpcPort  string
)

func init() {
	flag.StringVar(&addressGrpcPort, "adrgrpc", ":50051", "location http server")
	flag.StringVar(&addresServer, "adr", "localhost:3000", "location http server")
	flag.StringVar(&databaseMongoURL, "dbm", "mongodb://admin:admin@localhost:27017", "database url for conection mongo")
	flag.StringVar(&databaseRedisURL, "dbr", "localhost:6379", "database url for conection mongo")
}

// Разбираем конфигурацию по структурам
func ParseConfig() (*Config, error) {
	flag.Parse()

	var config Config

	config.DataBaseMongo.DatabaseURL = databaseMongoURL
	config.Server.Address = addresServer
	config.GrpcServerPort.Port = addressGrpcPort

	//Init by environment variables
	err := env.Parse(&config.DataBaseMongo)
	if err != nil {
		return nil, err
	}

	//Init by environment variables
	err = env.Parse(&config.DataBaseRedis)
	if err != nil {
		return nil, err
	}

	//Init by environment variables
	err = env.Parse(&config.GrpcServerPort)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
