package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/api"
	"github.com/btk-hackathon-24-debug-duo/project-setup/pkg/ai"
	db "github.com/btk-hackathon-24-debug-duo/project-setup/pkg/database"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type app struct {
	db           *sql.DB
	mongoClient  *mongo.Collection
	geminiClient *genai.GenerativeModel
}

func main() {
	_ = godotenv.Load()

	time.Sleep(10 * time.Second)
	port := ":8080"

	Db, err := db.SetupDb()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mongoClient, err := db.SetupMongoDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	geminiClient, err := ai.SetupGemini()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := app{
		db:           Db,
		mongoClient:  mongoClient,
		geminiClient: geminiClient,
	}

	router := api.NewRouter(app.db, app.mongoClient, app.geminiClient)

	r := router.NewRouter()

	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("user actions service started at " + port)

}
