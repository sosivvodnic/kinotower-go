package user_review_repository

import (
	"context"
	"fmt"
	"time"

	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	"github.com/sosivvodnic/kinotower-go/internal/features/user_reviews/domain"
)

type Repository interface {
	UserExists(ctx context.Context, userID int) (bool, error)
	FilmExists(ctx context.Context, filmID int) (bool, error)
	Create(ctx context.Context, userID int, filmID int, message string) (*domain.ReviewResponse, error)
	ListByUser(ctx context.Context, userID int) ([]domain.ReviewResponse, error)
	Delete(ctx context.Context, userID int, reviewID int) (bool, error)
}

type repository struct {
	db core_database.Database
}

func NewRepository(db core_database.Database) *repository {
	return &repository{db: db}
}

func (r *repository) UserExists(ctx context.Context, userID int) (bool, error) {
	var cnt int
	if err := r.db.GetContext(ctx, &cnt, `SELECT COUNT(*) FROM users WHERE id=$1 AND deleted_at IS NULL`, userID); err != nil {
		return false, fmt.Errorf("user exists: %w", err)
	}
	return cnt > 0, nil
}

func (r *repository) FilmExists(ctx context.Context, filmID int) (bool, error) {
	var cnt int
	if err := r.db.GetContext(ctx, &cnt, `SELECT COUNT(*) FROM films WHERE id=$1 AND deleted_at IS NULL`, filmID); err != nil {
		return false, fmt.Errorf("film exists: %w", err)
	}
	return cnt > 0, nil
}

type createRow struct {
	ID         int       `db:"id"`
	FilmID     int       `db:"film_id"`
	FilmName   string    `db:"film_name"`
	Message    string    `db:"message"`
	IsApproved bool      `db:"is_approved"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r *repository) Create(ctx context.Context, userID int, filmID int, message string) (*domain.ReviewResponse, error) {
	var row createRow
	err := r.db.GetContext(ctx, &row, `
INSERT INTO reviews (film_id, user_id, message, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING id, film_id, message, is_approved, created_at
`, filmID, userID, message)
	if err != nil {
		return nil, fmt.Errorf("create review: %w", err)
	}

	// film name
	if err := r.db.GetContext(ctx, &row.FilmName, `SELECT name FROM films WHERE id=$1`, filmID); err != nil {
		return nil, fmt.Errorf("get film name: %w", err)
	}

	return &domain.ReviewResponse{
		ID: row.ID,
		Film: domain.FilmShort{
			ID:   row.FilmID,
			Name: row.FilmName,
		},
		Message:    row.Message,
		IsApproved: row.IsApproved,
		CreatedAt:  domain.NewISOTime(row.CreatedAt),
	}, nil
}

type listRow struct {
	ID         int       `db:"id"`
	FilmID     int       `db:"film_id"`
	FilmName   string    `db:"film_name"`
	Message    string    `db:"message"`
	IsApproved bool      `db:"is_approved"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r *repository) ListByUser(ctx context.Context, userID int) ([]domain.ReviewResponse, error) {
	sqlq := `
SELECT
  rv.id,
  f.id AS film_id,
  f.name AS film_name,
  rv.message,
  rv.is_approved,
  rv.created_at
FROM reviews rv
JOIN films f ON f.id = rv.film_id
WHERE rv.user_id=$1 AND rv.deleted_at IS NULL
ORDER BY rv.created_at DESC
`
	var rows []listRow
	if err := r.db.SelectContext(ctx, &rows, sqlq, userID); err != nil {
		return nil, fmt.Errorf("list reviews: %w", err)
	}
	out := make([]domain.ReviewResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.ReviewResponse{
			ID: row.ID,
			Film: domain.FilmShort{
				ID:   row.FilmID,
				Name: row.FilmName,
			},
			Message:    row.Message,
			IsApproved: row.IsApproved,
			CreatedAt:  domain.NewISOTime(row.CreatedAt),
		})
	}
	return out, nil
}

func (r *repository) Delete(ctx context.Context, userID int, reviewID int) (bool, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE reviews SET deleted_at=NOW() WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL`, reviewID, userID)
	if err != nil {
		return false, fmt.Errorf("delete review: %w", err)
	}
	aff, _ := res.RowsAffected()
	return aff > 0, nil
}

