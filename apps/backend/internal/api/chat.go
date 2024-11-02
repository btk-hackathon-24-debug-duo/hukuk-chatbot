package api

import (
	"context"
	"database/sql"
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
	db           *sql.DB
}

func NewChatHandlers(mongo *mongo.Collection, gemini *genai.GenerativeModel, db *sql.DB) *ChatHandlers {
	return &ChatHandlers{
		mongoClient:  mongo,
		geminiClient: gemini,
		db:           db,
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

	if payload.ChatId == "" || payload.Message == "" {
		utils.JSONError(w, http.StatusBadRequest, "Chat ID and message are required")
		return
	}

	message := &models.Message{
		ChatId:   payload.ChatId,
		UserId:   userID,
		Message:  payload.Message,
		Category: payload.Category,
	}

	chatRepo := repository.NewChatRepository(h.mongoClient, h.db)

	_, err := chatRepo.GetChat(message.ChatId, userID)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "This chat does not exists"+err.Error())
		return
	}

	_, err = chatRepo.CreateChatMessage(message)
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

	var data models.DataModel
	ansBytes, err := json.Marshal(result)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot marshal answer")
		return
	}
	err = json.Unmarshal(ansBytes, &data)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot unmarshal answer")
		return

	}

	answer := data.Candidates[0].Content.Parts[0]

	message = &models.Message{
		ChatId:   message.ChatId,
		UserId:   "AI",
		Message:  answer,
		Category: payload.Category,
	}

	_, err = chatRepo.CreateChatMessage(message)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot add massage")
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

	chatId := r.URL.Query().Get("chat_id")

	chatRepo := repository.NewChatRepository(h.mongoClient, h.db)

	chat, err := chatRepo.GetChat(chatId, userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot get messages"+err.Error())
		return
	}

	if chat.Id == "" {
		utils.JSONError(w, http.StatusBadRequest, "This chat does not exists")
		return
	}

	messages, err := chatRepo.GetMessages(chatId)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot get messages"+err.Error())
		return
	}

	type response struct {
		Prompt string `json:"prompt"`
		Answer string `json:"answer"`
	}

	var messagesResponse []response

	prompt := ""
	answer := ""
	i := 0
	for _, message := range messages {
		if message.UserId == userID {
			prompt = message.Message
			if i > 0 {
				messagesResponse = append(messagesResponse, response{Prompt: prompt, Answer: answer})
			}
			answer = ""
		} else {
			answer = message.Message
		}
		i++
	}

	utils.JSONResponse(w, http.StatusOK, messagesResponse)
}

func (h *ChatHandlers) SendFirstMessageHandler(w http.ResponseWriter, r *http.Request) {

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

	var payload models.SendFirstMessagePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if payload.Message == "" || payload.Name == "" {
		utils.JSONError(w, http.StatusBadRequest, "Chat ID, message and name are required")
		return
	}

	chatRepo := repository.NewChatRepository(h.mongoClient, h.db)
	chat_id, err := chatRepo.NewChat(userID, payload.Name)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot create chat")
		return
	}

	message := &models.Message{
		ChatId:   chat_id,
		UserId:   userID,
		Message:  payload.Message,
		Category: payload.Category,
	}

	_, err = chatRepo.CreateChatMessage(message)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot add massage")
		return
	}

	ctx := context.Background()
	ans, err := h.geminiClient.GenerateContent(ctx, genai.Text(message.Message))
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot answer")
		return
	}

	var data models.DataModel
	ansBytes, err := json.Marshal(ans)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot marshal answer")
		return
	}
	err = json.Unmarshal(ansBytes, &data)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot unmarshal answer")
		return

	}

	answer := data.Candidates[0].Content.Parts[0]

	message = &models.Message{
		ChatId:   chat_id,
		UserId:   "AI",
		Message:  answer,
		Category: payload.Category,
	}

	_, err = chatRepo.CreateChatMessage(message)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot add massage")
		return
	}

	utils.JSONResponse(w, http.StatusOK, ans)
}
func (h *ChatHandlers) GetChats(w http.ResponseWriter, r *http.Request) {
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

	chatRepo := repository.NewChatRepository(h.mongoClient, h.db)

	result, err := chatRepo.GetChats(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot get chats")
		return
	}

	utils.JSONResponse(w, http.StatusOK, result)
}

func (h *ChatHandlers) NewChat(w http.ResponseWriter, r *http.Request) {
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

	var payload struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chatRepo := repository.NewChatRepository(h.mongoClient, h.db)

	result, err := chatRepo.NewChat(userID, payload.Name)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot create chat")
		return
	}

	utils.JSONResponse(w, http.StatusOK, result)
}

func (h *ChatHandlers) UpdateChatNameHandler(w http.ResponseWriter, r *http.Request) {
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

	var chat models.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chatRepo := repository.NewChatRepository(h.mongoClient, h.db)

	err := chatRepo.UpdateChatName(chat.Id, userID, chat.Name)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot update chat name"+err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Chat name updated")
}
