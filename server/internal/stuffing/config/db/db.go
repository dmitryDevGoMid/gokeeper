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

type MongoDBClient interface {
	GetCollection(name string) *mongo.Collection
	GetClient() *mongo.Client
}

type mongoDBClient struct {
	client *mongo.Client
}

func NewConnectMongoDB(cfg *config.Config) *mongoDBClient {
	initDataBase(cfg)
	return &mongoDBClient{client: client}
}

func initDataBase(cfg *config.Config) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println(cfg.DataBase.DatabaseURL)
	clientOptions := options.Client().ApplyURI(cfg.DataBase.DatabaseURL)
	client, _ = mongo.Connect(ctx, clientOptions)
}

func (m mongoDBClient) GetClient() *mongo.Client {
	return client
}

func (m mongoDBClient) GetCollection(name string) *mongo.Collection {
	return m.client.Database("gokeeper").Collection(name)
}
