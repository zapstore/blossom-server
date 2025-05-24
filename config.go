package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WorkingDirectory string `envconfig:"WORKING_DIR"`
	Port             string `envconfig:"PORT"`
}

func LoadConfig() {
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to read from env: %s", err)
		return
	}
}
