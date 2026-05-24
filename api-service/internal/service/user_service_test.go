package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"api-service/internal/model"
	"api-service/internal/service"
	"api-service/internal/service/mocks"
)

func TestUserService_AddCity_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	expected := &model.City{ID: 1, UserID: 7, City: "Almaty"}

	cityRepo.On("Add", mock.Anything, 7, "Almaty").Return(expected, nil)

	svc := service.NewUserService(userRepo, cityRepo)
	got, err := svc.AddCity(context.Background(), 7, "Almaty")

	require.NoError(t, err)
	assert.Equal(t, expected, got)
	cityRepo.AssertExpectations(t)
}

func TestUserService_AddCity_EmptyCity(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)

	svc := service.NewUserService(userRepo, cityRepo)
	_, err := svc.AddCity(context.Background(), 7, "")

	require.Error(t, err)
	cityRepo.AssertNotCalled(t, "Add")
}

func TestUserService_AddCity_RepoError(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	cityRepo.On("Add", mock.Anything, 7, "Almaty").Return(nil, errors.New("db error"))

	svc := service.NewUserService(userRepo, cityRepo)
	_, err := svc.AddCity(context.Background(), 7, "Almaty")

	require.Error(t, err)
}

func TestUserService_GetCities_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	expected := []*model.City{{ID: 1, UserID: 7, City: "Almaty"}, {ID: 2, UserID: 7, City: "Astana"}}
	cityRepo.On("GetByUser", mock.Anything, 7).Return(expected, nil)

	svc := service.NewUserService(userRepo, cityRepo)
	got, err := svc.GetCities(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestUserService_DeleteCity_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	cityRepo.On("Delete", mock.Anything, 7, 1).Return(nil)

	svc := service.NewUserService(userRepo, cityRepo)
	err := svc.DeleteCity(context.Background(), 7, 1)

	require.NoError(t, err)
}

func TestUserService_DeleteCity_NotFound(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	cityRepo.On("Delete", mock.Anything, 7, 99).Return(errors.New("city not found"))

	svc := service.NewUserService(userRepo, cityRepo)
	err := svc.DeleteCity(context.Background(), 7, 99)

	require.Error(t, err)
}

func TestUserService_GetByID_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	expected := &model.User{ID: 1, Name: "Ali"}
	userRepo.On("GetByID", mock.Anything, 1).Return(expected, nil)

	svc := service.NewUserService(userRepo, cityRepo)
	got, err := svc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	userRepo.On("GetByID", mock.Anything, 999).Return(nil, errors.New("user not found"))

	svc := service.NewUserService(userRepo, cityRepo)
	_, err := svc.GetByID(context.Background(), 999)

	require.Error(t, err)
}

func TestUserService_GetAll_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	expected := []*model.User{{ID: 1}, {ID: 2}}
	userRepo.On("GetAll", mock.Anything).Return(expected, nil)

	svc := service.NewUserService(userRepo, cityRepo)
	got, err := svc.GetAll(context.Background())

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestUserService_Delete_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	cityRepo := new(mocks.MockCityRepository)
	userRepo.On("SoftDelete", mock.Anything, 5).Return(nil)

	svc := service.NewUserService(userRepo, cityRepo)
	err := svc.Delete(context.Background(), 5)

	require.NoError(t, err)
}
