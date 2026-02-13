package database

import (
	"log"
	"tower/model/maindb"
)

func MustInit() *Container {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("❌ Failed to load database config: %v", err)
	}

	container, err := NewContainer(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to databases: %v", err)
	}

	err = container.MainDB.AutoMigrate(
		&maindb.User{},
		&maindb.RefreshToken{},
		&maindb.Verification{},
	)
	if err != nil {
		log.Fatalf("❌ Migration Failed: %v", err)
	}

	return container
}
