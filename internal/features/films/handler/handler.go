package film_handler

import (
	"net/http"

	film_service "github.com/sosivvodnic/kinotower-go/internal/features/films/service"
)

type FilmHandler interface {
	GetFilms(w http.ResponseWriter, r *http.Request)
}

type filmHandler struct {
	filmService film_service.FilmService
}

func NewFilmHandler(filmService film_service.FilmService) *filmHandler {
	return &filmHandler{
		filmService: filmService,
	}
}