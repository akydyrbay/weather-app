package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"weather-app/internal/model"
	"weather-app/internal/service"
	"weather-app/internal/service/mocks"
)

func TestUserWeatherService_GetHistory_Success(t *testing.T) {
	users := new(mocks.MockUserRepository)
	cities := new(mocks.MockCityRepository)
	hist := new(mocks.MockHistoryRepository)
	weather := service.NewWeatherService(new(mockProvider))

	users.On("GetByID", mock.Anything, 1).Return(&model.User{ID: 1}, nil)
	hist.On("Get", mock.Anything, 1, model.HistoryFilter{Limit: 10}).
		Return([]*model.WeatherHistory{{ID: 1, City: "Almaty"}}, nil)

	svc := service.NewUserWeatherService(users, cities, hist, weather)
	uid, items, err := svc.GetHistory(context.Background(), 1, model.HistoryFilter{Limit: 10})

	require.NoError(t, err)
	assert.Equal(t, 1, uid)
	assert.Len(t, items, 1)
}

func TestUserWeatherService_GetHistory_UserNotFound(t *testing.T) {
	users := new(mocks.MockUserRepository)
	cities := new(mocks.MockCityRepository)
	hist := new(mocks.MockHistoryRepository)
	weather := service.NewWeatherService(new(mockProvider))

	users.On("GetByID", mock.Anything, 99).Return(nil, errors.New("not found"))

	svc := service.NewUserWeatherService(users, cities, hist, weather)
	_, _, err := svc.GetHistory(context.Background(), 99, model.HistoryFilter{})

	require.Error(t, err)
	hist.AssertNotCalled(t, "Get")
}

func TestUserWeatherService_GetHistory_RepoError(t *testing.T) {
	users := new(mocks.MockUserRepository)
	cities := new(mocks.MockCityRepository)
	hist := new(mocks.MockHistoryRepository)
	weather := service.NewWeatherService(new(mockProvider))

	users.On("GetByID", mock.Anything, 1).Return(&model.User{ID: 1}, nil)
	hist.On("Get", mock.Anything, 1, mock.Anything).Return(nil, errors.New("db down"))

	svc := service.NewUserWeatherService(users, cities, hist, weather)
	_, _, err := svc.GetHistory(context.Background(), 1, model.HistoryFilter{})

	require.Error(t, err)
}

func TestUserWeatherService_GetWeatherForUser_NoCities(t *testing.T) {
	users := new(mocks.MockUserRepository)
	cities := new(mocks.MockCityRepository)
	hist := new(mocks.MockHistoryRepository)
	weather := service.NewWeatherService(new(mockProvider))

	users.On("GetByID", mock.Anything, 1).Return(&model.User{ID: 1}, nil)
	cities.On("GetByUser", mock.Anything, 1).Return([]*model.City{}, nil)

	svc := service.NewUserWeatherService(users, cities, hist, weather)
	res, err := svc.GetWeatherForUser(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 1, res.UserID)
	assert.Empty(t, res.Results)
}

func TestUserWeatherService_GetWeatherForUser_UserNotFound(t *testing.T) {
	users := new(mocks.MockUserRepository)
	cities := new(mocks.MockCityRepository)
	hist := new(mocks.MockHistoryRepository)
	weather := service.NewWeatherService(new(mockProvider))

	users.On("GetByID", mock.Anything, 42).Return(nil, errors.New("not found"))

	svc := service.NewUserWeatherService(users, cities, hist, weather)
	_, err := svc.GetWeatherForUser(context.Background(), 42)

	require.Error(t, err)
}
