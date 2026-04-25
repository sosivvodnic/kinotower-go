package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (r *Router) categoryRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", r.categoryHandler.GetCategories)
	return router
}

