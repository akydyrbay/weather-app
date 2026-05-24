package model

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

type City struct {
	ID     int
	UserID int
	City   string
}

type WeatherHistory struct {
	ID          int
	UserID      int
	City        string
	Temperature float64
	Description string
	RequestedAt time.Time
}

type HistoryFilter struct {
	City   string
	Limit  int
	Offset int
}

type Weather struct {
	City        string
	Latitude    float64
	Longitude   float64
	Temperature float64
	WindSpeed   float64
	WeatherCode int
	Time        string
	Description string
	Clothing    string
}

type UserWeather struct {
	UserID  int
	Results []*Weather
}
