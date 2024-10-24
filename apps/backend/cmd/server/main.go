package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/api"
	db "github.com/btk-hackathon-24-debug-duo/project-setup/pkg/database"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type app struct {
	db          *sql.DB
	mongoClient *mongo.Client
}

func main() {
	_ = godotenv.Load()

	port := ":8080"

	Db, err := db.GetDb()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mongoClient, err := db.SetupMongoDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := app{
		db:          Db,
		mongoClient: mongoClient,
	}

	router := api.NewRouter(app.db, app.mongoClient)

	r := router.NewRouter()

	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("user actions service started at " + port)

}
