package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"weather-api/internal/auth"
	"weather-api/internal/repository"
)

type AuthService struct {
	repo   *repository.UserRepo
	secret string
}

func NewAuthService(repo *repository.UserRepo, secret string) *AuthService {
	return &AuthService{repo: repo, secret: secret}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*repository.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, fmt.Errorf("name, email and password are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	return s.repo.CreateWithPassword(ctx, name, email, string(hash))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, hash, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	return auth.GenerateToken(user.ID, user.Email, user.Role, s.secret)
}
