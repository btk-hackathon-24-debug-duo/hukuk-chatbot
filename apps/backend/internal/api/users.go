package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/repository"
	"github.com/btk-hackathon-24-debug-duo/project-setup/pkg/utils"
	"github.com/golang-jwt/jwt"
)

type Handlers struct {
	db *sql.DB
}

type Claims struct {
	User models.User `json:"user"`
	jwt.StandardClaims
}

var jwtKey = []byte("your_secret_key")

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{
		db: db,
	}
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	usersRepo := repository.NewUsersRepository(h.db)

	var user models.User

	user.Email = r.URL.Query().Get("email")
	user.Password = utils.HashPassword(r.URL.Query().Get(("password")))

	exists, err := usersRepo.GetUserWithEmailPassword(user)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		utils.JSONError(w, http.StatusUnauthorized, "User is not exists")
		return
	}

	utils.JSONResponse(w, http.StatusOK, exists)
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	usersRepo := repository.NewUsersRepository(h.db)

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := usersRepo.CreateUser(user)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		User: result,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, tokenString)
}
