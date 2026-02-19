package fnEcho

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tower/pkg/database"
)

func WaitForShutdown(srv *http.Server, db *database.Container) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("👋 Server exited properly")
}
