package category_repository

import (
	"context"
	"fmt"

	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	"github.com/sosivvodnic/kinotower-go/internal/features/categories/domain"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]domain.Category, error)
}

type categoryRepository struct {
	db core_database.Database
}

func NewCategoryRepository(db core_database.Database) *categoryRepository {
	return &categoryRepository{db: db}
}

type categoryRow struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	ParentID   *int   `db:"parent_id"`
	ParentName *string `db:"parent_name"`
	FilmCount  int    `db:"film_count"`
}

func (r *categoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	sqlq := `
SELECT
  c.id,
  c.name,
  p.id AS parent_id,
  p.name AS parent_name,
  COALESCE(COUNT(DISTINCT f.id), 0) AS film_count
FROM categories c
LEFT JOIN categories p ON p.id = c.parent_id
LEFT JOIN categories_films cf ON cf.category_id = c.id
LEFT JOIN films f ON f.id = cf.film_id AND f.deleted_at IS NULL
WHERE c.deleted_at IS NULL
GROUP BY c.id, c.name, p.id, p.name
ORDER BY c.id ASC
`

	var rows []categoryRow
	if err := r.db.SelectContext(ctx, &rows, sqlq); err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	out := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		var parent *domain.ParentCategory
		if row.ParentID != nil && row.ParentName != nil {
			parent = &domain.ParentCategory{ID: *row.ParentID, Name: *row.ParentName}
		}
		out = append(out, domain.Category{
			ID:             row.ID,
			Name:           row.Name,
			ParentCategory: parent,
			FilmCount:      row.FilmCount,
		})
	}
	return out, nil
}

