package config

import (
	"log"
	"tower/pkg/fnEnv"
)

type EchoConfig struct {
	Port      string
	JwtSecret string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type S3Config struct {
	Domain     string
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
}

type MailgunConfig struct {
	Domain string
	APIKey string
	Sender string
}

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

type AppConfig struct {
	Env string

	Echo     EchoConfig
	MySQL84  DBConfig
	MySQL51  DBConfig
	S3       S3Config
	Mailgun  MailgunConfig
	Telegram TelegramConfig
}

var App *AppConfig

func LoadEnv() {
	fnEnv.Load()

	App = &AppConfig{
		Env: fnEnv.GetString("ENV", "development"),
		Echo: EchoConfig{
			Port:      fnEnv.GetString("PORT", "8080"),
			JwtSecret: fnEnv.MustGetString("JWT_SECRET"),
		},
		MySQL84: DBConfig{
			Host:     fnEnv.MustGetString("MYSQL84_HOST"),
			Port:     fnEnv.MustGetString("MYSQL84_PORT"),
			Username: fnEnv.MustGetString("MYSQL84_USERNAME"),
			Password: fnEnv.MustGetString("MYSQL84_PASSWORD"),
			Database: fnEnv.MustGetString("MYSQL84_DATABASE"),
		},
		MySQL51: DBConfig{
			Host:     fnEnv.MustGetString("MYSQL51_HOST"),
			Port:     fnEnv.MustGetString("MYSQL51_PORT"),
			Username: fnEnv.MustGetString("MYSQL51_USERNAME"),
			Password: fnEnv.MustGetString("MYSQL51_PASSWORD"),
			Database: fnEnv.MustGetString("MYSQL51_DATABASE"),
		},
		S3: S3Config{
			Domain:     fnEnv.MustGetString("S3_DOMAIN"),
			Endpoint:   fnEnv.MustGetString("S3_ENDPOINT"),
			AccessKey:  fnEnv.MustGetString("S3_ACCESS_KEY"),
			SecretKey:  fnEnv.MustGetString("S3_SECRET_KEY"),
			BucketName: fnEnv.MustGetString("S3_BUCKET_NAME"),
		},
		Mailgun: MailgunConfig{
			Domain: fnEnv.MustGetString("MAILGUN_DOMAIN"),
			APIKey: fnEnv.MustGetString("MAILGUN_API_KEY"),
			Sender: fnEnv.MustGetString("MAIL_SENDER_ADDRESS"),
		},
		Telegram: TelegramConfig{
			BotToken: fnEnv.GetString("TELEGRAM_BOT_TOKEN", ""),
			ChatID:   fnEnv.GetString("TELEGRAM_CHAT_ID", ""),
		},
	}

	log.Println("[Config] 프로젝트 환경변수 매핑 및 무결성 검증 완료 🚀")
}
