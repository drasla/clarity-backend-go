package config

import (
	"database/sql"
	"log"
	"tower/model/maindb"
	"tower/pkg/fnMySQL"

	"gorm.io/gorm"
)

type ProjectDB struct {
	MainDB *gorm.DB
	SmsDB  *sql.DB
}

func (db *ProjectDB) Close() {
	if db.MainDB != nil {
		if sqlDB, err := db.MainDB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	if db.SmsDB != nil {
		_ = db.SmsDB.Close()
	}
}

func InitDatabase() *ProjectDB {
	mainCfg := &fnMySQL.Config{
		Host:     App.MySQL84.Host,
		Port:     App.MySQL84.Port,
		User:     App.MySQL84.Username,
		Password: App.MySQL84.Password,
		DBName:   App.MySQL84.Database,
		Charset:  "utf8mb4",
	}
	mainDB, err := fnMySQL.ConnectGORM(mainCfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MainDB: %v", err)
	}
	log.Println("[Registry] ✅ Connected to Main DB (MySQL 8.4)")

	smsCfg := &fnMySQL.Config{
		Host:     App.MySQL51.Host,
		Port:     App.MySQL51.Port,
		User:     App.MySQL51.Username,
		Password: App.MySQL51.Password,
		DBName:   App.MySQL51.Database,
		Charset:  "utf8",
	}
	smsDB, err := fnMySQL.ConnectSQL(smsCfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to SmsDB: %v", err)
	}
	log.Println("[Registry] ✅ Connected to SMS DB (MySQL 5.1)")

	err = mainDB.AutoMigrate(
		&maindb.User{},
		&maindb.RefreshToken{},
		&maindb.Verification{},
		&maindb.File{},
		&maindb.Inquiry{},
		&maindb.EmailTemplate{},
		&maindb.EmailLog{},
	)
	if err != nil {
		log.Fatalf("❌ Migration Failed: %v", err)
	}

	seedData(mainDB)

	return &ProjectDB{
		MainDB: mainDB,
		SmsDB:  smsDB,
	}
}
