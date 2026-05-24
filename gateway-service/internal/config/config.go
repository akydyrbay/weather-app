package config

import "os"

type Config struct {
	Addr           string
	OpenMeteoURL   string
	GeocodingURL   string
	RequestTimeout string
}

func Load() *Config {
	return &Config{
		Addr:           getEnv("ADDR", ":8081"),
		OpenMeteoURL:   getEnv("OPEN_METEO_URL", "https://api.open-meteo.com/v1/forecast"),
		GeocodingURL:   getEnv("GEOCODING_URL", "https://geocoding-api.open-meteo.com/v1/search"),
		RequestTimeout: getEnv("REQUEST_TIMEOUT", "10s"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
