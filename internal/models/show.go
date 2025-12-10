package models

type Show struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	OriginalName   string  `json:"original_name"`
	Overview       string  `json:"overview"`
	FirstAirDate   string  `json:"first_air_date"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
	OriginalLanguage string `json:"original_language"`
}

type ShowDiscoverResponse struct {
	Page         int    `json:"page"`
	Results      []Show `json:"results"`
	TotalPages   int    `json:"total_pages"`
	TotalResults int    `json:"total_results"`
}

type ShowDetails struct {
	Name string `json:"name"`
}

type ShowExternalIDs struct {
	ID     int    `json:"id"`
	ImdbID string `json:"imdb_id"`
	TvdbID int    `json:"tvdb_id"`
}

type ShowWatchProviderResponse struct {
	ID      int                        `json:"id"`
	Results map[string]RegionProviders `json:"results"`
}

func (s *Show) GetYear() string {
	if len(s.FirstAirDate) >= 4 {
		return "(" + s.FirstAirDate[:4] + ")"
	}
	return ""
}

func (s *Show) GetTitle() string {
	return s.Name
}