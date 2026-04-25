package core_router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sosivvodnic/kinotower-go/internal/core/auth"
	mw "github.com/sosivvodnic/kinotower-go/internal/core/middleware"
	auth_handler "github.com/sosivvodnic/kinotower-go/internal/features/auth/handler"
	category_handler "github.com/sosivvodnic/kinotower-go/internal/features/categories/handler"
	country_handler "github.com/sosivvodnic/kinotower-go/internal/features/countries/handler"
	film_handler "github.com/sosivvodnic/kinotower-go/internal/features/films/handler"
	gender_handler "github.com/sosivvodnic/kinotower-go/internal/features/genders/handler"
)

type Router struct {
	jwtMgr          *auth.Manager
	authHandler     auth_handler.Handler
	filmHandler     film_handler.FilmHandler
	categoryHandler category_handler.CategoryHandler
	countryHandler  country_handler.CountryHandler
	genderHandler   gender_handler.GenderHandler
}

func NewRouter(
	jwtMgr *auth.Manager,
	authHandler auth_handler.Handler,
	filmHandler film_handler.FilmHandler,
	categoryHandler category_handler.CategoryHandler,
	countryHandler country_handler.CountryHandler,
	genderHandler gender_handler.GenderHandler,
) *Router {
	return &Router{
		jwtMgr:          jwtMgr,
		authHandler:     authHandler,
		filmHandler:     filmHandler,
		categoryHandler: categoryHandler,
		countryHandler:  countryHandler,
		genderHandler:   genderHandler,
	}
}

func (r *Router) RegisterRoute() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(mw.RequestLogger) // ← наш красивый логер вместо стандартного

	router.Route("/api/v1", func(rl chi.Router) {
		rl.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})
		rl.Mount("/films", r.filmRoutes())
		rl.Mount("/categories", r.categoryRoutes())
		rl.Mount("/countries", r.countryRoutes())
		rl.Mount("/genders", r.genderRoutes())
		rl.Mount("/auth", r.authRoutes())
	})

	return router
}
