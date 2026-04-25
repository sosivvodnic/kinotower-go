package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (r *Router) userRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/{id}", r.userHandler.GetUser)
	router.Put("/", r.userHandler.UpdateMe)
	router.Delete("/", r.userHandler.DeleteMe)
	return router
}

