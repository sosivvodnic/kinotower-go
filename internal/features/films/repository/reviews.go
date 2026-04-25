package film_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sosivvodnic/kinotower-go/internal/features/films/domain"
)

type reviewRow struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	FIO       string    `db:"fio"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *filmRepository) GetApprovedReviewsByFilmID(ctx context.Context, filmID int) ([]domain.Review, error) {
	sqlq := `
SELECT
  rv.id,
  u.id AS user_id,
  u.fio,
  rv.message,
  rv.created_at
FROM reviews rv
JOIN users u ON u.id = rv.user_id
WHERE rv.film_id = $1
  AND rv.is_approved = TRUE
  AND rv.deleted_at IS NULL
ORDER BY rv.created_at DESC
`

	var rows []reviewRow
	if err := r.db.SelectContext(ctx, &rows, sqlq, filmID); err != nil {
		return nil, fmt.Errorf("select reviews: %w", err)
	}

	out := make([]domain.Review, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Review{
			ID: row.ID,
			User: domain.ReviewUser{
				ID:  row.UserID,
				FIO: row.FIO,
			},
			Message:   row.Message,
			CreatedAt: domain.NewISOTime(row.CreatedAt),
		})
	}
	return out, nil
}

func (r *filmRepository) FilmExists(ctx context.Context, filmID int) (bool, error) {
	var cnt int
	if err := r.db.GetContext(ctx, &cnt, `SELECT COUNT(*) FROM films WHERE id=$1 AND deleted_at IS NULL`, filmID); err != nil {
		return false, fmt.Errorf("check film exists: %w", err)
	}
	return cnt > 0, nil
}

