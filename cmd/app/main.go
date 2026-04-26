package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"weather-api/internal/client"
	"weather-api/internal/handler"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL env var is required")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	// repo
	userRepo := repository.NewUserRepo(db)
	cityRepo := repository.NewCityRepo(db)
	historyRepo := repository.NewHistoryRepo(db)

	// weather client and service
	httpClient := &http.Client{Timeout: 10 * time.Second}
	weatherClient := client.NewWeatherClient(httpClient)
	weatherService := service.NewWeatherService(weatherClient)

	// user services
	userService := service.NewUserService(userRepo)
	userWeatherService := service.NewUserWeatherService(userRepo, cityRepo, historyRepo, weatherService)

	// handlers
	weatherHandler := handler.NewWeatherHandler(weatherService)
	userHandler := handler.NewUserHandler(userService, userWeatherService, cityRepo)

	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// weather
	router.Get("/weather", weatherHandler.GetWeather)
	router.Get("/weather/{city}", weatherHandler.GetCityWeather)
	router.Get("/weather/country/{country}", weatherHandler.GetCountryWeather)
	router.Get("/weather/country/{country}/top", weatherHandler.GetCountryWeatherTop)

	// users
	router.Post("/users", userHandler.CreateUser)
	router.Get("/users", userHandler.GetUsers)
	router.Get("/users/{id}", userHandler.GetUser)
	router.Put("/users/{id}", userHandler.UpdateUser)
	router.Delete("/users/{id}", userHandler.DeleteUser)

	// user cities
	router.Post("/users/{id}/cities", userHandler.AddCity)
	router.Get("/users/{id}/cities", userHandler.GetCities)
	router.Delete("/users/{id}/cities/{city_id}", userHandler.DeleteCity)

	// user weather & history
	router.Get("/users/{id}/weather", userHandler.GetUserWeather)
	router.Get("/users/{id}/weather/history", userHandler.GetWeatherHistory)

	addr := ":8080"
	log.Printf("server started on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
