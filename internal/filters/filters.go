package filters

import (
	"strconv"
	"strings"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func MeetsRatingCriteria(voteAverage float64, voteCount int, minRating float64, minVotes int) bool {
	return voteAverage >= minRating && voteCount >= minVotes
}

func ParseProviders(input string) map[string]bool {
	parts := strings.Split(input, ",")
	providerMap := make(map[string]bool)
	for _, part := range parts {
		key := strings.ToLower(strings.TrimSpace(part))
		providerMap[key] = true
	}
	return providerMap
}

func CheckAvailability(providerData models.RegionProviders, desiredProviders map[string]bool) ([]string, bool) {
	available := []string{}

	for _, p := range providerData.Flatrate {
		if isProviderMatched(p.ProviderName, desiredProviders) {
			available = append(available, p.ProviderName)
		}
	}

	return available, len(available) > 0
}

func isProviderMatched(providerName string, desiredProviders map[string]bool) bool {
	providerLower := strings.ToLower(providerName)

	for desired := range desiredProviders {
		if desired == "amazon" || desired == "amazonprime" {
			if strings.Contains(providerLower, "amazon prime video") {
				return true
			}
		} else if desired == providerLower {
			return true
		}
	}

	return false
}

func FilterByGenre(movie *models.Movie, desiredGenre string, genreMap map[string]int) bool {
	if desiredGenre == "" {
		return true // No filter
	}

	desiredGenreLower := strings.ToLower(strings.TrimSpace(desiredGenre))

	// Check if it's a genre ID (numeric)
	if genreID, err := strconv.Atoi(desiredGenreLower); err == nil {
		// Check by ID in GenreIDs array
		for _, id := range movie.GenreIDs {
			if id == genreID {
				return true
			}
		}
		// Check by ID in Genres array
		for _, genre := range movie.Genres {
			if genre.ID == genreID {
				return true
			}
		}
		return false
	}

	// Check by name in Genres array
	for _, genre := range movie.Genres {
		if strings.ToLower(genre.Name) == desiredGenreLower {
			return true
		}
	}

	// Check by name using genre map and GenreIDs
	if genreID, exists := genreMap[desiredGenreLower]; exists {
		for _, id := range movie.GenreIDs {
			if id == genreID {
				return true
			}
		}
	}

	return false
}

// BuildGenreMap creates a map of genre names to IDs for quick lookup
func BuildGenreMap(genres []models.Genre) map[string]int {
	genreMap := make(map[string]int)
	for _, genre := range genres {
		genreMap[strings.ToLower(genre.Name)] = genre.ID
	}
	return genreMap
}

// GetGenreNames converts genre IDs to names using a genre map
func GetGenreNames(genreIDs []int, genres []models.Genre) []string {
	idToName := make(map[int]string)
	for _, genre := range genres {
		idToName[genre.ID] = genre.Name
	}

	names := make([]string, 0, len(genreIDs))
	for _, id := range genreIDs {
		if name, exists := idToName[id]; exists {
			names = append(names, name)
		}
	}
	return names
}
