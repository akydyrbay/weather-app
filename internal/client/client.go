package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"weather-api/internal/service"
)

const cacheTTL = 5 * time.Minute

type cacheEntry struct {
	resp      *service.ProviderWeatherResponse
	expiresAt time.Time
}

type WeatherClient struct {
	httpClient *http.Client
	baseURL    string
	geoURL     string
	cache      sync.Map
}

func NewWeatherClient(httpClient *http.Client) *WeatherClient {
	return &WeatherClient{
		httpClient: httpClient,
		baseURL:    "https://api.open-meteo.com/v1/forecast",
		geoURL:     "https://geocoding-api.open-meteo.com/v1/search",
	}
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

type openMeteoResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Weathercode int     `json:"weathercode"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
}

func (c *WeatherClient) GetCurrentWeather(ctx context.Context, lat, lon float64) (*service.ProviderWeatherResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", lat))
	q.Set("longitude", fmt.Sprintf("%.4f", lon))
	q.Set("current_weather", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call external api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external api returned status: %d", resp.StatusCode)
	}

	var result openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode external api response: %w", err)
	}

	return &service.ProviderWeatherResponse{
		Latitude:    lat,
		Longitude:   lon,
		Temperature: result.CurrentWeather.Temperature,
		WindSpeed:   result.CurrentWeather.Windspeed,
		WeatherCode: result.CurrentWeather.Weathercode,
		Time:        result.CurrentWeather.Time,
	}, nil
}

func (c *WeatherClient) GetCurrentCityWeather(ctx context.Context, city string) (*service.ProviderWeatherResponse, error) {
	if v, ok := c.cache.Load(city); ok {
		entry := v.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.resp, nil
		}
	}

	u, err := url.Parse(c.geoURL)
	if err != nil {
		return nil, fmt.Errorf("parse geo url: %w", err)
	}

	q := u.Query()
	q.Set("name", city)
	q.Set("count", "1")
	q.Set("language", "en")
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create geo request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call geo api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geo api returned status: %d", resp.StatusCode)
	}

	var geoResult struct {
		Results []struct {
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&geoResult); err != nil {
		return nil, fmt.Errorf("decode geo response: %w", err)
	}

	if len(geoResult.Results) == 0 {
		return nil, fmt.Errorf("city not found")
	}

	lat := geoResult.Results[0].Latitude
	lon := geoResult.Results[0].Longitude
	cityName := geoResult.Results[0].Name

	weatherResp, err := c.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return nil, err
	}

	weatherResp.City = cityName
	weatherResp.Latitude = lat
	weatherResp.Longitude = lon

	c.cache.Store(city, cacheEntry{resp: weatherResp, expiresAt: time.Now().Add(cacheTTL)})

	return weatherResp, nil
}

func (c *WeatherClient) GetCurrentCountryWeather(ctx context.Context, country string) ([]*service.ProviderWeatherResponse, error) {
	cities, ok := citiesByCountry[country]
	if !ok {
		return nil, fmt.Errorf("country not supported: %s", country)
	}

	results := make([]*service.ProviderWeatherResponse, 0, len(cities))
	for _, city := range cities {
		weather, err := c.GetCurrentCityWeather(ctx, city)
		if err != nil {
			continue
		}
		results = append(results, weather)
	}

	return results, nil
}

func (c *WeatherClient) GetCurrentCountryWeatherTop(ctx context.Context, country string) ([]*service.ProviderWeatherResponse, error) {
	all, err := c.GetCurrentCountryWeather(ctx, country)
	if err != nil {
		return nil, err
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Temperature > all[j].Temperature
	})

	if len(all) > 3 {
		all = all[:3]
	}

	return all, nil
}
