package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"gateway-service/internal/dto"
)

type weatherSvc interface {
	ByCoords(ctx context.Context, lat, lon float64) (*dto.WeatherResponse, error)
	ByCity(ctx context.Context, city string) (*dto.WeatherResponse, error)
	ByCountry(ctx context.Context, country string) ([]*dto.WeatherResponse, error)
	ByCountryTop(ctx context.Context, country string) ([]*dto.WeatherResponse, error)
}

type WeatherHandler struct {
	svc weatherSvc
}

func NewWeatherHandler(svc weatherSvc) *WeatherHandler {
	return &WeatherHandler{svc: svc}
}

func (h *WeatherHandler) ByCoords(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid lat"})
		return
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid lon"})
		return
	}

	resp, err := h.svc.ByCoords(r.Context(), lat, lon)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *WeatherHandler) ByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "city is required"})
		return
	}
	resp, err := h.svc.ByCity(r.Context(), city)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *WeatherHandler) ByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	resp, err := h.svc.ByCountry(r.Context(), country)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *WeatherHandler) ByCountryTop(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	resp, err := h.svc.ByCountryTop(r.Context(), country)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
