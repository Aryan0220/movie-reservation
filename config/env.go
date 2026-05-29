package config

import (
	"os"
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() string {
	err := godotenv.Load()
	if err != nil {
		return "No .env file Found."
	}
	return ".Env File Found"
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func PrintLog(message string, kind string) {
	if GetEnv("APP_ENV") == "development" {
		logMessage := "[" + kind + "] " + message
		log.Println(logMessage)
	}
}
