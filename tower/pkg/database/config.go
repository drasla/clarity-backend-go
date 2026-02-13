package database

import (
	"fmt"
	"tower/pkg/env"
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
		Host:     env.GetString("MYSQL84_HOST", "127.0.0.1"),
		Port:     env.GetString("MYSQL84_PORT", "8054"),
		User:     env.GetString("MYSQL84_USERNAME", "root"),
		Password: env.GetString("MYSQL84_PASSWORD", ""),
		DBName:   env.GetString("MYSQL84_DATABASE", ""),
		Charset:  "utf8mb4",
	}

	smsConfig := ConnectionInfo{
		Host:     env.GetString("MYSQL51_HOST", "127.0.0.1"),
		Port:     env.GetString("MYSQL51_PORT", "8051"),
		User:     env.GetString("MYSQL51_USERNAME", "root"),
		Password: env.GetString("MYSQL51_PASSWORD", ""),
		DBName:   env.GetString("MYSQL51_DATABASE", ""),
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
