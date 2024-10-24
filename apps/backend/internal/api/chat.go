package api

import (
	"encoding/json"
	"net/http"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/repository"
	"github.com/btk-hackathon-24-debug-duo/project-setup/pkg/utils"
)

func (h *Handlers) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
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

	result, err := chatRepo.CreateChatMessage(message)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Cannot add massage")
		return
	}

	utils.JSONResponse(w, http.StatusOK, result)

}
