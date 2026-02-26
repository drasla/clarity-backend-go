package fnEnv

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
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

var once sync.Once

func Load() {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Printf("[fnEnv] Info: No .env file found, relying on OS environment variables")
		} else {
			log.Println("[fnEnv] .env file loaded successfully 📄")
		}
	})
}

func MustGetString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("[fnEnv] 🚨 치명적 오류: 필수 환경변수 '%s'가 설정되지 않았습니다.", key)
	}
	return val
}

func GetString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
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

func GetBool(key string, fallback bool) bool {
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
