package film_repository

import (
	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
)

type FilmRepository interface {
	GetFilms() ([]domain.Film, error)
}

type filmRepository struct {
	db core_database.Database
}

func NewFilmRepository(db core_database.Database) *filmRepository {
	return &filmRepository{db: db}
}
