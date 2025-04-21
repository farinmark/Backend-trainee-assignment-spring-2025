package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBUrl     string `envconfig:"DB_URL" required:"true"`
	JWTSecret string `envconfig:"JWT_SECRET" required:"true"`
}

func Load() *Config {
	var cfg Config
	if err := envconfig.Process("pvz", &cfg); err != nil {
		log.Fatalf("config load error: %v", err)
	}
	return &cfg
}
