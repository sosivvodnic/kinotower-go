package user_repository

import (
	"context"
	"database/sql"
	"fmt"

	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
	"github.com/sosivvodnic/kinotower-go/internal/features/users/domain"
)

type Repository interface {
	GetUserByID(ctx context.Context, id int) (*domain.UserResponse, error)
	UpdateUser(ctx context.Context, id int, fio, email string, birthday *string, genderID int) error
	SoftDeleteUser(ctx context.Context, id int) error
}

type repository struct {
	db core_database.Database
}

func NewRepository(db core_database.Database) *repository {
	return &repository{db: db}
}

func (r *repository) GetUserByID(ctx context.Context, id int) (*domain.UserResponse, error) {
	sqlq := `
SELECT
  u.id,
  u.fio,
  u.email,
  TO_CHAR(u.birthday, 'YYYY-MM-DD') AS birthday,
  g.id AS gender_id,
  g.name AS gender_name,
  COALESCE(rc.review_count, 0) AS review_count,
  COALESCE(rt.rating_count, 0) AS rating_count
FROM users u
JOIN genders g ON g.id = u.gender_id
LEFT JOIN (
  SELECT user_id, COUNT(*) AS review_count
  FROM reviews
  WHERE deleted_at IS NULL
  GROUP BY user_id
) rc ON rc.user_id = u.id
LEFT JOIN (
  SELECT user_id, COUNT(*) AS rating_count
  FROM ratings
  GROUP BY user_id
) rt ON rt.user_id = u.id
WHERE u.id = $1 AND u.deleted_at IS NULL
LIMIT 1
`

	var out domain.UserResponse
	if err := r.db.GetContext(ctx, &out, sqlq, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &out, nil
}

func (r *repository) UpdateUser(ctx context.Context, id int, fio, email string, birthday *string, genderID int) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE users
SET fio=$1, email=$2, birthday=$3::date, gender_id=$4
WHERE id=$5 AND deleted_at IS NULL
`, fio, email, birthday, genderID, id)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *repository) SoftDeleteUser(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

