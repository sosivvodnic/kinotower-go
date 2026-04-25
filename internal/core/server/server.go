package core_server

import (
	"net/http"

	core_database "github.com/Otvetov/kinotower-go/internal/core/database"
	core_router "github.com/Otvetov/kinotower-go/internal/core/router"
	film_handler "github.com/Otvetov/kinotower-go/internal/features/films/handler"
	film_repository "github.com/Otvetov/kinotower-go/internal/features/films/repository"
	film_service "github.com/Otvetov/kinotower-go/internal/features/films/service"
)

type Server struct {
	http.Server
}

func NewServer(db core_database.Database) *Server {
	cfg := NewConfigMust()
	filmRepository := film_repository.NewFilmRepository(db)
	filmService := film_service.NewFilmService(filmRepository)
	filmHandler := film_handler.NewFilmHandler(filmService)

	router := core_router.NewRouter(filmHandler)

	return &Server{
		Server: http.Server{
			Addr:    cfg.Addr,
			Handler: router.RegisterRoute(),
		},
	}
}
