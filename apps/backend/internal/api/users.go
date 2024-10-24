package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/repository"
	"github.com/btk-hackathon-24-debug-duo/project-setup/pkg/utils"
)

type UserHandlers struct {
	db *sql.DB
}

func NewUserHandlers(db *sql.DB) *UserHandlers {
	return &UserHandlers{
		db: db,
	}
}

func (h *UserHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	usersRepo := repository.NewUsersRepository(h.db)

	var User models.User

	User.Email = r.URL.Query().Get("email")
	User.Password = utils.HashPassword(r.URL.Query().Get(("password")))

	result, err := usersRepo.GetUserWithEmailPassword(User)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if result.Id == "" {
		utils.JSONError(w, http.StatusUnauthorized, "User is not exists")
		return
	}

	tokenString, err := utils.CreateJWTToken(result)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, tokenString)
}

func (h *UserHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	usersRepo := repository.NewUsersRepository(h.db)

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := usersRepo.CreateUser(user)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokenString, err := utils.CreateJWTToken(result)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, tokenString)
}
