package repository

import (
	"context"
	"fmt"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	mongoClient *mongo.Collection
}

func NewChatRepository(mongo *mongo.Collection) *ChatRepository {
	return &ChatRepository{
		mongoClient: mongo,
	}
}

func (r *ChatRepository) CreateChatMessage(message *models.Message) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()
	result, err := r.mongoClient.InsertOne(ctx, message)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func (r *ChatRepository) GetMessages(id string) ([]models.Message, error) {
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetLimit(10)

	Id := bson.M{"id": id}

	result, err := r.mongoClient.Find(ctx, Id, findOptions)
	if err != nil {
		return nil, err
	}

	var messages []models.Message
	if err = result.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
