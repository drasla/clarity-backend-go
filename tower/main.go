package main

import (
	"log"
	"tower/pkg/database"
	"tower/pkg/fnEcho"
	"tower/pkg/handler"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Info: No .fnEnv file found, relying on environment variables")
	}

	db := database.MustInit()
	defer db.Close()

	errHandler := handler.NewErrorHandler(db.MainDB)

	fnEcho.StartEchoServer(db, errHandler)
}
