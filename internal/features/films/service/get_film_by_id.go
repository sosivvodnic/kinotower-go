package film_service

import (
	"context"

	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
)

func (s *filmService) GetFilmByID(ctx context.Context, id int) (*domain.Film, error) {
	return s.filmRepository.GetFilmByID(ctx, id)
}

func (s *filmService) GetApprovedReviewsByFilmID(ctx context.Context, filmID int) ([]domain.Review, error) {
	return s.filmRepository.GetApprovedReviewsByFilmID(ctx, filmID)
}

func (s *filmService) FilmExists(ctx context.Context, filmID int) (bool, error) {
	return s.filmRepository.FilmExists(ctx, filmID)
}

