package service

import (
	"context"
	"fmt"

	"weather-api/internal/model"
)

type UserWeatherService struct {
	users   UserRepository
	cities  CityRepository
	history HistoryRepository
	weather *WeatherService
}

func NewUserWeatherService(
	users UserRepository,
	cities CityRepository,
	history HistoryRepository,
	weather *WeatherService,
) *UserWeatherService {
	return &UserWeatherService{
		users:   users,
		cities:  cities,
		history: history,
		weather: weather,
	}
}

func (s *UserWeatherService) GetWeatherForUser(ctx context.Context, userID int) (*model.UserWeather, error) {
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	cities, err := s.cities.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user cities: %w", err)
	}

	if len(cities) == 0 {
		return &model.UserWeather{UserID: userID, Results: []*model.Weather{}}, nil
	}

	type fetchResult struct {
		weather *model.Weather
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

	results := make([]*model.Weather, 0, len(cities))
	for range cities {
		r := <-ch
		if r.err == nil {
			results = append(results, r.weather)
		}
	}

	for _, w := range results {
		_ = s.history.Save(ctx, userID, w.City, w.Temperature, w.Description)
	}

	return &model.UserWeather{UserID: userID, Results: results}, nil
}

func (s *UserWeatherService) GetHistory(ctx context.Context, userID int, filter model.HistoryFilter) (int, []*model.WeatherHistory, error) {
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return 0, nil, fmt.Errorf("user not found: %w", err)
	}

	history, err := s.history.Get(ctx, userID, filter)
	if err != nil {
		return 0, nil, err
	}

	return userID, history, nil
}
