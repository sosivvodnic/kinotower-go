package user_rating_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
	user_rating_domain "github.com/sosivvodnic/kinotower-go/internal/features/user_ratings/domain"
	user_rating_service "github.com/sosivvodnic/kinotower-go/internal/features/user_ratings/service"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc user_rating_service.Service
}

func NewHandler(svc user_rating_service.Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user-id"))
	if err != nil || userID <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "User not found")
		return
	}
	var req user_rating_domain.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	resp, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		switch err {
		case user_rating_service.ErrUserNotFound:
			httpjson.WriteError(w, http.StatusNotFound, "User not found")
		case user_rating_service.ErrScoreExists:
			httpjson.WriteJSON(w, http.StatusUnauthorized, user_rating_domain.InvalidResponse{
				Status:  "invalid",
				Message: "Score exist",
			})
		case user_rating_service.ErrValidation, user_rating_service.ErrFilmNotFound:
			httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		default:
			httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, resp)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user-id"))
	if err != nil || userID <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "User not found")
		return
	}
	ratings, err := h.svc.List(r.Context(), userID)
	if err != nil {
		switch err {
		case user_rating_service.ErrUserNotFound:
			httpjson.WriteError(w, http.StatusNotFound, "User not found")
		default:
			httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, user_rating_domain.ListResponse{Ratings: ratings})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user-id"))
	if err != nil || userID <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "User not found")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "Rating not found")
		return
	}

	if err := h.svc.Delete(r.Context(), userID, id); err != nil {
		switch err {
		case user_rating_service.ErrUserNotFound:
			httpjson.WriteError(w, http.StatusNotFound, "User not found")
		case user_rating_service.ErrRatingNotFound:
			httpjson.WriteError(w, http.StatusNotFound, "Rating not found")
		default:
			httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

