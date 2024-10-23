package api

import (
	"database/sql"
	"net/http"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/middleware"
	"github.com/gorilla/mux"
)

type Router struct {
	db *sql.DB
}

func NewRouter(db *sql.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r *Router) NewRouter() *mux.Router {
	h := NewHandlers(r.db)

	router := mux.NewRouter()

	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.EnsureValidToken)

	router.HandleFunc("/api/login", h.LoginHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/register", h.RegisterHandler).Methods(http.MethodPost, http.MethodOptions)

	return router
}
