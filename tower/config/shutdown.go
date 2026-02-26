package config

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tower/pkg/fnEcho"
	"tower/pkg/fnMySQL"
)

func WaitForShutdown(srv *http.Server, db *ProjectDB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("========================================")
	log.Println("🛑 Graceful Shutdown initiated...")
	log.Println("========================================")

	fnEcho.Shutdown(srv)

	if db != nil {
		if err := fnMySQL.CloseGORM(db.MainDB); err != nil {
			log.Printf("⚠️ Error closing Main DB: %v", err)
		} else {
			log.Println("✅ Main DB (MySQL 8.4) closed properly")
		}

		if err := fnMySQL.CloseSQL(db.SmsDB); err != nil {
			log.Printf("⚠️ Error closing SMS DB: %v", err)
		} else {
			log.Println("✅ SMS DB (MySQL 5.1) closed properly")
		}
	}

	log.Println("👋 Application exited properly")
}
