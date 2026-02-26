package fnEcho

import (
	"context"
	"log"
	"net/http"
	"time"
)

func Shutdown(srv *http.Server) {
	log.Println("[fnEcho] 🛑 Shutting down Echo server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[fnEcho] ❌ Server forced to shutdown: %v", err)
	} else {
		log.Println("[fnEcho] ✅ Echo server shutdown completed")
	}
}
