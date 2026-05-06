package config

import (
	"os"
	"github.com/joho/godotenv"
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
