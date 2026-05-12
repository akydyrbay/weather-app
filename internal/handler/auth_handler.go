package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"weather-api/internal/dto"
	"weather-api/internal/model"
)

type authSvc interface {
	Register(ctx context.Context, name, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	svc authSvc
}

func NewAuthHandler(svc authSvc) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	user, err := h.svc.Register(r.Context(), body.Name, body.Email, body.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, dto.UserFromModel(user))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	token, err := h.svc.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, dto.LoginResponse{AccessToken: token})
}
