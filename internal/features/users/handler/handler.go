package user_handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
	mw "github.com/sosivvodnic/kinotower-go/internal/core/middleware"
	"github.com/sosivvodnic/kinotower-go/internal/features/users/domain"
	user_service "github.com/sosivvodnic/kinotower-go/internal/features/users/service"
)

type Handler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateMe(w http.ResponseWriter, r *http.Request)
	DeleteMe(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc user_service.Service
}

func NewHandler(svc user_service.Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	u, err := h.svc.GetUserByID(r.Context(), id)
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if u == nil {
		httpjson.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, u)
}

func (h *handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := mw.UserIDFromContext(r.Context())
	if !ok || userID <= 0 {
		httpjson.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}
	if _, err := time.Parse("2006-01-02", req.Birthday); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.svc.UpdateUser(r.Context(), userID, req); err != nil {
		if err == user_service.ErrValidation {
			httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
			return
		}
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]any{"status": "success"})
}

func (h *handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := mw.UserIDFromContext(r.Context())
	if !ok || userID <= 0 {
		httpjson.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.svc.DeleteUser(r.Context(), userID); err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// 204 without JSON is allowed by spec
	w.WriteHeader(http.StatusNoContent)
}

