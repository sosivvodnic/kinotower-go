package user_rating_repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	"github.com/sosivvodnic/kinotower-go/internal/features/user_ratings/domain"
)

type Repository interface {
	UserExists(ctx context.Context, userID int) (bool, error)
	FilmExists(ctx context.Context, filmID int) (bool, error)
	RatingExists(ctx context.Context, userID int, filmID int) (bool, error)
	Create(ctx context.Context, userID int, filmID int, ball int) (*domain.RatingResponse, error)
	ListByUser(ctx context.Context, userID int) ([]domain.RatingResponse, error)
	Delete(ctx context.Context, userID int, ratingID int) (bool, error)
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

func (r *repository) RatingExists(ctx context.Context, userID int, filmID int) (bool, error) {
	var cnt int
	if err := r.db.GetContext(ctx, &cnt, `SELECT COUNT(*) FROM ratings WHERE user_id=$1 AND film_id=$2`, userID, filmID); err != nil {
		return false, fmt.Errorf("rating exists: %w", err)
	}
	return cnt > 0, nil
}

type createRow struct {
	ID        int       `db:"id"`
	FilmID    int       `db:"film_id"`
	Ball      int       `db:"ball"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *repository) Create(ctx context.Context, userID int, filmID int, ball int) (*domain.RatingResponse, error) {
	var row createRow
	err := r.db.GetContext(ctx, &row, `
INSERT INTO ratings (film_id, user_id, ball, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING id, film_id, ball, created_at
`, filmID, userID, ball)
	if err != nil {
		return nil, fmt.Errorf("create rating: %w", err)
	}

	var filmName string
	if err := r.db.GetContext(ctx, &filmName, `SELECT name FROM films WHERE id=$1`, filmID); err != nil {
		return nil, fmt.Errorf("get film name: %w", err)
	}

	return &domain.RatingResponse{
		ID: row.ID,
		Film: domain.FilmShort{
			ID:   row.FilmID,
			Name: filmName,
		},
		Score:     row.Ball,
		CreatedAt: domain.NewISOTime(row.CreatedAt),
	}, nil
}

type listRow struct {
	ID        int       `db:"id"`
	FilmID    int       `db:"film_id"`
	FilmName  string    `db:"film_name"`
	Ball      int       `db:"ball"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *repository) ListByUser(ctx context.Context, userID int) ([]domain.RatingResponse, error) {
	sqlq := `
SELECT
  rt.id,
  f.id AS film_id,
  f.name AS film_name,
  rt.ball,
  rt.created_at
FROM ratings rt
JOIN films f ON f.id = rt.film_id
WHERE rt.user_id=$1
ORDER BY rt.created_at DESC
`
	var rows []listRow
	if err := r.db.SelectContext(ctx, &rows, sqlq, userID); err != nil {
		return nil, fmt.Errorf("list ratings: %w", err)
	}
	out := make([]domain.RatingResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.RatingResponse{
			ID: row.ID,
			Film: domain.FilmShort{
				ID:   row.FilmID,
				Name: row.FilmName,
			},
			Score:     row.Ball,
			CreatedAt: domain.NewISOTime(row.CreatedAt),
		})
	}
	return out, nil
}

func (r *repository) Delete(ctx context.Context, userID int, ratingID int) (bool, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM ratings WHERE id=$1 AND user_id=$2`, ratingID, userID)
	if err != nil {
		return false, fmt.Errorf("delete rating: %w", err)
	}
	aff, _ := res.RowsAffected()
	return aff > 0, nil
}

var _ = sql.ErrNoRows

