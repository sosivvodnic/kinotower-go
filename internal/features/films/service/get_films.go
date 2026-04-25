package film_service

import (
	"context"

	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
)

func (s *filmService) GetFilms(ctx context.Context, f domain.FilmFilter) ([]domain.Film, int, error) {
	return s.filmRepository.GetFilms(ctx, f)
}
