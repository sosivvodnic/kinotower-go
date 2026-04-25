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
	user_rating_handler "github.com/sosivvodnic/kinotower-go/internal/features/user_ratings/handler"
	user_review_handler "github.com/sosivvodnic/kinotower-go/internal/features/user_reviews/handler"
	user_handler "github.com/sosivvodnic/kinotower-go/internal/features/users/handler"
)

type Router struct {
	jwtMgr            *auth.Manager
	authHandler       auth_handler.Handler
	filmHandler       film_handler.FilmHandler
	categoryHandler   category_handler.CategoryHandler
	countryHandler    country_handler.CountryHandler
	genderHandler     gender_handler.GenderHandler
	userHandler       user_handler.Handler
	userReviewHandler user_review_handler.Handler
	userRatingHandler user_rating_handler.Handler
}

func NewRouter(
	jwtMgr *auth.Manager,
	authHandler auth_handler.Handler,
	filmHandler film_handler.FilmHandler,
	categoryHandler category_handler.CategoryHandler,
	countryHandler country_handler.CountryHandler,
	genderHandler gender_handler.GenderHandler,
	userHandler user_handler.Handler,
	userReviewHandler user_review_handler.Handler,
	userRatingHandler user_rating_handler.Handler,
) *Router {
	return &Router{
		jwtMgr:            jwtMgr,
		authHandler:       authHandler,
		filmHandler:       filmHandler,
		categoryHandler:   categoryHandler,
		countryHandler:    countryHandler,
		genderHandler:     genderHandler,
		userHandler:       userHandler,
		userReviewHandler: userReviewHandler,
		userRatingHandler: userRatingHandler,
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

		// closed endpoints
		rl.With(mw.RequireAuth(r.jwtMgr)).Mount("/users", r.userRoutes())
	})

	return router
}
