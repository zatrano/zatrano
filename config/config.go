package config

import (
	"log"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default environment variables.")
	}
}

func Get(key string, defaultValue ...string) string {
	val := os.Getenv(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func GetInt(key string, defaultValue ...int) int {
	valStr := Get(key)
	if valStr == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	val, err := strconv.Atoi(valStr)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}