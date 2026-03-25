package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Service ServiceConfig `envPrefix:"SERVICE_"`

	Providers    []string     `env:"PROVIDERS" envSeparator:","`
	GarudaConfig GarudaConfig `envPrefix:"GARUDA_"`
}

type GarudaConfig struct {
	MockPath string `env:"MOCK_PATH"`
}

type ServiceConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: No .env file found, relying on OS environment")
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	return &cfg
}
