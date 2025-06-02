package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WorkingDirectory string
	Port             string
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config = Config{
		WorkingDirectory: os.Getenv("WORKING_DIR"),
		Port:             os.Getenv("PORT"),
	}
}
