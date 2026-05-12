package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"weather-api/internal/dto"
	"weather-api/internal/model"
)

type weatherSvc interface {
	GetWeather(ctx context.Context, lat, lon float64) (*model.Weather, error)
	GetCityWeather(ctx context.Context, city string) (*model.Weather, error)
	GetCountryWeather(ctx context.Context, country string) ([]*model.Weather, error)
	GetCountryWeatherTop(ctx context.Context, country string) ([]*model.Weather, error)
}

type WeatherHandler struct {
	service weatherSvc
}

func NewWeatherHandler(service weatherSvc) *WeatherHandler {
	return &WeatherHandler{service: service}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "query params lat and lon are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid lat"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid lon"})
		return
	}

	result, err := h.service.GetWeather(r.Context(), lat, lon)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.WeatherFromModel(result))
}

func (h *WeatherHandler) GetCityWeather(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "city parameter is required"})
		return
	}

	result, err := h.service.GetCityWeather(r.Context(), city)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.WeatherFromModel(result))
}

func (h *WeatherHandler) GetCountryWeather(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	if country == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "country parameter is required"})
		return
	}

	results, err := h.service.GetCountryWeather(r.Context(), country)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.WeathersFromModel(results))
}

func (h *WeatherHandler) GetCountryWeatherTop(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	if country == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "country parameter is required"})
		return
	}

	results, err := h.service.GetCountryWeatherTop(r.Context(), country)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dto.WeathersFromModel(results))
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		http.Error(w, `{"error":"failed to encode json"}`, http.StatusInternalServerError)
	}
}
