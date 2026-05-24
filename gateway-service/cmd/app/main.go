package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"gateway-service/internal/client"
	"gateway-service/internal/config"
	"gateway-service/internal/handler"
	"gateway-service/internal/middleware"
	"gateway-service/internal/service"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	cfg := config.Load()

	timeout, err := time.ParseDuration(cfg.RequestTimeout)
	if err != nil {
		logger.Fatal("parse REQUEST_TIMEOUT", zap.Error(err))
	}

	httpClient := &http.Client{Timeout: timeout}
	openMeteo := client.NewOpenMeteoClient(httpClient, cfg.OpenMeteoURL, cfg.GeocodingURL)
	weatherService := service.NewWeatherService(openMeteo)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	router := chi.NewRouter()
	router.Use(middleware.Logging(logger))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	router.Get("/weather", weatherHandler.ByCoords)
	router.Get("/weather/city/{city}", weatherHandler.ByCity)
	router.Get("/weather/country/{country}", weatherHandler.ByCountry)
	router.Get("/weather/country/{country}/top", weatherHandler.ByCountryTop)

	logger.Info("gateway started", zap.String("addr", cfg.Addr))
	if err := http.ListenAndServe(cfg.Addr, router); err != nil {
		logger.Fatal("listen", zap.Error(err))
	}
}
