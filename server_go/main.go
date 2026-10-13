package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"todo-list-api/db"
	"todo-list-api/features/todo"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	loadEnv()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./todo.db"
	}

	store := &db.Store{}
	if err := store.Connect(dbPath); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()

	svc := todo.NewService(store.DB())
	h := todo.NewHandler(svc)
	r := todo.NewRouter(h)

	log.Println("Starting server on :8787")

	if err := http.ListenAndServe(":8787", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
