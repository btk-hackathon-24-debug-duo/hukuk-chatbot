package repository

import (
	"context"
	"fmt"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository struct {
	mongoClient *mongo.Client
}

func NewChatRepository(mongo *mongo.Client) *ChatRepository {
	return &ChatRepository{
		mongoClient: mongo,
	}
}

func (r *ChatRepository) CreateChatMessage(message *models.Message) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()
	result, err := r.mongoClient.Database("chat").Collection("messages").InsertOne(ctx, message)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}
