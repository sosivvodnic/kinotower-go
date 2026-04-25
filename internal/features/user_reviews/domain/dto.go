package domain

import "github.com/sosivvodnic/kinotower-go/internal/core/httpjson"

type CreateRequest struct {
	FilmID  int    `json:"film_id"`
	Message string `json:"message"`
}

type FilmShort struct {
	ID   int    `json:"id" db:"film_id"`
	Name string `json:"name" db:"film_name"`
}

type ReviewResponse struct {
	ID         int              `json:"id" db:"id"`
	Film       FilmShort        `json:"film"`
	Message    string           `json:"message" db:"message"`
	IsApproved bool             `json:"is_approved" db:"is_approved"`
	CreatedAt  httpjson.ISOTime `json:"created_at" db:"created_at"`
}

type ListResponse struct {
	Reviews []ReviewResponse `json:"reviews"`
}

