package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (r *Router) countryRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", r.countryHandler.GetCountries)
	return router
}
