package service

import (
	"context"
	"fmt"

	"weather-api/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, name, email string) (*repository.User, error) {
	if name == "" || email == "" {
		return nil, fmt.Errorf("name and email are required")
	}
	return s.repo.Create(ctx, name, email)
}

func (s *UserService) GetAll(ctx context.Context) ([]*repository.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id int) (*repository.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id int, name, email string) (*repository.User, error) {
	if name == "" || email == "" {
		return nil, fmt.Errorf("name and email are required")
	}
	return s.repo.Update(ctx, id, name, email)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
