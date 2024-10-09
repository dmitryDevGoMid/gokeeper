package user

import (
	"context"
	"errors"

	db "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/config/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrDataNotFound = errors.New("data not found")
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ID_User      string             `json:"id_user" bson:"id_user,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Password     string             `json:"password" bson:"password"`
	Token        string             `json:"token" bson:"token"`
	TokenRefresh string             `json:"token_refresh" bson:"token_refresh"`
	Request      string             `json:"request" bson:"request"`
	PublicKey    string             `json:"publickey" bson:"publickey`
}

type ResponseSaveData struct {
	ID   string `json:"id" bson:"id_str,omitempty"`
	Data string `json:"data" bson:"data"`
}

type SaveData struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	User_ID  string             `json:"user_id" bson:"user_id"`
	TypeData string             `json:"type_data" bson:"type_data"`
	Data     string             `json:"data" bson:"data"`
}

type UserRepository interface {
	Create(context.Context, *User) error
	GetByUsername(context.Context, string) (*User, error)
	UpdatePassword(context.Context, string, string) error
	DeleteByUsername(context.Context, string) error

	//Password
	CreatePasswordByUser(ctx context.Context, user *SaveData) error
	GetPasswordByUser(ctx context.Context, user *User) (*[]ResponseSaveData, error)
	UpdatePasswordByKey(ctx context.Context, user *SaveData) error
	DelerePasswordById(ctx context.Context, id primitive.ObjectID) error

	//Cards
	CreateCardByUser(ctx context.Context, user *SaveData) error
	UpdateCardByKey(ctx context.Context, user *SaveData) error
	GetCardsByUser(ctx context.Context, user *User) (*[]ResponseSaveData, error)
	DelereCardById(ctx context.Context, id primitive.ObjectID) error
}

type userRepository struct {
	client     *db.MongoDBClient
	collection *mongo.Collection
}

func NewUserRepository(client *db.MongoDBClient) UserRepository {
	collection := client.GetCollection("users")
	return &userRepository{client: client, collection: collection}
}
