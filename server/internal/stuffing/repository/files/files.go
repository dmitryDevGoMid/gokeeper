package files

import (
	"context"
	"errors"
	"time"

	db "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/config/db"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrDataNotFound = errors.New("data not found")
)

type Metadata struct {
	ClientID  string `bson:"client_id"`
	CountPart int    `bson:"count_part"`
	UID       string `bson:"uid"`
}

type Files struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ID_File   string             `json:"id_file" bson:"id_file,omitempty"`
	Filename  string             `json:"filename" bson:"filename"`
	ChunkSize int64              `json:"chunk_size" bson:"chunkSize"`
	//Metadata   string             `json:"metadata" bson: "metadata"`
	Metadata   Metadata  `json:"metadata" bson:"metadata"`
	UploadDate time.Time `json:"upload_date" bson:"uploadDate"`
	Length     int64     `json:"length" bson:"length"`
}

type FilesRepository interface {
	GetByUserIdListFiles(ctx context.Context, user *user.User) (*[]Files, error)
	DeleteFilesByID(ctx context.Context, id string) error
}

type filesRepository struct {
	client db.MongoDBClient
}

func NewFilesRepository(client db.MongoDBClient) FilesRepository {
	return &filesRepository{client: client}
}
