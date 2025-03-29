package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address           string        `yaml:"address"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if len(configPath) == 0 {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file not exist: %s", configPath)
	}

	cfg := Config{}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Dont parse configfile: %s", cfg)
	}

	return cfg
}
