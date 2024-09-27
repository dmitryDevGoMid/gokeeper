package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *userRepository) Create(ctx context.Context, user *User) error {
	fmt.Println("Creating===>>>", user)
	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) CreatePasswordByUser(ctx context.Context, user *SaveData) error {
	fmt.Println("Creating===>>>", user)
	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) CreateCardByUser(ctx context.Context, user *SaveData) error {
	fmt.Println("Creating Card===>>>", user)
	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}
