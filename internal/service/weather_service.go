package service

import (
	"context"
	"fmt"
)

type WeatherProvider interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*ProviderWeatherResponse, error)
	GetCurrentCityWeather(ctx context.Context, city string) (*ProviderWeatherResponse, error)
	GetCurrentCountryWeather(ctx context.Context, country string) ([]*ProviderWeatherResponse, error)
	GetCurrentCountryWeatherTop(ctx context.Context, country string) ([]*ProviderWeatherResponse, error)
}

type ProviderWeatherResponse struct {
	City        string
	Latitude    float64
	Longitude   float64
	Temperature float64
	WindSpeed   float64
	WeatherCode int
	Time        string
}

type WeatherResult struct {
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
	Description string  `json:"description"`
	Clothing    string  `json:"clothing,omitempty"`
}

type WeatherService struct {
	provider WeatherProvider
}

func NewWeatherService(provider WeatherProvider) *WeatherService {
	return &WeatherService{
		provider: provider,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, lat, lon float64) (*WeatherResult, error) {
	resp, err := s.provider.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("get weather from provider: %w", err)
	}

	return &WeatherResult{
		Latitude:    lat,
		Longitude:   lon,
		Temperature: resp.Temperature,
		WindSpeed:   resp.WindSpeed,
		WeatherCode: resp.WeatherCode,
		Time:        resp.Time,
		Description: mapWeatherCode(resp.WeatherCode),
	}, nil
}

func (s *WeatherService) GetCityWeather(ctx context.Context, city string) (*WeatherResult, error) {
	resp, err := s.provider.GetCurrentCityWeather(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("get city weather from provider: %w", err)
	}

	clothing := getClothingRecommendation(resp.Temperature)

	return &WeatherResult{
		City:        city,
		Latitude:    resp.Latitude,
		Longitude:   resp.Longitude,
		Temperature: resp.Temperature,
		WindSpeed:   resp.WindSpeed,
		WeatherCode: resp.WeatherCode,
		Time:        resp.Time,
		Description: mapWeatherCode(resp.WeatherCode),
		Clothing:    clothing,
	}, nil
}

func (s *WeatherService) GetCountryWeather(ctx context.Context, country string) ([]*WeatherResult, error) {
	resps, err := s.provider.GetCurrentCountryWeather(ctx, country)
	if err != nil {
		return nil, fmt.Errorf("get country weather from provider: %w", err)
	}

	results := make([]*WeatherResult, len(resps))
	for i, resp := range resps {
		results[i] = &WeatherResult{
			City:        resp.City,
			Latitude:    resp.Latitude,
			Longitude:   resp.Longitude,
			Temperature: resp.Temperature,
			WindSpeed:   resp.WindSpeed,
			WeatherCode: resp.WeatherCode,
			Time:        resp.Time,
			Description: mapWeatherCode(resp.WeatherCode),
		}
	}
	return results, nil
}

func (s *WeatherService) GetCountryWeatherTop(ctx context.Context, country string) ([]*WeatherResult, error) {
	resps, err := s.provider.GetCurrentCountryWeatherTop(ctx, country)
	if err != nil {
		return nil, fmt.Errorf("get country weather top from provider: %w", err)
	}

	results := make([]*WeatherResult, len(resps))
	for i, resp := range resps {
		results[i] = &WeatherResult{
			City:        resp.City,
			Latitude:    resp.Latitude,
			Longitude:   resp.Longitude,
			Temperature: resp.Temperature,
			WindSpeed:   resp.WindSpeed,
			WeatherCode: resp.WeatherCode,
			Time:        resp.Time,
			Description: mapWeatherCode(resp.WeatherCode),
		}
	}
	return results, nil
}
func getClothingRecommendation(temp float64) string {
	if temp < 5 {
		return "тёплая одежда"
	} else if temp < 15 {
		return "куртка"
	} else {
		return "лёгкая одежда"
	}
}

func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1, 2, 3:
		return "Переменная облачность"
	case 45, 48:
		return "Туман"
	case 51, 53, 55:
		return "Морось"
	case 61, 63, 65:
		return "Дождь"
	case 71, 73, 75:
		return "Снег"
	case 95:
		return "Гроза"
	default:
		return "Неизвестно"
	}
}
