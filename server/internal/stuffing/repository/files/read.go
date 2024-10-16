package files

import (
	"context"
	"errors"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *filesRepository) GetByUserIdListFiles(ctx context.Context, user *user.User) (*[]Files, error) {
	// Установите соединение с MongoDB
	client := r.client.GetClient()

	//defer client.Disconnect(ctx)
	opts := options.GridFSBucket().SetName("custom name")

	// Инициализируйте GridFS
	bucket, err := gridfs.NewBucket(client.Database("files"), opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrDataNotFound
		}
		return nil, err
	}

	// Определите фильтр для поиска файлов по полю metadata
	filter := bson.D{
		{Key: "metadata.client_id", Value: user.ID_User},
	}

	// Определяем опции для поиска
	findOpts := options.GridFSFind()
	var results []Files

	// Находим файлы, соответствующие фильтру
	cursor, err := bucket.Find(filter, findOpts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrDataNotFound
		}
		return nil, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}
