package database

import "log"

func MustInit() *Container {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("❌ Failed to load database config: %v", err)
	}

	container, err := NewContainer(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to databases: %v", err)
	}

	return container
}
