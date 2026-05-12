package service

import (
	"context"

	"weather-api/internal/model"
)

type UserRepository interface {
	CreateWithPassword(ctx context.Context, name, email, passwordHash string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, string, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	SoftDelete(ctx context.Context, id int) error
}

type CityRepository interface {
	Add(ctx context.Context, userID int, city string) (*model.City, error)
	GetByUser(ctx context.Context, userID int) ([]*model.City, error)
	Delete(ctx context.Context, userID, cityID int) error
}

type HistoryRepository interface {
	Save(ctx context.Context, userID int, city string, temperature float64, description string) error
	Get(ctx context.Context, userID int, f model.HistoryFilter) ([]*model.WeatherHistory, error)
}
