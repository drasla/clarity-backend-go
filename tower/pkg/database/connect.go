package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Container struct {
	MainDB *gorm.DB
	SmsDB  *sql.DB
}

func NewContainer(cfg *Config) (*Container, error) {
	mainDB, err := connectMainDB(cfg.Main.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to main db: %w", err)
	}

	smsDB, err := connectSmsDB(cfg.Sms.DSN())
	if err != nil {
		closeMainDB(mainDB)
		return nil, fmt.Errorf("failed to connect to sms db: %w", err)
	}

	return &Container{
		MainDB: mainDB,
		SmsDB:  smsDB,
	}, nil
}

func connectMainDB(dsn string) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Connected to Main DB (MySQL 8.4)")
	return db, nil
}

func connectSmsDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("✅ Connected to SMS DB (MySQL 5.1)")
	return db, nil
}

func (c *Container) Close() {
	if c.MainDB != nil {
		closeMainDB(c.MainDB)
	}
	if c.SmsDB != nil {
		if err := c.SmsDB.Close(); err != nil {
			log.Printf("⚠️ Error closing SmsDB: %v", err)
		}
	}
}

func closeMainDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("⚠️ Failed to retrieve SQL DB from GORM to close: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("⚠️ Error closing MainDB: %v", err)
	}
}
