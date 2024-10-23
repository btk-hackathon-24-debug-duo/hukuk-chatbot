package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/api"
	db "github.com/btk-hackathon-24-debug-duo/project-setup/pkg/database"
	"github.com/joho/godotenv"
)

type app struct {
	db *sql.DB
}

func main() {
	_ = godotenv.Load()

	port := ":8080"

	db, err := db.GetDb()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app := app{
		db: db,
	}

	router := api.NewRouter(app.db)

	r := router.NewRouter()

	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("user actions service started at " + port)

}
