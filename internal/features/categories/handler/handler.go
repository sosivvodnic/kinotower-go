package category_handler

import (
	"net/http"

	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
	"github.com/sosivvodnic/kinotower-go/internal/features/categories/domain"
	category_service "github.com/sosivvodnic/kinotower-go/internal/features/categories/service"
)

type CategoryHandler interface {
	GetCategories(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	svc category_service.CategoryService
}

func NewCategoryHandler(svc category_service.CategoryService) *categoryHandler {
	return &categoryHandler{svc: svc}
}

func (h *categoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.List(r.Context())
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, domain.ListResponse{Categories: cats})
}

