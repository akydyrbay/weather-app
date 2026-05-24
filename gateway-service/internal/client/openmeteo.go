package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OpenMeteoClient struct {
	httpClient *http.Client
	baseURL    string
	geoURL     string
}

func NewOpenMeteoClient(httpClient *http.Client, baseURL, geoURL string) *OpenMeteoClient {
	return &OpenMeteoClient{httpClient: httpClient, baseURL: baseURL, geoURL: geoURL}
}

type ForecastResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Weathercode int     `json:"weathercode"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
}

type GeoResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c *OpenMeteoClient) Forecast(ctx context.Context, lat, lon float64) (*ForecastResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", lat))
	q.Set("longitude", fmt.Sprintf("%.4f", lon))
	q.Set("current_weather", "true")
	u.RawQuery = q.Encode()

	var resp ForecastResponse
	if err := c.getJSON(ctx, u.String(), &resp); err != nil {
		return nil, err
	}
	resp.Latitude = lat
	resp.Longitude = lon
	return &resp, nil
}

func (c *OpenMeteoClient) Geocode(ctx context.Context, name string) (*GeoResult, error) {
	u, err := url.Parse(c.geoURL)
	if err != nil {
		return nil, fmt.Errorf("parse geo url: %w", err)
	}
	q := u.Query()
	q.Set("name", name)
	q.Set("count", "1")
	q.Set("language", "en")
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	var resp struct {
		Results []GeoResult `json:"results"`
	}
	if err := c.getJSON(ctx, u.String(), &resp); err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("city not found: %s", name)
	}
	return &resp.Results[0], nil
}

func (c *OpenMeteoClient) getJSON(ctx context.Context, fullURL string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call open-meteo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("open-meteo status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
