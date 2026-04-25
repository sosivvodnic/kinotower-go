package domain

type Country struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	FilmCount int    `json:"filmCount" db:"film_count"`
}

type ListResponse struct {
	Countries []Country `json:"countries"`
}

