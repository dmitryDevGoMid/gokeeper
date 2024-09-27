package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *userRepository) GetByKey(ctx context.Context, userRequest *User) (*User, error) {
	filter := bson.M{"publickey": userRequest.PublicKey}
	var user User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrDataNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	filter := bson.M{"username": username}
	var user User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrDataNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetCardsByUser(ctx context.Context, user *User) (*[]ResponseSaveData, error) {
	filter := bson.M{
		"user_id":   user.ID_User,
		"type_data": "card",
	}

	// Создаем конвейер для агрегации
	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$project": bson.M{
				"data":   1,
				"id_str": "$_id",
			},
		},
	}

	// Выполняем агрегацию
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ResponseSaveData
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func (r *userRepository) GetPasswordByUser(ctx context.Context, user *User) (*[]ResponseSaveData, error) {
	filter := bson.M{
		"user_id":   user.ID_User,
		"type_data": "password",
	}

	// Создаем конвейер для агрегации
	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$project": bson.M{
				"data":   1,
				"id_str": "$_id",
			},
		},
	}

	// Выполняем агрегацию
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ResponseSaveData
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}
