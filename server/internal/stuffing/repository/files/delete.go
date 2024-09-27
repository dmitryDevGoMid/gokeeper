package files

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *filesRepository) DeleteFilesByID(ctx context.Context, id string) error {

	// Установите соединение с MongoDB
	client := r.client.GetClient()

	//defer client.Disconnect(ctx)
	opts := options.GridFSBucket().SetName("custom name")

	// Инициализируйте GridFS
	bucket, err := gridfs.NewBucket(client.Database("files"), opts)
	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
	}

	// ID файла, который нужно удалить
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	// Удаление файла из GridFS по его ID
	err = bucket.Delete(fileID)
	if err != nil {
		log.Fatal(err)
	}

	return nil

}
