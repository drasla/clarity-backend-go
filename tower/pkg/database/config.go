package database

import (
	"fmt"
	"tower/pkg/fnEnv"
)

type ConnectionInfo struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

func (c ConnectionInfo) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}

type Config struct {
	Main ConnectionInfo
	Sms  ConnectionInfo
}

func LoadConfigFromEnv() (*Config, error) {
	mainConfig := ConnectionInfo{
		Host:     fnEnv.App.MySQL84Host,
		Port:     fnEnv.App.MySQL84Port,
		User:     fnEnv.App.MySQL84Username,
		Password: fnEnv.App.MySQL84Password,
		DBName:   fnEnv.App.MySQL84Database,
		Charset:  "utf8mb4",
	}

	smsConfig := ConnectionInfo{
		Host:     fnEnv.App.MySQL51Host,
		Port:     fnEnv.App.MySQL51Port,
		User:     fnEnv.App.MySQL51Username,
		Password: fnEnv.App.MySQL51Password,
		DBName:   fnEnv.App.MySQL51Database,
		Charset:  "utf8",
	}

	if mainConfig.User == "" || smsConfig.User == "" {
		return nil, fmt.Errorf("database user is required in environment variables")
	}

	return &Config{
		Main: mainConfig,
		Sms:  smsConfig,
	}, nil
}
