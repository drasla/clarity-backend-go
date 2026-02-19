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
		if closeErr := closeMainDB(mainDB); closeErr != nil {
			log.Printf("⚠️ Failed to close MainDB during rollback: %v", closeErr)
		}

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

func closeMainDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB from gorm: %w", err)
	}
	return sqlDB.Close()
}

func (c *Container) Close() error {
	var errs []error

	if c.MainDB != nil {
		if err := closeMainDB(c.MainDB); err != nil {
			log.Printf("⚠️ Error closing MainDB: %v", err)
			errs = append(errs, err)
		}
	}
	if c.SmsDB != nil {
		if err := c.SmsDB.Close(); err != nil {
			log.Printf("⚠️ Error closing SmsDB: %v", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close databases: %v", errs)
	}
	return nil
}
