package fnEnv

import (
	"log"
	"os"
	"strconv"
	"sync"
)

type AppConfig struct {
	Port string
	Env  string

	MySQL84Host     string
	MySQL84Port     string
	MySQL84Username string
	MySQL84Password string
	MySQL84Database string

	MySQL51Host     string
	MySQL51Port     string
	MySQL51Username string
	MySQL51Password string
	MySQL51Database string

	JwtSecret string

	S3Domain     string
	S3Endpoint   string
	S3AccessKey  string
	S3SecretKey  string
	S3BucketName string

	TelegramBotToken string
	TelegramChatID   string

	MailgunDomain string
	MailgunAPIKey string
	MailSender    string
}

var (
	App  *AppConfig
	once sync.Once
)

func Load() {
	once.Do(func() {
		App = &AppConfig{
			Port: getString("PORT", "8080"),
			Env:  getString("ENV", "development"),

			TelegramBotToken: getString("TELEGRAM_BOT_TOKEN", ""),
			TelegramChatID:   getString("TELEGRAM_CHAT_ID", ""),

			MySQL84Host:     mustGetString("MYSQL84_HOST"),
			MySQL84Port:     mustGetString("MYSQL84_PORT"),
			MySQL84Username: mustGetString("MYSQL84_USERNAME"),
			MySQL84Password: mustGetString("MYSQL84_PASSWORD"),
			MySQL84Database: mustGetString("MYSQL84_DATABASE"),

			MySQL51Host:     mustGetString("MYSQL51_HOST"),
			MySQL51Port:     mustGetString("MYSQL51_PORT"),
			MySQL51Username: mustGetString("MYSQL51_USERNAME"),
			MySQL51Password: mustGetString("MYSQL51_PASSWORD"),
			MySQL51Database: mustGetString("MYSQL51_DATABASE"),

			JwtSecret: mustGetString("JWT_SECRET"),

			S3Domain:     mustGetString("S3_DOMAIN"),
			S3Endpoint:   mustGetString("S3_ENDPOINT"),
			S3AccessKey:  mustGetString("S3_ACCESS_KEY"),
			S3SecretKey:  mustGetString("S3_SECRET_KEY"),
			S3BucketName: mustGetString("S3_BUCKET_NAME"),

			MailgunDomain: mustGetString("MAILGUN_DOMAIN"),
			MailgunAPIKey: mustGetString("MAILGUN_API_KEY"),
			MailSender:    mustGetString("MAIL_SENDER_ADDRESS"),
		}

		log.Println("[fnEnv] 모든 환경변수 로드 및 무결성 검증 완료 🚀")
	})
}

func mustGetString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("[fnEnv] 🚨 치명적 오류: 필수 환경변수 '%s'가 설정되지 않았습니다.", key)
	}
	return val
}

func getString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return value
}

func getBool(key string, fallback bool) bool {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return fallback
	}
	return value
}
