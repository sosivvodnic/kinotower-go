package country_handler

import (
	"net/http"

	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
	"github.com/sosivvodnic/kinotower-go/internal/features/countries/domain"
	country_service "github.com/sosivvodnic/kinotower-go/internal/features/countries/service"
)

type CountryHandler interface {
	GetCountries(w http.ResponseWriter, r *http.Request)
}

type countryHandler struct {
	svc country_service.CountryService
}

func NewCountryHandler(svc country_service.CountryService) *countryHandler {
	return &countryHandler{svc: svc}
}

func (h *countryHandler) GetCountries(w http.ResponseWriter, r *http.Request) {
	countries, err := h.svc.List(r.Context())
	if err != nil {
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, domain.ListResponse{Countries: countries})
}
