package user_review_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
	user_review_domain "github.com/sosivvodnic/kinotower-go/internal/features/user_reviews/domain"
	user_review_service "github.com/sosivvodnic/kinotower-go/internal/features/user_reviews/service"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc user_review_service.Service
}

func NewHandler(svc user_review_service.Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user-id"))
	if err != nil || userID <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	var req user_review_domain.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	resp, err := h.svc.Create(r.Context(), userID, req)
	if err != nil {
		switch err {
		case user_review_service.ErrUserNotFound:
			httpjson.WriteError(w, http.StatusNotFound, "User not found")
		case user_review_service.ErrFilmNotFound:
			httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		case user_review_service.ErrValidation:
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

	revs, err := h.svc.List(r.Context(), userID)
	if err != nil {
		switch err {
		case user_review_service.ErrUserNotFound:
			httpjson.WriteError(w, http.StatusNotFound, "User not found")
		default:
			httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, user_review_domain.ListResponse{Reviews: revs})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user-id"))
	if err != nil || userID <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "User not found")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		httpjson.WriteError(w, http.StatusNotFound, "Review not found")
		return
	}

	if err := h.svc.Delete(r.Context(), userID, id); err != nil {
		switch err {
		case user_review_service.ErrUserNotFound:
			httpjson.WriteError(w, http.StatusNotFound, "User not found")
		case user_review_service.ErrReviewNotFound:
			httpjson.WriteError(w, http.StatusNotFound, "Review not found")
		default:
			httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

