package dbmongo

import (
	"context"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type MongoDBClient interface {
	GetCollection(name string) *mongo.Collection
	GetBucket() (*gridfs.Bucket, error)
	GetCollectionFiles(name string) *mongo.Collection
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
	clientOptions := options.Client().ApplyURI(cfg.DataBaseMongo.DatabaseURL)
	clientOptions.SetMaxConnIdleTime(5 * time.Minute)
	clientOptions.SetMaxConnecting(0)
	clientOptions.SetMaxPoolSize(200)
	clientOptions.SetMinPoolSize(10)
	client, _ = mongo.Connect(ctx, clientOptions)
}

func getClient() *mongo.Client {
	return client
}

func (m mongoDBClient) GetCollection(name string) *mongo.Collection {
	return m.client.Database("gokeeper").Collection(name)
}

func (m mongoDBClient) GetCollectionFiles(name string) *mongo.Collection {
	return m.client.Database("files").Collection(name)
}

func (m mongoDBClient) GetBucket() (*gridfs.Bucket, error) {
	db := m.client.Database("files")
	opts := options.GridFSBucket().SetName("custom name")
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		panic(err)
	}

	return bucket, nil
}
