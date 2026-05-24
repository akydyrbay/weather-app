package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"api-service/internal/model"
)

// MockUserRepository implements service.UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateWithPassword(ctx context.Context, name, email, passwordHash string) (*model.User, error) {
	args := m.Called(ctx, name, email, passwordHash)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, string, error) {
	args := m.Called(ctx, email)
	u, _ := args.Get(0).(*model.User)
	return u, args.String(1), args.Error(2)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	users, _ := args.Get(0).([]*model.User)
	return users, args.Error(1)
}

func (m *MockUserRepository) SoftDelete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCityRepository implements service.CityRepository.
type MockCityRepository struct {
	mock.Mock
}

func (m *MockCityRepository) Add(ctx context.Context, userID int, city string) (*model.City, error) {
	args := m.Called(ctx, userID, city)
	c, _ := args.Get(0).(*model.City)
	return c, args.Error(1)
}

func (m *MockCityRepository) GetByUser(ctx context.Context, userID int) ([]*model.City, error) {
	args := m.Called(ctx, userID)
	cities, _ := args.Get(0).([]*model.City)
	return cities, args.Error(1)
}

func (m *MockCityRepository) Delete(ctx context.Context, userID, cityID int) error {
	args := m.Called(ctx, userID, cityID)
	return args.Error(0)
}

// MockHistoryRepository implements service.HistoryRepository.
type MockHistoryRepository struct {
	mock.Mock
}

func (m *MockHistoryRepository) Save(ctx context.Context, userID int, city string, temperature float64, description string) error {
	args := m.Called(ctx, userID, city, temperature, description)
	return args.Error(0)
}

func (m *MockHistoryRepository) Get(ctx context.Context, userID int, f model.HistoryFilter) ([]*model.WeatherHistory, error) {
	args := m.Called(ctx, userID, f)
	h, _ := args.Get(0).([]*model.WeatherHistory)
	return h, args.Error(1)
}
