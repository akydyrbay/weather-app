package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"api-service/internal/model"
	"api-service/internal/service"
	"api-service/internal/service/mocks"
)

func TestAuthService_Register_Success(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	expected := &model.User{ID: 1, Name: "Ali", Email: "ali@x.com", Role: "user"}

	repo.On("CreateWithPassword",
		mock.Anything, "Ali", "ali@x.com", mock.AnythingOfType("string"),
	).Return(expected, nil)

	svc := service.NewAuthService(repo, "secret")
	user, err := svc.Register(context.Background(), "Ali", "ali@x.com", "password")

	require.NoError(t, err)
	assert.Equal(t, expected, user)
	repo.AssertExpectations(t)
}

func TestAuthService_Register_ValidationErrors(t *testing.T) {
	cases := []struct {
		name, n, e, p string
	}{
		{"empty name", "", "x@x.com", "pass"},
		{"empty email", "Ali", "", "pass"},
		{"empty password", "Ali", "x@x.com", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := new(mocks.MockUserRepository)
			svc := service.NewAuthService(repo, "secret")

			_, err := svc.Register(context.Background(), c.n, c.e, c.p)

			require.Error(t, err)
			repo.AssertNotCalled(t, "CreateWithPassword")
		})
	}
}

func TestAuthService_Register_RepoError(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	repo.On("CreateWithPassword",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(nil, errors.New("db down"))

	svc := service.NewAuthService(repo, "secret")
	_, err := svc.Register(context.Background(), "n", "e", "p")

	require.Error(t, err)
	repo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo := new(mocks.MockUserRepository)
	repo.On("GetByEmail", mock.Anything, "x@x.com").
		Return(&model.User{ID: 1, Email: "x@x.com", Role: "user"}, string(hash), nil)

	svc := service.NewAuthService(repo, "secret")
	token, err := svc.Login(context.Background(), "x@x.com", "pass")

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	repo.On("GetByEmail", mock.Anything, "missing@x.com").
		Return(nil, "", errors.New("user not found"))

	svc := service.NewAuthService(repo, "secret")
	_, err := svc.Login(context.Background(), "missing@x.com", "pass")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo := new(mocks.MockUserRepository)
	repo.On("GetByEmail", mock.Anything, "x@x.com").
		Return(&model.User{ID: 1, Email: "x@x.com", Role: "user"}, string(hash), nil)

	svc := service.NewAuthService(repo, "secret")
	_, err := svc.Login(context.Background(), "x@x.com", "wrong")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}
