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
	"weather-api/internal/middleware"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL env var is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env var is required")
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

	// services
	httpClient := &http.Client{Timeout: 10 * time.Second}
	weatherClient := client.NewWeatherClient(httpClient)
	weatherService := service.NewWeatherService(weatherClient)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtSecret)
	userWeatherService := service.NewUserWeatherService(userRepo, cityRepo, historyRepo, weatherService)

	// handlers
	weatherHandler := handler.NewWeatherHandler(weatherService)
	userHandler := handler.NewUserHandler(userService, userWeatherService, cityRepo)
	authHandler := handler.NewAuthHandler(authService)

	auth := middleware.Auth(jwtSecret)
	adminOnly := middleware.RequireRole("admin")

	router := chi.NewRouter()

	// public
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
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

		r.Get("/users/weather", userHandler.GetUserWeather)
		r.Get("/users/weather/history", userHandler.GetWeatherHistory)
	})

	// only admin routes
	router.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOnly)

		r.Get("/users", userHandler.GetUsers)
		r.Get("/users/{id}", userHandler.GetUser)
		r.Delete("/users/{id}", userHandler.DeleteUser)
	})

	addr := ":8080"
	log.Printf("server started on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
