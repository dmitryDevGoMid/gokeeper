package db

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type MongoDBClient interface {
	GetCollection(name string) *mongo.Collection
	GetBucket() *gridfs.Bucket
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

func getClient() *mongo.Client {
	return client
}

func (m mongoDBClient) GetCollection(name string) *mongo.Collection {
	return m.client.Database("gokeeper").Collection(name)
}

func (m mongoDBClient) GetBucket() (*gridfs.Bucket, error) {
	// Получаем экземпляр GridFS
	/*db := m.client.Database("files")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, err
		//log.Fatal(err)
	}

	return bucket, nil*/
	db := m.client.Database("files")
	opts := options.GridFSBucket().SetName("custom name")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}

	return bucket, nil
}
