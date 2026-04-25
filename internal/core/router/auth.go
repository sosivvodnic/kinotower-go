package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	mw "github.com/sosivvodnic/kinotower-go/internal/core/middleware"
)

func (r *Router) authRoutes() http.Handler {
	router := chi.NewRouter()
	router.Post("/signup", r.authHandler.SignUp)
	router.Post("/signin", r.authHandler.SignIn)
	router.With(mw.RequireAuth(r.jwtMgr)).Post("/signout", r.authHandler.SignOut)
	return router
}
