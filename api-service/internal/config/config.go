package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Addr        string
	GatewayURL  string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL env var is required")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET env var is required")
	}
	gatewayURL := os.Getenv("GATEWAY_URL")
	if gatewayURL == "" {
		return nil, fmt.Errorf("GATEWAY_URL env var is required")
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return &Config{DatabaseURL: dbURL, JWTSecret: secret, Addr: addr, GatewayURL: gatewayURL}, nil
}
