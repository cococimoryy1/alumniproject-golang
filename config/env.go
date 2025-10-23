package config

import (
    "log"
	"os"
    "github.com/joho/godotenv"
)

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️  .env file not found or not loaded")
    } else {
        log.Println("✅ .env file loaded successfully")
    }
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
