package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Token    string             `bson:"token"`
}

type UserRepository interface {
	Create(context.Context, *User) error
	GetByUsername(context.Context, string) (*User, error)
	UpdatePassword(context.Context, string, string) error
	DeleteByUsername(context.Context, string) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{collection: collection}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	filter := bson.M{"username": username}
	var user User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, username, password string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"password": password}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *userRepository) DeleteByUsername(ctx context.Context, username string) error {
	filter := bson.M{"username": username}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
