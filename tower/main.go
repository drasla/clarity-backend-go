package main

import (
	"log"
	"tower/pkg/database"
	"tower/pkg/handler"
	"tower/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Info: No .env file found, relying on environment variables")
	}

	db := database.MustInit()
	defer db.Close()

	errHandler := handler.NewErrorHandler(db.MainDB)

	server.StartEchoServer(db, errHandler)
}
