package models

type Movie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Name          string  `json:"name"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	ReleaseDate   string  `json:"release_date"`
	FirstAirDate  string  `json:"first_air_date"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	GenreIDs      []int   `json:"genre_ids"`
	Genres        []Genre `json:"genres"`
}

type DiscoverResponse struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

type MovieDetails struct {
	Title string `json:"title"`
}

type ExternalIDs struct {
	ID     int    `json:"id"`
	ImdbID string `json:"imdb_id"`
}

func (m *Movie) GetYear() string {
	date := m.ReleaseDate
	if date == "" {
		date = m.FirstAirDate
	}
	if len(date) >= 4 {
		return "(" + date[:4] + ")"
	}
	return ""
}

func (m *Movie) GetTitle() string {
	if m.Title != "" {
		return m.Title
	}
	return m.Name
}

func (m *Movie) GetGenreNames() []string {
	names := make([]string, len(m.Genres))
	for i, genre := range m.Genres {
		names[i] = genre.Name
	}
	return names
}