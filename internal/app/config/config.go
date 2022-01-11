package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Address     string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	DBURL       string `env:"DATABASE_DSN" envDefault:""`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		log.Printf("env prsing failed:+%v\n", err)
		return nil, err
	}

	return cfg, nil
}

func (c *Config) InitFlags() {
	flag.StringVar(&c.Address, "a", c.Address, "host to listen on")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "base url")
	flag.StringVar(&c.StoragePath, "f", c.StoragePath, "file storage path")
	flag.StringVar(&c.DBURL, "d", c.DBURL, "data base url")
	flag.Parse()
}
