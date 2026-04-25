package gender_repository

import (
	"context"
	"fmt"

	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	"github.com/sosivvodnic/kinotower-go/internal/features/genders/domain"
)

type GenderRepository interface {
	List(ctx context.Context) ([]domain.Gender, error)
}

type genderRepository struct {
	db core_database.Database
}

func NewGenderRepository(db core_database.Database) *genderRepository {
	return &genderRepository{db: db}
}

func (r *genderRepository) List(ctx context.Context) ([]domain.Gender, error) {
	var out []domain.Gender
	if err := r.db.SelectContext(ctx, &out, `SELECT id, name FROM genders ORDER BY id ASC`); err != nil {
		return nil, fmt.Errorf("list genders: %w", err)
	}
	return out, nil
}

