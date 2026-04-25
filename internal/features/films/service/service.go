package film_service

import (
	"github.com/Otvetov/kinotower-go/internal/features/films/domain"
	film_repository "github.com/Otvetov/kinotower-go/internal/features/films/repository"
)

type FilmService interface {
	GetFilms() ([]domain.Film, error)
}

type filmService struct {
	filmRepository film_repository.FilmRepository
}

func NewFilmService(filmRepo film_repository.FilmRepository) *filmService {
	return &filmService{
		filmRepository: filmRepo,
	}
}