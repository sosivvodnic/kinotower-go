package core_database

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"true"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Database string `envconfig:"POSTGRES_DB" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("POSTGRES", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process envconfig: %w", err)
	}
	return &cfg, nil
}

// Исправлено название функции
func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("cannot create database config: %w", err))
	}
	return cfg
}
