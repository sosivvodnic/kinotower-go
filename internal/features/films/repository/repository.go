package film_repository

import (
	"context"

	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
)

type FilmRepository interface {
	GetFilms(ctx context.Context, f domain.FilmFilter) ([]domain.Film, int, error)
	GetFilmByID(ctx context.Context, id int) (*domain.Film, error)
	GetApprovedReviewsByFilmID(ctx context.Context, filmID int) ([]domain.Review, error)
	FilmExists(ctx context.Context, filmID int) (bool, error)
}

type filmRepository struct {
	db core_database.Database
}

func NewFilmRepository(db core_database.Database) *filmRepository {
	return &filmRepository{db: db}
}
