package main

import (
	"tower/config"
	"tower/pkg/fnError"
)

func main() {
	config.LoadEnv()

	db := config.InitDatabase()
	defer db.Close()

	execSchema := config.NewExecutableSchema(db)
	errHandler := fnError.NewErrorHandler(db.MainDB)
	srv := config.StartWebServer(errHandler, execSchema)

	config.WaitForShutdown(srv, db)
}
