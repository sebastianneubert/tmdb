package models

type Actor struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	KnownFor    []Movie `json:"known_for"`
	Popularity  float64 `json:"popularity"`
	ProfilePath string  `json:"profile_path"`
}

type ActorSearchResponse struct {
	Page         int     `json:"page"`
	Results      []Actor `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

type MovieCredit struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	ReleaseDate   string  `json:"release_date"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	Character     string  `json:"character"`
	Overview      string  `json:"overview"`
	// GenreIds      []Genre `json:"genre_ids"`
}

type ActorCreditsResponse struct {
	ID   int           `json:"id"`
	Cast []MovieCredit `json:"cast"`
}

func (mc *MovieCredit) GetYear() string {
	if len(mc.ReleaseDate) >= 4 {
		return "(" + mc.ReleaseDate[:4] + ")"
	}
	return ""
}