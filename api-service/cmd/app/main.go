package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"api-service/internal/client"
	"api-service/internal/config"
	"api-service/internal/handler"
	"api-service/internal/middleware"
	"api-service/internal/repository"
	"api-service/internal/service"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("load config", zap.Error(err))
	}

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("connect to db", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		logger.Fatal("ping db", zap.Error(err))
	}

	userRepo := repository.NewUserRepo(db)
	cityRepo := repository.NewCityRepo(db)
	historyRepo := repository.NewHistoryRepo(db)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	gatewayClient := client.NewGatewayClient(httpClient, cfg.GatewayURL)
	weatherService := service.NewWeatherService(gatewayClient)
	userService := service.NewUserService(userRepo, cityRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userWeatherService := service.NewUserWeatherService(userRepo, cityRepo, historyRepo, weatherService)

	weatherHandler := handler.NewWeatherHandler(weatherService)
	userHandler := handler.NewUserHandler(userService, userWeatherService)
	authHandler := handler.NewAuthHandler(authService)

	auth := middleware.Auth(cfg.JWTSecret)
	adminOnly := middleware.RequireRole("admin")

	router := chi.NewRouter()
	router.Use(middleware.Logging(logger))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	router.Post("/auth/register", authHandler.Register)
	router.Post("/auth/login", authHandler.Login)

	router.Get("/weather", weatherHandler.GetWeather)
	router.Get("/weather/{city}", weatherHandler.GetCityWeather)
	router.Get("/weather/country/{country}", weatherHandler.GetCountryWeather)
	router.Get("/weather/country/{country}/top", weatherHandler.GetCountryWeatherTop)

	router.Group(func(r chi.Router) {
		r.Use(auth)

		r.Get("/users/me", userHandler.GetMe)

		r.Post("/cities", userHandler.AddCity)
		r.Get("/cities", userHandler.GetCities)
		r.Delete("/cities/{city_id}", userHandler.DeleteCity)

		r.Get("/weather/history", userHandler.GetWeatherHistory)
		r.Get("/users/weather", userHandler.GetUserWeather)
	})

	router.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOnly)

		r.Get("/users", userHandler.GetUsers)
		r.Get("/users/{id}", userHandler.GetUser)
		r.Delete("/users/{id}", userHandler.DeleteUser)
	})

	logger.Info("server started", zap.String("addr", cfg.Addr))
	if err := http.ListenAndServe(cfg.Addr, router); err != nil {
		logger.Fatal("listen", zap.Error(err))
	}
}
