package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"api-service/internal/service"
)

type GatewayClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewGatewayClient(httpClient *http.Client, baseURL string) *GatewayClient {
	return &GatewayClient{httpClient: httpClient, baseURL: baseURL}
}

type weatherDTO struct {
	City        string  `json:"city,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
}

func (c *GatewayClient) GetCurrentWeather(ctx context.Context, lat, lon float64) (*service.ProviderWeatherResponse, error) {
	q := url.Values{}
	q.Set("lat", fmt.Sprintf("%.4f", lat))
	q.Set("lon", fmt.Sprintf("%.4f", lon))

	var w weatherDTO
	if err := c.getJSON(ctx, "/weather?"+q.Encode(), &w); err != nil {
		return nil, err
	}
	return toProvider(w), nil
}

func (c *GatewayClient) GetCurrentCityWeather(ctx context.Context, city string) (*service.ProviderWeatherResponse, error) {
	var w weatherDTO
	if err := c.getJSON(ctx, "/weather/city/"+url.PathEscape(city), &w); err != nil {
		return nil, err
	}
	return toProvider(w), nil
}

func (c *GatewayClient) GetCurrentCountryWeather(ctx context.Context, country string) ([]*service.ProviderWeatherResponse, error) {
	var ws []weatherDTO
	if err := c.getJSON(ctx, "/weather/country/"+url.PathEscape(country), &ws); err != nil {
		return nil, err
	}
	return toProviderList(ws), nil
}

func (c *GatewayClient) GetCurrentCountryWeatherTop(ctx context.Context, country string) ([]*service.ProviderWeatherResponse, error) {
	var ws []weatherDTO
	if err := c.getJSON(ctx, "/weather/country/"+url.PathEscape(country)+"/top", &ws); err != nil {
		return nil, err
	}
	return toProviderList(ws), nil
}

func (c *GatewayClient) getJSON(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("create gateway request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gateway returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode gateway response: %w", err)
	}
	return nil
}

func toProvider(w weatherDTO) *service.ProviderWeatherResponse {
	return &service.ProviderWeatherResponse{
		City:        w.City,
		Latitude:    w.Latitude,
		Longitude:   w.Longitude,
		Temperature: w.Temperature,
		WindSpeed:   w.WindSpeed,
		WeatherCode: w.WeatherCode,
		Time:        w.Time,
	}
}

func toProviderList(ws []weatherDTO) []*service.ProviderWeatherResponse {
	out := make([]*service.ProviderWeatherResponse, len(ws))
	for i, w := range ws {
		out[i] = toProvider(w)
	}
	return out
}
