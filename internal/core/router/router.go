package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) RegisterRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/api/v1", func (rl chi.Router)  {
		rl.Get("/", func (w http.ResponseWriter, r *http.Request)  {
				w.Write([]byte("ping"))
		})

		rl.Mount("/films", r.filmRoutes())
		rl.Mount("/genders", r.genderRoutes())
	})

	return router
}
