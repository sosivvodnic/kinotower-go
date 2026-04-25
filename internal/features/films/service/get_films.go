package film_service

import "github.com/Otvetov/kinotower-go/internal/features/films/domain"

func (s *filmService) GetFilms() ([]domain.Film, error) {
	return []domain.Film{}, nil
}
