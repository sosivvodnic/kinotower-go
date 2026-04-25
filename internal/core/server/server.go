package core_server

import (
	"net/http"
	"time"

	"github.com/sosivvodnic/kinotower-go/internal/core/auth"
	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	core_router "github.com/sosivvodnic/kinotower-go/internal/core/router"
	auth_handler "github.com/sosivvodnic/kinotower-go/internal/features/auth/handler"
	auth_repository "github.com/sosivvodnic/kinotower-go/internal/features/auth/repository"
	auth_service "github.com/sosivvodnic/kinotower-go/internal/features/auth/service"
	category_handler "github.com/sosivvodnic/kinotower-go/internal/features/categories/handler"
	category_repository "github.com/sosivvodnic/kinotower-go/internal/features/categories/repository"
	category_service "github.com/sosivvodnic/kinotower-go/internal/features/categories/service"
	country_handler "github.com/sosivvodnic/kinotower-go/internal/features/countries/handler"
	country_repository "github.com/sosivvodnic/kinotower-go/internal/features/countries/repository"
	country_service "github.com/sosivvodnic/kinotower-go/internal/features/countries/service"
	film_handler "github.com/sosivvodnic/kinotower-go/internal/features/films/handler"
	film_repository "github.com/sosivvodnic/kinotower-go/internal/features/films/repository"
	film_service "github.com/sosivvodnic/kinotower-go/internal/features/films/service"
	gender_handler "github.com/sosivvodnic/kinotower-go/internal/features/genders/handler"
	gender_repository "github.com/sosivvodnic/kinotower-go/internal/features/genders/repository"
	gender_service "github.com/sosivvodnic/kinotower-go/internal/features/genders/service"
)

type Server struct {
	http.Server
}

func NewServer(db core_database.Database) *Server {
	cfg := NewConfigMust()
	authCfg := NewAuthConfigMust()
	jwtMgr := auth.NewManager(authCfg.JWTSecret, 24*time.Hour)

	filmRepository := film_repository.NewFilmRepository(db)
	filmService := film_service.NewFilmService(filmRepository)
	filmHandler := film_handler.NewFilmHandler(filmService)

	authRepo := auth_repository.NewRepository(db)
	authSvc := auth_service.NewService(authRepo, jwtMgr)
	authHandler := auth_handler.NewHandler(authSvc, jwtMgr)

	categoryRepository := category_repository.NewCategoryRepository(db)
	categoryService := category_service.NewCategoryService(categoryRepository)
	categoryHandler := category_handler.NewCategoryHandler(categoryService)

	countryRepository := country_repository.NewCountryRepository(db)
	countryService := country_service.NewCountryService(countryRepository)
	countryHandler := country_handler.NewCountryHandler(countryService)

	genderRepository := gender_repository.NewGenderRepository(db)
	genderService := gender_service.NewGenderService(genderRepository)
	genderHandler := gender_handler.NewGenderHandler(genderService)

	router := core_router.NewRouter(jwtMgr, authHandler, filmHandler, categoryHandler, countryHandler, genderHandler)

	return &Server{
		Server: http.Server{
			Addr:    cfg.Addr,
			Handler: router.RegisterRoute(),
		},
	}
}
