package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"weather-api/internal/repository"
	"weather-api/internal/service"
)

type UserHandler struct {
	users   *service.UserService
	weather *service.UserWeatherService
	cities  *repository.CityRepo
}

func NewUserHandler(
	users *service.UserService,
	weather *service.UserWeatherService,
	cities *repository.CityRepo,
) *UserHandler {
	return &UserHandler{users: users, weather: weather, cities: cities}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	user, err := h.users.Create(r.Context(), body.Name, body.Email)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.GetAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}

	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}

	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	user, err := h.users.Update(r.Context(), id, body.Name, body.Email)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.users.Delete(r.Context(), id); err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) AddCity(w http.ResponseWriter, r *http.Request) {
	userID, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}

	var body struct {
		City string `json:"city"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.City == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "city is required"})
		return
	}

	city, err := h.cities.Add(r.Context(), userID, body.City)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, city)
}
func (h *UserHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	userID, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}

	cities, err := h.cities.GetByUser(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, cities)
}

func (h *UserHandler) DeleteCity(w http.ResponseWriter, r *http.Request) {
	userID, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	cityID, err := parseID(r, "city_id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid city id"})
		return
	}

	if err := h.cities.Delete(r.Context(), userID, cityID); err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "city not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetUserWeather(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.weather.GetWeatherForUser(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *UserHandler) GetWeatherHistory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}

	filter := repository.HistoryFilter{
		City: r.URL.Query().Get("city"),
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err := strconv.Atoi(l)
		if err != nil || limit < 0 {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid limit"})
			return
		}
		filter.Limit = limit
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		offset, err := strconv.Atoi(o)
		if err != nil || offset < 0 {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid offset"})
			return
		}
		filter.Offset = offset
	}

	result, err := h.weather.GetHistory(r.Context(), id, filter)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func parseID(r *http.Request, param string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, param))
}
