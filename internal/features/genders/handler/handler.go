package gender_handler

import (
	"net/http"

	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
	"github.com/sosivvodnic/kinotower-go/internal/features/genders/domain"
	gender_service "github.com/sosivvodnic/kinotower-go/internal/features/genders/service"
)

type GenderHandler interface {
	GetGenders(w http.ResponseWriter, r *http.Request)
}

type genderHandler struct {
	svc gender_service.GenderService
}

func NewGenderHandler(svc gender_service.GenderService) *genderHandler {
	return &genderHandler{svc: svc}
}

func (h *genderHandler) GetGenders(w http.ResponseWriter, r *http.Request) {
	genders, err := h.svc.List(r.Context())
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, domain.ListResponse{Genders: genders})
}
