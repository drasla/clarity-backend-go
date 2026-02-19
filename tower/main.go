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

	errHandler := handler.NewErrorHandler(db.MainDB)

	srv := fnEcho.StartEchoServer(db, errHandler)

	fnEcho.WaitForShutdown(srv, db)
}
