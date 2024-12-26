package config

import (
	"log"

	"github.com/joho/godotenv"
)

var DBDSN string

func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found, using environment variables")
	}
}

func LoadConfigGForMockDB() {
	// Load .env file from the root directory
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("No .env file found, using environment variables")
	}
}
