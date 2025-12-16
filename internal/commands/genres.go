package commands

import (
	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/sebastianneubert/tmdb/internal/models"
)

// LoadGenres fetches the genre list from TMDB API and returns both the list and a map for quick lookup
// Returns empty slices/maps if the API call fails (doesn't crash, just skips genre functionality)
func LoadGenres(client *api.Client) ([]models.Genre, map[string]int) {
	genreResp, err := client.GetGenres("de-DE")
	if err != nil {
		return []models.Genre{}, map[string]int{}
	}
	return genreResp.Genres, filters.BuildGenreMap(genreResp.Genres)
}
