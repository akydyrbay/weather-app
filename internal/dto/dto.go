package dto

import (
	"time"

	"weather-api/internal/model"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func UserFromModel(u *model.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

func UsersFromModel(users []*model.User) []*UserResponse {
	out := make([]*UserResponse, len(users))
	for i, u := range users {
		out[i] = UserFromModel(u)
	}
	return out
}

type CityRequest struct {
	City string `json:"city"`
}

type CityResponse struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	City   string `json:"city"`
}

func CityFromModel(c *model.City) *CityResponse {
	if c == nil {
		return nil
	}
	return &CityResponse{ID: c.ID, UserID: c.UserID, City: c.City}
}

func CitiesFromModel(cs []*model.City) []*CityResponse {
	out := make([]*CityResponse, len(cs))
	for i, c := range cs {
		out[i] = CityFromModel(c)
	}
	return out
}

type HistoryItemResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	Description string    `json:"description"`
	RequestedAt time.Time `json:"requested_at"`
}

type HistoryResponse struct {
	UserID  int                    `json:"user_id"`
	City    string                 `json:"city,omitempty"`
	History []*HistoryItemResponse `json:"history"`
}

func HistoryItemFromModel(h *model.WeatherHistory) *HistoryItemResponse {
	if h == nil {
		return nil
	}
	return &HistoryItemResponse{
		ID:          h.ID,
		UserID:      h.UserID,
		City:        h.City,
		Temperature: h.Temperature,
		Description: h.Description,
		RequestedAt: h.RequestedAt,
	}
}

func HistoryItemsFromModel(items []*model.WeatherHistory) []*HistoryItemResponse {
	out := make([]*HistoryItemResponse, len(items))
	for i, h := range items {
		out[i] = HistoryItemFromModel(h)
	}
	return out
}

type WeatherResponse struct {
	City        string  `json:"city,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
	Description string  `json:"description"`
	Clothing    string  `json:"clothing,omitempty"`
}

func WeatherFromModel(w *model.Weather) *WeatherResponse {
	if w == nil {
		return nil
	}
	return &WeatherResponse{
		City:        w.City,
		Latitude:    w.Latitude,
		Longitude:   w.Longitude,
		Temperature: w.Temperature,
		WindSpeed:   w.WindSpeed,
		WeatherCode: w.WeatherCode,
		Time:        w.Time,
		Description: w.Description,
		Clothing:    w.Clothing,
	}
}

func WeathersFromModel(ws []*model.Weather) []*WeatherResponse {
	out := make([]*WeatherResponse, len(ws))
	for i, w := range ws {
		out[i] = WeatherFromModel(w)
	}
	return out
}

type UserWeatherResponse struct {
	UserID  int                `json:"user_id"`
	Results []*WeatherResponse `json:"results"`
}

func UserWeatherFromModel(uw *model.UserWeather) *UserWeatherResponse {
	if uw == nil {
		return nil
	}
	return &UserWeatherResponse{
		UserID:  uw.UserID,
		Results: WeathersFromModel(uw.Results),
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
