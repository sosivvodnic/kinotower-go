package domain

import "github.com/sosivvodnic/kinotower-go/internal/core/httpjson"

type CreateRequest struct {
	FilmID int `json:"film_id"`
	Ball   int `json:"ball"`
}

type FilmShort struct {
	ID   int    `json:"id" db:"film_id"`
	Name string `json:"name" db:"film_name"`
}

type RatingResponse struct {
	ID        int              `json:"id" db:"id"`
	Film      FilmShort        `json:"film"`
	Score     int              `json:"score" db:"ball"`
	CreatedAt httpjson.ISOTime `json:"created_at" db:"created_at"`
}

type ListResponse struct {
	Ratings []RatingResponse `json:"ratings"`
}

type InvalidResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

