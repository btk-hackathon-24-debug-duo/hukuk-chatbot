package repository

import (
	"database/sql"
	"time"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"github.com/btk-hackathon-24-debug-duo/project-setup/pkg/utils"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (u *UsersRepository) CreateUser(user models.User) (models.User, error) {

	stmt := `INSERT INTO users (first_name, last_name, email, password, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING id, first_name, last_name, email;`
	var User models.User
	err := u.db.QueryRow(stmt, user.FirstName, user.LastName, user.Email, utils.HashPassword(user.Password), time.Now(), time.Now()).Scan(&User.Id, &User.FirstName, &User.LastName, &User.Email)
	if err != nil {
		return models.User{}, err
	}

	return User, nil
}

func (u *UsersRepository) GetUserWithEmailPassword(user models.User) (bool, error) {

	var exists bool

	stmt := `SELECT 1 FROM users WHERE email=$1 AND password=$2`
	err := u.db.QueryRow(stmt, user.Email, user.Password).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
