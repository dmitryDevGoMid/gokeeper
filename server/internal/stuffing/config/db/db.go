package db

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type MongoDBClient struct {
	client *mongo.Client
}

func NewConnectMongoDB(cfg *config.Config) *MongoDBClient {
	initDataBase(cfg)
	return &MongoDBClient{client: client}
}

func initDataBase(cfg *config.Config) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println(cfg.DataBase.DatabaseURL)
	clientOptions := options.Client().ApplyURI(cfg.DataBase.DatabaseURL)
	client, _ = mongo.Connect(ctx, clientOptions)
}

func (m MongoDBClient) GetClient() *mongo.Client {
	return client
}

func (m MongoDBClient) GetCollection(name string) *mongo.Collection {
	return m.client.Database("gokeeper").Collection(name)
}
