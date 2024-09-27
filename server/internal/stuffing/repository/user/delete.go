package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *userRepository) DeleteByUsername(ctx context.Context, username string) error {
	filter := bson.M{"username": username}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *userRepository) DelerePasswordById(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{
		"_id":       id,
		"type_data": "password",
	}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with ID %s", id.Hex())
	}

	return nil
}

func (r *userRepository) DelereCardById(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{
		"_id":       id,
		"type_data": "card",
	}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with ID %s", id.Hex())
	}

	return nil
}
