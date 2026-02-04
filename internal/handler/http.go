package handler

import (
	"driver-service/internal/domain"
	"driver-service/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type DriverHandler struct {
	svc *service.DriverService
}

func NewDriverHandler(svc *service.DriverService) *DriverHandler {
	return &DriverHandler{svc: svc}
}

func (h *DriverHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /drivers", h.Create)
	mux.HandleFunc("GET /drivers", h.List)
	mux.HandleFunc("POST /drivers/{id}/suspend", h.Suspend)
}

func (h *DriverHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.Driver
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	driver, err := h.svc.CreateDriver(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(driver)
}

func (h *DriverHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	resp, err := h.svc.ListDrivers(r.Context(), page, limit)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *DriverHandler) Suspend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var body struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.svc.SuspendDriver(r.Context(), id, body.Reason); err != nil {
		// Differentiate between Not Found and other errors in prod
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"suspended"}`))
}
