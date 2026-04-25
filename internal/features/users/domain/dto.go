package domain

type Gender struct {
	ID   int    `json:"id" db:"gender_id"`
	Name string `json:"name" db:"gender_name"`
}

type UserResponse struct {
	ID          int     `json:"id" db:"id"`
	FIO         string  `json:"fio" db:"fio"`
	Email       string  `json:"email" db:"email"`
	Birthday    *string `json:"birthday" db:"birthday"` // yyyy-mm-dd
	Gender      Gender  `json:"gender"`
	ReviewCount int     `json:"reviewCount" db:"review_count"`
	RatingCount int     `json:"ratingCount" db:"rating_count"`
}

type UpdateUserRequest struct {
	FIO      string `json:"fio"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	GenderID int    `json:"gender_id"`
}
