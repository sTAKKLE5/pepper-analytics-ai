package utils

import (
	"log"
	"os"
	"strconv"
)

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable not set: %s", key)
	}
	return value
}

func GetRequiredEnvAsInt(key string) int {
	value := GetRequiredEnv(key)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Environment variable %s must be a valid integer: %v", key, err)
	}
	return intValue
}
