package main

import (
	"log"
	"tower/pkg/database"
	"tower/pkg/fnEcho"
	"tower/pkg/fnEnv"
	"tower/pkg/fnError"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Info: No .fnEnv file found, relying on environment variables")
	}
	fnEnv.Load()

	db := database.MustInit()

	errHandler := fnError.NewErrorHandler(db.MainDB)

	srv := fnEcho.StartEchoServer(db, errHandler)

	fnEcho.WaitForShutdown(srv, db)
}
