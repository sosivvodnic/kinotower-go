package core_database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // добавьте импорт драйвера
)

type Database struct {
	*sqlx.DB
}

func NewDatabase() (*Database, error) {
	cfg := NewConfigMust() // ← исправлено название

	// Port теперь int, используем %d
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{DB: db}, nil
}
