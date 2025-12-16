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

type ActorCreditsResponse struct {
	ID   int     `json:"id"`
	Cast []Movie `json:"cast"`
}
