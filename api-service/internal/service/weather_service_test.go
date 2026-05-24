package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"api-service/internal/service"
)

type mockProvider struct{ mock.Mock }

func (m *mockProvider) GetCurrentWeather(ctx context.Context, lat, lon float64) (*service.ProviderWeatherResponse, error) {
	args := m.Called(ctx, lat, lon)
	r, _ := args.Get(0).(*service.ProviderWeatherResponse)
	return r, args.Error(1)
}

func (m *mockProvider) GetCurrentCityWeather(ctx context.Context, city string) (*service.ProviderWeatherResponse, error) {
	args := m.Called(ctx, city)
	r, _ := args.Get(0).(*service.ProviderWeatherResponse)
	return r, args.Error(1)
}

func (m *mockProvider) GetCurrentCountryWeather(ctx context.Context, country string) ([]*service.ProviderWeatherResponse, error) {
	args := m.Called(ctx, country)
	r, _ := args.Get(0).([]*service.ProviderWeatherResponse)
	return r, args.Error(1)
}

func (m *mockProvider) GetCurrentCountryWeatherTop(ctx context.Context, country string) ([]*service.ProviderWeatherResponse, error) {
	args := m.Called(ctx, country)
	r, _ := args.Get(0).([]*service.ProviderWeatherResponse)
	return r, args.Error(1)
}

func TestWeatherService_GetWeather_Success(t *testing.T) {
	p := new(mockProvider)
	p.On("GetCurrentWeather", mock.Anything, 1.0, 2.0).
		Return(&service.ProviderWeatherResponse{Temperature: 20, WeatherCode: 0, Time: "t"}, nil)

	svc := service.NewWeatherService(p)
	res, err := svc.GetWeather(context.Background(), 1.0, 2.0)

	require.NoError(t, err)
	assert.Equal(t, 20.0, res.Temperature)
	assert.Equal(t, "Ясно", res.Description)
}

func TestWeatherService_GetWeather_ProviderError(t *testing.T) {
	p := new(mockProvider)
	p.On("GetCurrentWeather", mock.Anything, 1.0, 2.0).Return(nil, errors.New("boom"))

	svc := service.NewWeatherService(p)
	_, err := svc.GetWeather(context.Background(), 1.0, 2.0)

	require.Error(t, err)
}

func TestWeatherService_GetCityWeather_Success(t *testing.T) {
	p := new(mockProvider)
	p.On("GetCurrentCityWeather", mock.Anything, "Almaty").
		Return(&service.ProviderWeatherResponse{Temperature: 3, WeatherCode: 61, Latitude: 43, Longitude: 76}, nil)

	svc := service.NewWeatherService(p)
	res, err := svc.GetCityWeather(context.Background(), "Almaty")

	require.NoError(t, err)
	assert.Equal(t, "Almaty", res.City)
	assert.Equal(t, "Дождь", res.Description)
	assert.Equal(t, "тёплая одежда", res.Clothing)
}

func TestWeatherService_GetCityWeather_ProviderError(t *testing.T) {
	p := new(mockProvider)
	p.On("GetCurrentCityWeather", mock.Anything, "X").Return(nil, errors.New("not found"))

	svc := service.NewWeatherService(p)
	_, err := svc.GetCityWeather(context.Background(), "X")

	require.Error(t, err)
}

func TestWeatherService_GetCountryWeather_Success(t *testing.T) {
	p := new(mockProvider)
	p.On("GetCurrentCountryWeather", mock.Anything, "KZ").Return(
		[]*service.ProviderWeatherResponse{
			{City: "A", Temperature: 10, WeatherCode: 1},
			{City: "B", Temperature: 20, WeatherCode: 95},
		}, nil)

	svc := service.NewWeatherService(p)
	res, err := svc.GetCountryWeather(context.Background(), "KZ")

	require.NoError(t, err)
	require.Len(t, res, 2)
	assert.Equal(t, "Переменная облачность", res[0].Description)
	assert.Equal(t, "Гроза", res[1].Description)
}

func TestWeatherService_GetCountryWeather_ProviderError(t *testing.T) {
	p := new(mockProvider)
	p.On("GetCurrentCountryWeather", mock.Anything, "XX").Return(nil, errors.New("unsupported"))

	svc := service.NewWeatherService(p)
	_, err := svc.GetCountryWeather(context.Background(), "XX")

	require.Error(t, err)
}

func TestWeatherService_GetCountryWeatherTop_Success(t *testing.T) {
	p := new(mockProvider)
	p.On("GetCurrentCountryWeatherTop", mock.Anything, "KZ").Return(
		[]*service.ProviderWeatherResponse{{City: "A", Temperature: 25, WeatherCode: 45}}, nil)

	svc := service.NewWeatherService(p)
	res, err := svc.GetCountryWeatherTop(context.Background(), "KZ")

	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, "Туман", res[0].Description)
}
