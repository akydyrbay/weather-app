package service

import (
	"context"
	"fmt"

	"weather-api/internal/repository"
)

type UserWeatherService struct {
	users   *repository.UserRepo
	cities  *repository.CityRepo
	history *repository.HistoryRepo
	weather *WeatherService
}

func NewUserWeatherService(
	users *repository.UserRepo,
	cities *repository.CityRepo,
	history *repository.HistoryRepo,
	weather *WeatherService,
) *UserWeatherService {
	return &UserWeatherService{
		users:   users,
		cities:  cities,
		history: history,
		weather: weather,
	}
}

type UserWeatherResult struct {
	UserID  int              `json:"user_id"`
	Results []*WeatherResult `json:"results"`
}

func (s *UserWeatherService) GetWeatherForUser(ctx context.Context, userID int) (*UserWeatherResult, error) {
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	cities, err := s.cities.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user cities: %w", err)
	}

	if len(cities) == 0 {
		return &UserWeatherResult{UserID: userID, Results: []*WeatherResult{}}, nil
	}

	type fetchResult struct {
		weather *WeatherResult
		err     error
	}

	ch := make(chan fetchResult, len(cities))
	for _, c := range cities {
		city := c.City
		go func() {
			w, err := s.weather.GetCityWeather(ctx, city)
			ch <- fetchResult{weather: w, err: err}
		}()
	}

	results := make([]*WeatherResult, 0, len(cities))
	for range cities {
		r := <-ch
		if r.err == nil {
			results = append(results, r.weather)
		}
	}

	for _, w := range results {
		_ = s.history.Save(ctx, userID, w.City, w.Temperature, w.Description)
	}

	return &UserWeatherResult{UserID: userID, Results: results}, nil
}

type HistoryResponse struct {
	UserID  int                          `json:"user_id"`
	City    string                       `json:"city,omitempty"`
	History []*repository.WeatherHistory `json:"history"`
}

func (s *UserWeatherService) GetHistory(ctx context.Context, userID int, filter repository.HistoryFilter) (*HistoryResponse, error) {
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	history, err := s.history.Get(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	return &HistoryResponse{
		UserID:  userID,
		City:    filter.City,
		History: history,
	}, nil
}
