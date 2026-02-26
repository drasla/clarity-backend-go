package fnMySQL

import (
	"database/sql"
	"log"

	"gorm.io/gorm"
)

func CloseGORM(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("[fnMySQL] ⚠️ Failed to get sql.DB from GORM: %v", err)
		return err
	}
	return sqlDB.Close()
}

func CloseSQL(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
