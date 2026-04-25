package domain

import "github.com/sosivvodnic/kinotower-go/internal/core/httpjson"

type Country struct {
	ID   int    `json:"id" db:"country_id"`
	Name string `json:"name" db:"country_name"`
}

type Category struct {
	ID   int    `json:"id" db:"category_id"`
	Name string `json:"name" db:"category_name"`
}

type Film struct {
	ID            int              `json:"id" db:"id"`
	Name          string           `json:"name" db:"name"`
	Duration      int              `json:"duration" db:"duration"`
	YearOfIssue   int              `json:"year_of_issue" db:"year_of_issue"`
	Age           int              `json:"age" db:"age"`
	LinkImg       *string          `json:"link_img" db:"link_img"`
	LinkKinopoisk *string          `json:"link_kinopoisk" db:"link_kinopoisk"`
	LinkVideo     string           `json:"link_video" db:"link_video"`
	CreatedAt     httpjson.ISOTime `json:"created_at" db:"created_at"`

	Country     Country    `json:"country"`
	Categories  []Category `json:"categories"`
	RatingAvg   *float64   `json:"ratingAvg" db:"rating_avg"`
	ReviewCount int        `json:"reviewCount" db:"review_count"`
}

type FilmListResponse struct {
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int    `json:"total"`
	Films []Film `json:"films"`
}

type ReviewUser struct {
	ID  int    `json:"id" db:"user_id"`
	FIO string `json:"fio" db:"fio"`
}

type Review struct {
	ID        int              `json:"id" db:"id"`
	User      ReviewUser       `json:"user"`
	Message   string           `json:"message" db:"message"`
	CreatedAt httpjson.ISOTime `json:"created_at" db:"created_at"`
}
