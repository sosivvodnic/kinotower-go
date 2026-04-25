package film_service

import (
	"context"

	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
	film_repository "github.com/sosivvodnic/kinotower-go/internal/features/films/repository"
)

type FilmService interface {
	GetFilms(ctx context.Context, f domain.FilmFilter) ([]domain.Film, int, error)
	GetFilmByID(ctx context.Context, id int) (*domain.Film, error)
	GetApprovedReviewsByFilmID(ctx context.Context, filmID int) ([]domain.Review, error)
	FilmExists(ctx context.Context, filmID int) (bool, error)
}

type filmService struct {
	filmRepository film_repository.FilmRepository
}

func NewFilmService(filmRepo film_repository.FilmRepository) *filmService {
	return &filmService{
		filmRepository: filmRepo,
	}
}
