package core_server

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type AuthConfig struct {
	JWTSecret string `envconfig:"JWT_SECRET" required:"true"`
}

func NewAuthConfig() (*AuthConfig, error) {
	var cfg AuthConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("cannot process auth config: %w", err)
	}
	return &cfg, nil
}

func NewAuthConfigMust() *AuthConfig {
	cfg, err := NewAuthConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

