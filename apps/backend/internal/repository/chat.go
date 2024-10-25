package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	mongoClient *mongo.Collection
	db          *sql.DB
}

func NewChatRepository(mongo *mongo.Collection, db *sql.DB) *ChatRepository {
	return &ChatRepository{
		mongoClient: mongo,
		db:          db,
	}
}

func (r *ChatRepository) CreateChatMessage(message *models.Message) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()
	message.AiModel = os.Getenv("AI_MODEL")
	result, err := r.mongoClient.InsertOne(ctx, message)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func (r *ChatRepository) GetMessages(chat_id, user_id string) ([]models.Message, error) {
	ctx := context.TODO()

	findOptions := options.Find()
	findOptions.SetLimit(10)

	Id := bson.M{"chatid": chat_id, "userid": user_id}

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

func (r *ChatRepository) GetChats(id string) ([]models.Chat, error) {
	stmt := `SELECT id, user_id, name FROM chats WHERE user_id = $1`

	//returning multiple records
	rows, err := r.db.Query(stmt, id)
	if err != nil {
		return []models.Chat{}, err
	}
	defer rows.Close()
	//rows to chat array
	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		err := rows.Scan(&chat.Id, &chat.UserId, &chat.Name)
		if err != nil {
			return []models.Chat{}, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *ChatRepository) NewChat(id, name string) (string, error) {
	chat_id := ""
	stmt := `INSERT INTO chats(user_id,name,created_at,updated_at) VALUES($1,$2,$3,$4) RETURNING id`
	err := r.db.QueryRow(stmt, id, name, time.Now(), time.Now()).Scan(&chat_id)
	if err != nil {
		return "", err
	}
	return chat_id, nil
}
