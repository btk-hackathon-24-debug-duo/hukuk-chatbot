package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/repository"
	"github.com/btk-hackathon-24-debug-duo/project-setup/pkg/utils"
	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatHandlers struct {
	mongoClient  *mongo.Collection
	geminiClient *genai.GenerativeModel
}

func NewChatHandlers(mongo *mongo.Collection, gemini *genai.GenerativeModel) *ChatHandlers {
	return &ChatHandlers{
		mongoClient:  mongo,
		geminiClient: gemini,
	}
}

func (h *ChatHandlers) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetTokenClaims(r)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "Token claims missing")
		return
	}

	userID, ok := utils.GetUserIDFromClaims(claims)
	if !ok {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var payload models.SendMessagePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := &models.Message{
		Id:       userID,
		Message:  payload.Message,
		Category: payload.Category,
	}

	chatRepo := repository.NewChatRepository(h.mongoClient)

	_, err := chatRepo.CreateChatMessage(message)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot add massage")
		return
	}

	ctx := context.Background()
	result, err := h.geminiClient.GenerateContent(ctx, genai.Text(message.Message))
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot answer")
		return
	}

	utils.JSONResponse(w, http.StatusOK, result)

}

func (h *ChatHandlers) GetMessages(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetTokenClaims(r)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "Token claims missing")
		return
	}

	userID, ok := utils.GetUserIDFromClaims(claims)
	if !ok {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	chatRepo := repository.NewChatRepository(h.mongoClient)

	messages, err := chatRepo.GetMessages(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot get messages"+err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, messages)
}
