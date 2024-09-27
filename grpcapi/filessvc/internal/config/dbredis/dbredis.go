package dbredis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/config"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

type RedisDBClient interface {
	GetNewClient() *redis.Client
	GetClient() *redis.Client
}

type redisDBClient struct {
	cfg *config.Config
}

func NewConnectRedis(cfg *config.Config) RedisDBClient {
	initDataBase(cfg)
	return &redisDBClient{cfg: cfg}
}

func initDataBase(cfg *config.Config) {
	client = Connect(cfg)
	// Получение информации о сервере
	info, err := client.Info(context.TODO(), "server").Result()
	if err != nil {
		log.Fatal(err)
	}

	// Вывод информации о сервере
	fmt.Println(info)
}

func (red *redisDBClient) GetClient() *redis.Client {
	return client
}

func Connect(cfg *config.Config) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:        cfg.DataBaseRedis.DatabaseURL,
		Password:    "",               // no password set
		DB:          0,                // use default DB
		PoolSize:    10,               // set pool size to 10 connections
		DialTimeout: 60 * time.Second, // set dial timeout to 5 seconds
		ReadTimeout: 60 * time.Second, // set read timeout to 5 seconds
	})

	return client
}

func (red *redisDBClient) GetNewClient() *redis.Client {
	rdb := Connect(red.cfg)
	return rdb
}
