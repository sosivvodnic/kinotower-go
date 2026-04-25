package auth_repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	core_database "github.com/sosivvodnic/kinotower-go/internal/core/database"
)

type User struct {
	ID        int       `db:"id"`
	FIO       string    `db:"fio"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Birthday  *time.Time `db:"birthday"`
	GenderID  int       `db:"gender_id"`
	CreatedAt time.Time `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Repository interface {
	CreateUser(ctx context.Context, fio, email, passwordHash string, birthday *time.Time, genderID int) (id int, createdFIO string, err error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db core_database.Database
}

func NewRepository(db core_database.Database) *repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, fio, email, passwordHash string, birthday *time.Time, genderID int) (int, string, error) {
	var id int
	var outFIO string
	err := r.db.QueryRowxContext(ctx, `
INSERT INTO users (fio, birthday, gender_id, email, password, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING id, fio
`, fio, birthday, genderID, email, passwordHash).Scan(&id, &outFIO)
	if err != nil {
		return 0, "", fmt.Errorf("create user: %w", err)
	}
	return id, outFIO, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.GetContext(ctx, &u, `
SELECT id, fio, email, password, birthday, gender_id, created_at, deleted_at
FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1
`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

