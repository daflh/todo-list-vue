package todo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreateRequest struct {
	Task     string `json:"task"`
	Priority string `json:"priority"`
}

type UpdateRequest struct {
	Task     *string `json:"task"`
	Priority *string `json:"priority"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type DeleteRequest struct {
	IDs []int64 `json:"ids"`
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	todos, err := h.svc.list(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list todos: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, todos)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	if req.Task == "" {
		http.Error(w, "Task is required", http.StatusBadRequest)
		return
	}

	todo, err := h.svc.create(r.Context(), req.Task, req.Priority)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create todo: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, todo)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	todo, err := h.svc.update(r.Context(), id, req.Task, req.Priority)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Todo not found for update", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update todo: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	todo, err := h.svc.updateStatus(r.Context(), id, req.Status)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Todo not found for update", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update todo: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "No IDs provided for deletion", http.StatusBadRequest)
		return
	}

	deleted, err := h.svc.delete(r.Context(), req.IDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete todos: %v", err), http.StatusInternalServerError)
		return
	}
	if deleted == 0 {
		http.Error(w, "Todo not found for deletion", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
