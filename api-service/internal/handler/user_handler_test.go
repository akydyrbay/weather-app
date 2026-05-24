package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"api-service/internal/auth"
	"api-service/internal/handler"
	"api-service/internal/middleware"
	"api-service/internal/model"
)

type mockUserSvc struct{ mock.Mock }

func (m *mockUserSvc) GetAll(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	u, _ := args.Get(0).([]*model.User)
	return u, args.Error(1)
}

func (m *mockUserSvc) GetByID(ctx context.Context, id int) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockUserSvc) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockUserSvc) AddCity(ctx context.Context, userID int, city string) (*model.City, error) {
	args := m.Called(ctx, userID, city)
	c, _ := args.Get(0).(*model.City)
	return c, args.Error(1)
}

func (m *mockUserSvc) GetCities(ctx context.Context, userID int) ([]*model.City, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]*model.City)
	return c, args.Error(1)
}

func (m *mockUserSvc) DeleteCity(ctx context.Context, userID, cityID int) error {
	return m.Called(ctx, userID, cityID).Error(0)
}

type mockUserWeatherSvc struct{ mock.Mock }

func (m *mockUserWeatherSvc) GetWeatherForUser(ctx context.Context, userID int) (*model.UserWeather, error) {
	args := m.Called(ctx, userID)
	u, _ := args.Get(0).(*model.UserWeather)
	return u, args.Error(1)
}

func (m *mockUserWeatherSvc) GetHistory(ctx context.Context, userID int, f model.HistoryFilter) (int, []*model.WeatherHistory, error) {
	args := m.Called(ctx, userID, f)
	h, _ := args.Get(1).([]*model.WeatherHistory)
	return args.Int(0), h, args.Error(2)
}

func newUserRouter(usr *mockUserSvc, uw *mockUserWeatherSvc) http.Handler {
	r := chi.NewRouter()
	h := handler.NewUserHandler(usr, uw)
	r.Get("/users/me", h.GetMe)
	return r
}

func TestUserHandler_GetMe_Success(t *testing.T) {
	usr := new(mockUserSvc)
	uw := new(mockUserWeatherSvc)
	usr.On("GetByID", mock.Anything, 42).
		Return(&model.User{ID: 42, Name: "Ali", Email: "a@b.com", Role: "user"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req = req.WithContext(middleware.WithClaims(req.Context(), &auth.Claims{UserID: 42}))
	w := httptest.NewRecorder()

	newUserRouter(usr, uw).ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.EqualValues(t, 42, body["id"])
	assert.Equal(t, "a@b.com", body["email"])
	_, hasPwd := body["password"]
	assert.False(t, hasPwd)
}

func TestUserHandler_GetMe_NotFound(t *testing.T) {
	usr := new(mockUserSvc)
	uw := new(mockUserWeatherSvc)
	usr.On("GetByID", mock.Anything, 99).Return(nil, errors.New("user not found"))

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req = req.WithContext(middleware.WithClaims(req.Context(), &auth.Claims{UserID: 99}))
	w := httptest.NewRecorder()

	newUserRouter(usr, uw).ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func decodeJSON(w *httptest.ResponseRecorder, v any) error {
	return json.NewDecoder(w.Body).Decode(v)
}
