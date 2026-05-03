package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"weather-api/internal/repository"
)

type authSvc interface {
	Register(ctx context.Context, name, email, password string) (*repository.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	svc authSvc
}

func NewAuthHandler(svc authSvc) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	user, err := h.svc.Register(r.Context(), body.Name, body.Email, body.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	token, err := h.svc.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"access_token": token})
}
