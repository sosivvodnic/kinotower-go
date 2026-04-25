package film_repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
)

type filmByIDRow struct {
	ID            int       `db:"id"`
	Name          string    `db:"name"`
	Duration      int       `db:"duration"`
	YearOfIssue   int       `db:"year_of_issue"`
	Age           int       `db:"age"`
	LinkImg       *string   `db:"link_img"`
	LinkKinopoisk *string   `db:"link_kinopoisk"`
	LinkVideo     string    `db:"link_video"`
	CreatedAt     time.Time `db:"created_at"`

	CountryID   int      `db:"country_id"`
	CountryName string   `db:"country_name"`
	RatingAvg   *float64 `db:"rating_avg"`
	ReviewCount int      `db:"review_count"`
}

func (r *filmRepository) GetFilmByID(ctx context.Context, id int) (*domain.Film, error) {
	sqlq := `
SELECT
  f.id,
  f.name,
  f.duration,
  f.year_of_issue,
  f.age,
  f.link_img,
  f.link_kinopoisk,
  f.link_video,
  f.created_at,
  c.id AS country_id,
  c.name AS country_name,
  ra.rating_avg,
  COALESCE(rc.review_count, 0) AS review_count
FROM films f
JOIN countries c ON c.id = f.country_id
LEFT JOIN (
  SELECT film_id, AVG(ball)::float8 AS rating_avg
  FROM ratings
  GROUP BY film_id
) ra ON ra.film_id = f.id
LEFT JOIN (
  SELECT film_id, COUNT(*) AS review_count
  FROM reviews
  WHERE is_approved = TRUE AND deleted_at IS NULL
  GROUP BY film_id
) rc ON rc.film_id = f.id
WHERE f.id = $1 AND f.deleted_at IS NULL
LIMIT 1
`

	var row filmByIDRow
	if err := r.db.GetContext(ctx, &row, sqlq, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get film: %w", err)
	}

	film := &domain.Film{
		ID:            row.ID,
		Name:          row.Name,
		Duration:      row.Duration,
		YearOfIssue:   row.YearOfIssue,
		Age:           row.Age,
		LinkImg:       row.LinkImg,
		LinkKinopoisk: row.LinkKinopoisk,
		LinkVideo:     row.LinkVideo,
		CreatedAt:     domain.NewISOTime(row.CreatedAt),
		Country: domain.Country{
			ID:   row.CountryID,
			Name: row.CountryName,
		},
		Categories:  []domain.Category{},
		RatingAvg:   row.RatingAvg,
		ReviewCount: row.ReviewCount,
	}

	catsByFilm, err := r.getCategoriesByFilmIDs(ctx, []int{film.ID})
	if err != nil {
		return nil, err
	}
	if cats, ok := catsByFilm[film.ID]; ok {
		film.Categories = cats
	}

	return film, nil
}
