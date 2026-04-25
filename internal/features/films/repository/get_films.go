package film_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
)

type filmRow struct {
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

func (r *filmRepository) GetFilms(ctx context.Context, f domain.FilmFilter) ([]domain.Film, int, error) {
	where := []string{"f.deleted_at IS NULL"}
	args := []any{}

	if f.Country > 0 {
		args = append(args, f.Country)
		where = append(where, fmt.Sprintf("f.country_id = $%d", len(args)))
	}
	if f.Category > 0 {
		args = append(args, f.Category)
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM categories_films cf WHERE cf.film_id=f.id AND cf.category_id=$%d)", len(args)))
	}
	if f.Search != nil && strings.TrimSpace(*f.Search) != "" {
		args = append(args, "%"+strings.TrimSpace(*f.Search)+"%")
		where = append(where, fmt.Sprintf("f.name ILIKE $%d", len(args)))
	}

	whereSQL := strings.Join(where, " AND ")

	// total
	var total int
	countSQL := "SELECT COUNT(*) FROM films f WHERE " + whereSQL
	if err := r.db.GetContext(ctx, &total, countSQL, args...); err != nil {
		return nil, 0, fmt.Errorf("count films: %w", err)
	}

	orderBy, err := filmsOrderBy(f.SortBy, f.SortDir)
	if err != nil {
		return nil, 0, err
	}

	// pagination
	limit := f.Size
	if limit <= 0 {
		limit = 10
	}
	offset := 0
	if f.Page > 1 {
		offset = (f.Page - 1) * limit
	}

	args = append(args, limit, offset)
	limitPos := len(args) - 1
	offsetPos := len(args)

	sql := `
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
WHERE ` + whereSQL + `
ORDER BY ` + orderBy + `
LIMIT $` + fmt.Sprintf("%d", limitPos) + ` OFFSET $` + fmt.Sprintf("%d", offsetPos) + `
`

	var rows []filmRow
	if err := r.db.SelectContext(ctx, &rows, sql, args...); err != nil {
		return nil, 0, fmt.Errorf("select films: %w", err)
	}

	films := make([]domain.Film, 0, len(rows))
	ids := make([]int, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
		films = append(films, domain.Film{
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
		})
	}

	if len(ids) == 0 {
		return films, total, nil
	}

	catsByFilm, err := r.getCategoriesByFilmIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	for i := range films {
		if cats, ok := catsByFilm[films[i].ID]; ok {
			films[i].Categories = cats
		}
	}

	return films, total, nil
}

func filmsOrderBy(sortBy, sortDir string) (string, error) {
	sb := strings.ToLower(strings.TrimSpace(sortBy))
	if sb == "" {
		sb = "name"
	}
	sd := strings.ToLower(strings.TrimSpace(sortDir))
	if sd == "" {
		sd = "asc"
	}
	if sd != "asc" && sd != "desc" {
		return "", fmt.Errorf("invalid sortDir")
	}
	switch sb {
	case "name":
		return "f.name " + sd, nil
	case "year":
		return "f.year_of_issue " + sd, nil
	case "rating":
		// Put NULL ratings last for better UX.
		return "ra.rating_avg " + sd + " NULLS LAST", nil
	default:
		return "", fmt.Errorf("invalid sortBy")
	}
}

type categoryRow struct {
	FilmID int    `db:"film_id"`
	ID     int    `db:"category_id"`
	Name   string `db:"category_name"`
}

func (r *filmRepository) getCategoriesByFilmIDs(ctx context.Context, filmIDs []int) (map[int][]domain.Category, error) {
	query, args, err := sqlx.In(`
SELECT
  cf.film_id,
  c.id AS category_id,
  c.name AS category_name
FROM categories_films cf
JOIN categories c ON c.id = cf.category_id
WHERE cf.film_id IN (?) AND c.deleted_at IS NULL
ORDER BY cf.id ASC
`, filmIDs)
	if err != nil {
		return nil, fmt.Errorf("build categories query: %w", err)
	}
	query = r.db.Rebind(query)

	var rows []categoryRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("select categories: %w", err)
	}

	out := map[int][]domain.Category{}
	for _, row := range rows {
		out[row.FilmID] = append(out[row.FilmID], domain.Category{
			ID:   row.ID,
			Name: row.Name,
		})
	}
	return out, nil
}
