package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string `yaml:"address" env-default:"localhost:8082"`
	Timeout     string `yaml:"timeout" env-default:"5s"`
	IdleTimeout string `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH env var is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	return &cfg
}
