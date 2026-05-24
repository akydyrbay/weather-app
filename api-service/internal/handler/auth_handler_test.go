package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"api-service/internal/dto"
	"api-service/internal/handler"
	"api-service/internal/model"
)

type mockAuthSvc struct{ mock.Mock }

func (m *mockAuthSvc) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	args := m.Called(ctx, name, email, password)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAuthSvc) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func newAuthRouter(svc *mockAuthSvc) http.Handler {
	r := chi.NewRouter()
	h := handler.NewAuthHandler(svc)
	r.Post("/auth/register", h.Register)
	return r
}

func TestAuthHandler_Register_Success(t *testing.T) {
	svc := new(mockAuthSvc)
	svc.On("Register", mock.Anything, "Ali", "a@b.com", "secret").
		Return(&model.User{ID: 1, Name: "Ali", Email: "a@b.com", Role: "user"}, nil)

	body := strings.NewReader(`{"name":"Ali","email":"a@b.com","password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	newAuthRouter(svc).ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.NotContains(t, w.Body.String(), "password")
	assert.Contains(t, w.Body.String(), `"email": "a@b.com"`)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	svc := new(mockAuthSvc)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{not json`))
	w := httptest.NewRecorder()

	newAuthRouter(svc).ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var resp dto.ErrorResponse
	require.NoError(t, decodeJSON(w, &resp))
	assert.NotEmpty(t, resp.Error)
	svc.AssertNotCalled(t, "Register")
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	svc := new(mockAuthSvc)
	svc.On("Register", mock.Anything, "", "a@b.com", "secret").
		Return(nil, errors.New("name, email and password are required"))

	body := strings.NewReader(`{"name":"","email":"a@b.com","password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", body)
	w := httptest.NewRecorder()

	newAuthRouter(svc).ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
