package api

import (
	"database/sql"
	"net/http"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	db          *sql.DB
	mongoClient *mongo.Client
}

func NewRouter(db *sql.DB, mongo *mongo.Client) *Router {
	return &Router{
		db:          db,
		mongoClient: mongo,
	}
}

func (r *Router) NewRouter() *mux.Router {
	h := NewHandlers(r.db, r.mongoClient)

	router := mux.NewRouter()

	router.Use(middleware.CorsMiddleware)

	router.HandleFunc("/api/login", h.LoginHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/register", h.RegisterHandler).Methods(http.MethodPost, http.MethodOptions)

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.EnsureValidToken)

	protected.HandleFunc("/chat/message", h.SendMessageHandler).Methods(http.MethodPost, http.MethodOptions)

	return router
}
