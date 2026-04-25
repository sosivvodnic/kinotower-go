package domain

type FilmFilter struct {
	Page     int
	Size     int
	SortBy   string // name|year|rating
	SortDir  string // asc|desc
	Category int
	Country  int
	Search   *string
}
