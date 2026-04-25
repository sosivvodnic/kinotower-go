package country_repository

import (
	"context"
	"fmt"

	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	"github.com/sosivvodnic/kinotower-go/internal/features/countries/domain"
)

type CountryRepository interface {
	List(ctx context.Context) ([]domain.Country, error)
}

type countryRepository struct {
	db core_database.Database
}

func NewCountryRepository(db core_database.Database) *countryRepository {
	return &countryRepository{db: db}
}

func (r *countryRepository) List(ctx context.Context) ([]domain.Country, error) {
	sqlq := `
SELECT
  co.id,
  co.name,
  COALESCE(COUNT(f.id), 0) AS film_count
FROM countries co
LEFT JOIN films f ON f.country_id = co.id AND f.deleted_at IS NULL
GROUP BY co.id, co.name
ORDER BY co.id ASC
`

	var out []domain.Country
	if err := r.db.SelectContext(ctx, &out, sqlq); err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}
	return out, nil
}
