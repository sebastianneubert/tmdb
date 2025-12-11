package models

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GenreListResponse struct {
	Genres []Genre `json:"genres"`
}