package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"api-service/internal/dto"
	"api-service/internal/middleware"
	"api-service/internal/model"
)

type userSvc interface {
	GetAll(ctx context.Context) ([]*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	Delete(ctx context.Context, id int) error
	AddCity(ctx context.Context, userID int, city string) (*model.City, error)
	GetCities(ctx context.Context, userID int) ([]*model.City, error)
	DeleteCity(ctx context.Context, userID, cityID int) error
}

type userWeatherSvc interface {
	GetWeatherForUser(ctx context.Context, userID int) (*model.UserWeather, error)
	GetHistory(ctx context.Context, userID int, filter model.HistoryFilter) (int, []*model.WeatherHistory, error)
}

type UserHandler struct {
	users   userSvc
	weather userWeatherSvc
}

func NewUserHandler(users userSvc, weather userWeatherSvc) *UserHandler {
	return &UserHandler{users: users, weather: weather}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.GetAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, dto.UsersFromModel(users))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, dto.UserFromModel(user))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err := h.users.Delete(r.Context(), id); err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFrom(r.Context())
	user, err := h.users.GetByID(r.Context(), claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, dto.UserFromModel(user))
}

func (h *UserHandler) AddCity(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFrom(r.Context())

	var body dto.CityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	city, err := h.users.AddCity(r.Context(), claims.UserID, body.City)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, dto.CityFromModel(city))
}

func (h *UserHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFrom(r.Context())

	cities, err := h.users.GetCities(r.Context(), claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, dto.CitiesFromModel(cities))
}

func (h *UserHandler) DeleteCity(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFrom(r.Context())

	cityID, err := parseID(r, "city_id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid city id"})
		return
	}

	if err := h.users.DeleteCity(r.Context(), claims.UserID, cityID); err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: "city not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetUserWeather(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFrom(r.Context())

	result, err := h.weather.GetWeatherForUser(r.Context(), claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, dto.UserWeatherFromModel(result))
}

func (h *UserHandler) GetWeatherHistory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFrom(r.Context())

	filter := model.HistoryFilter{
		City: r.URL.Query().Get("city"),
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err := strconv.Atoi(l)
		if err != nil || limit < 0 {
			writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid limit"})
			return
		}
		filter.Limit = limit
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		offset, err := strconv.Atoi(o)
		if err != nil || offset < 0 {
			writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid offset"})
			return
		}
		filter.Offset = offset
	}

	userID, history, err := h.weather.GetHistory(r.Context(), claims.UserID, filter)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.HistoryResponse{
		UserID:  userID,
		City:    filter.City,
		History: dto.HistoryItemsFromModel(history),
	})
}

func parseID(r *http.Request, param string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, param))
}
