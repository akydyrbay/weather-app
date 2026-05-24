package dto

type WeatherResponse struct {
	City        string  `json:"city,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
