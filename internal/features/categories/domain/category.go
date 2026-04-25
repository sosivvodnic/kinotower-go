package domain

type ParentCategory struct {
	ID   int    `json:"id" db:"parent_id"`
	Name string `json:"name" db:"parent_name"`
}

type Category struct {
	ID             int             `json:"id" db:"id"`
	Name           string          `json:"name" db:"name"`
	ParentCategory *ParentCategory `json:"parentCategory"`
	FilmCount      int             `json:"filmCount" db:"film_count"`
}

type ListResponse struct {
	Categories []Category `json:"categories"`
}

