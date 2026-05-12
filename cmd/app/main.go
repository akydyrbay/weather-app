package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"weather-api/internal/client"
	"weather-api/internal/config"
	"weather-api/internal/handler"
	"weather-api/internal/middleware"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	// repositories
	userRepo := repository.NewUserRepo(db)
	cityRepo := repository.NewCityRepo(db)
	historyRepo := repository.NewHistoryRepo(db)

	// services
	httpClient := &http.Client{Timeout: 10 * time.Second}
	weatherClient := client.NewWeatherClient(httpClient)
	weatherService := service.NewWeatherService(weatherClient)
	userService := service.NewUserService(userRepo, cityRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userWeatherService := service.NewUserWeatherService(userRepo, cityRepo, historyRepo, weatherService)

	// handlers
	weatherHandler := handler.NewWeatherHandler(weatherService)
	userHandler := handler.NewUserHandler(userService, userWeatherService)
	authHandler := handler.NewAuthHandler(authService)

	auth := middleware.Auth(cfg.JWTSecret)
	adminOnly := middleware.RequireRole("admin")

	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// public auth
	router.Post("/auth/register", authHandler.Register)
	router.Post("/auth/login", authHandler.Login)

	// public weather
	router.Get("/weather", weatherHandler.GetWeather)
	router.Get("/weather/{city}", weatherHandler.GetCityWeather)
	router.Get("/weather/country/{country}", weatherHandler.GetCountryWeather)
	router.Get("/weather/country/{country}/top", weatherHandler.GetCountryWeatherTop)

	// authenticated user routes
	router.Group(func(r chi.Router) {
		r.Use(auth)

		r.Get("/users/me", userHandler.GetMe)

		r.Post("/cities", userHandler.AddCity)
		r.Get("/cities", userHandler.GetCities)
		r.Delete("/cities/{city_id}", userHandler.DeleteCity)

		r.Get("/weather/history", userHandler.GetWeatherHistory)
		r.Get("/users/weather", userHandler.GetUserWeather)
	})

	// admin routes
	router.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOnly)

		r.Get("/users", userHandler.GetUsers)
		r.Get("/users/{id}", userHandler.GetUser)
		r.Delete("/users/{id}", userHandler.DeleteUser)
	})

	log.Printf("server started on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, router); err != nil {
		log.Fatal(err)
	}
}
