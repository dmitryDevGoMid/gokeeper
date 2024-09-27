package user

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *userRepository) UpdateCardByKey(ctx context.Context, card *SaveData) error {

	// Получение текущего документа
	var currentData SaveData
	filter := bson.D{{Key: "_id", Value: card.ID}}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&currentData)
	if err != nil {
		log.Fatal(err)
	}

	// Обновление документа
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: card.User_ID},
			{Key: "type_data", Value: card.TypeData},
			{Key: "data", Value: card.Data},
		}},
	}

	updateResult, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("update card error:", err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}

func (r *userRepository) UpdatePasswordByKey(ctx context.Context, passwd *SaveData) error {

	// Получение текущего документа
	var currentData SaveData
	filter := bson.D{{Key: "_id", Value: passwd.ID}}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&currentData)
	if err != nil {
		log.Fatal(err)
	}

	// Обновление документа
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: passwd.User_ID},
			{Key: "type_data", Value: passwd.TypeData},
			{Key: "data", Value: passwd.Data},
		}},
	}

	updateResult, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}

func (r *userRepository) UpdateRegisterByKey(ctx context.Context, userRequest *User, password string) error {
	filter := bson.M{"key": userRequest.PublicKey}
	update := bson.M{"$set": bson.M{"password": userRequest.Password, "username": userRequest.Username}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, username, password string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"password": password}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
