package service

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"gateway-service/internal/client"
	"gateway-service/internal/dto"
)

const cacheTTL = 5 * time.Minute

type Provider interface {
	Forecast(ctx context.Context, lat, lon float64) (*client.ForecastResponse, error)
	Geocode(ctx context.Context, name string) (*client.GeoResult, error)
}

type cacheEntry struct {
	resp      *dto.WeatherResponse
	expiresAt time.Time
}

type WeatherService struct {
	provider Provider
	cache    sync.Map
}

func NewWeatherService(p Provider) *WeatherService {
	return &WeatherService{provider: p}
}

var citiesByCountry = map[string][]string{
	"Kazakhstan": {"Almaty", "Astana", "Shymkent", "Karaganda", "Aktobe", "Atyrau", "Pavlodar"},
	"Russia":     {"Moscow", "Saint Petersburg", "Novosibirsk", "Yekaterinburg", "Kazan", "Chelyabinsk", "Omsk"},
	"USA":        {"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio"},
	"Germany":    {"Berlin", "Hamburg", "Munich", "Cologne", "Frankfurt", "Stuttgart", "Dusseldorf"},
	"France":     {"Paris", "Marseille", "Lyon", "Toulouse", "Nice", "Nantes", "Strasbourg"},
	"China":      {"Beijing", "Shanghai", "Guangzhou", "Shenzhen", "Chengdu", "Wuhan", "Xian"},
	"Japan":      {"Tokyo", "Osaka", "Yokohama", "Nagoya", "Sapporo", "Fukuoka", "Kyoto"},
	"India":      {"Mumbai", "Delhi", "Bangalore", "Hyderabad", "Chennai", "Kolkata", "Pune"},
	"Brazil":     {"Sao Paulo", "Rio de Janeiro", "Brasilia", "Salvador", "Fortaleza", "Belo Horizonte", "Manaus"},
	"Canada":     {"Toronto", "Montreal", "Vancouver", "Calgary", "Edmonton", "Ottawa", "Winnipeg"},
}

func (s *WeatherService) ByCoords(ctx context.Context, lat, lon float64) (*dto.WeatherResponse, error) {
	f, err := s.provider.Forecast(ctx, lat, lon)
	if err != nil {
		return nil, err
	}
	return forecastToDTO("", f), nil
}

func (s *WeatherService) ByCity(ctx context.Context, city string) (*dto.WeatherResponse, error) {
	if v, ok := s.cache.Load(city); ok {
		entry := v.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.resp, nil
		}
	}

	geo, err := s.provider.Geocode(ctx, city)
	if err != nil {
		return nil, err
	}

	f, err := s.provider.Forecast(ctx, geo.Latitude, geo.Longitude)
	if err != nil {
		return nil, err
	}

	resp := forecastToDTO(geo.Name, f)
	s.cache.Store(city, cacheEntry{resp: resp, expiresAt: time.Now().Add(cacheTTL)})
	return resp, nil
}

func (s *WeatherService) ByCountry(ctx context.Context, country string) ([]*dto.WeatherResponse, error) {
	cities, ok := citiesByCountry[country]
	if !ok {
		return nil, fmt.Errorf("country not supported: %s", country)
	}

	out := make([]*dto.WeatherResponse, 0, len(cities))
	for _, city := range cities {
		w, err := s.ByCity(ctx, city)
		if err != nil {
			continue
		}
		out = append(out, w)
	}
	return out, nil
}

func (s *WeatherService) ByCountryTop(ctx context.Context, country string) ([]*dto.WeatherResponse, error) {
	all, err := s.ByCountry(ctx, country)
	if err != nil {
		return nil, err
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Temperature > all[j].Temperature })
	if len(all) > 3 {
		all = all[:3]
	}
	return all, nil
}

func forecastToDTO(city string, f *client.ForecastResponse) *dto.WeatherResponse {
	return &dto.WeatherResponse{
		City:        city,
		Latitude:    f.Latitude,
		Longitude:   f.Longitude,
		Temperature: f.Current.Temperature,
		WindSpeed:   f.Current.Windspeed,
		WeatherCode: f.Current.Weathercode,
		Time:        f.Current.Time,
	}
}
