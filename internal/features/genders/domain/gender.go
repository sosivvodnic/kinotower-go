package domain

type Gender struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type ListResponse struct {
	Genders []Gender `json:"genders"`
}
