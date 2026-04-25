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

	router.Route("/{user-id}", func(ur chi.Router) {
		ur.Route("/reviews", func(rr chi.Router) {
			rr.Post("/", r.userReviewHandler.Create)
			rr.Get("/", r.userReviewHandler.List)
			rr.Delete("/{id}", r.userReviewHandler.Delete)
		})
		ur.Route("/ratings", func(rt chi.Router) {
			rt.Post("/", r.userRatingHandler.Create)
			rt.Get("/", r.userRatingHandler.List)
			rt.Delete("/{id}", r.userRatingHandler.Delete)
		})
	})
	return router
}
