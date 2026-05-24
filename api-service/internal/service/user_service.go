package service

import (
	"context"
	"fmt"

	"api-service/internal/model"
)

type UserService struct {
	users  UserRepository
	cities CityRepository
}

func NewUserService(users UserRepository, cities CityRepository) *UserService {
	return &UserService{users: users, cities: cities}
}

func (s *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return s.users.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.users.SoftDelete(ctx, id)
}

func (s *UserService) AddCity(ctx context.Context, userID int, city string) (*model.City, error) {
	if city == "" {
		return nil, fmt.Errorf("city is required")
	}
	return s.cities.Add(ctx, userID, city)
}

func (s *UserService) GetCities(ctx context.Context, userID int) ([]*model.City, error) {
	return s.cities.GetByUser(ctx, userID)
}

func (s *UserService) DeleteCity(ctx context.Context, userID, cityID int) error {
	return s.cities.Delete(ctx, userID, cityID)
}
