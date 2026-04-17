package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"weather-api/internal/client"
	"weather-api/internal/handler"
	"weather-api/internal/service"
)

func main() {
	router := chi.NewRouter()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	weatherClient := client.NewWeatherClient(httpClient)
	weatherService := service.NewWeatherService(weatherClient)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	router.Get("/weather", weatherHandler.GetWeather)
	router.Get("/weather/{city}", weatherHandler.GetCityWeather)
	router.Get("/weather/country/{country}", weatherHandler.GetCountryWeather)
	router.Get("/weather/country/{country}/top", weatherHandler.GetCountryWeatherTop)

	addr := ":8080"
	log.Printf("server started on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
