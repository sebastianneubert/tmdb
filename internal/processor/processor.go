package processor

import (
	"fmt"

	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/sebastianneubert/tmdb/internal/models"
)

// FilterConfig holds all configuration needed for processing and filtering movies
type FilterConfig struct {
	MinRating        float64
	MinVotes         int
	Region           string
	GenreFilter      string
	DesiredProviders map[string]bool
	GenreList        []models.Genre
	GenreMap         map[string]int
}

// MovieProcessor handles fetching, filtering, and processing movies
type MovieProcessor struct {
	client *api.Client
	config FilterConfig
}

// NewMovieProcessor creates a new MovieProcessor instance
func NewMovieProcessor(client *api.Client, config FilterConfig) *MovieProcessor {
	return &MovieProcessor{
		client: client,
		config: config,
	}
}

// ProcessMovieFunc is the callback function type for processing each movie that passes filters
// It receives the filtered movie, available providers, and genre names
type ProcessMovieFunc func(*models.Movie, []string, []string) error

// FetchFunc is the callback function type for fetching a page of movies from the API
type FetchFunc func(page int) (*models.DiscoverResponse, error)

// Process fetches movies page by page, applies all filters, and calls processFunc for each matching movie
// The apiCall parameter allows different API endpoints (top-rated, popular, search, etc.)
// The processFunc parameter allows different display/processing logic per command
func (mp *MovieProcessor) Process(apiCall FetchFunc, processFunc ProcessMovieFunc) error {
	resultsFound := 0

	for page := 1; page <= config.MaxPagesToSearch && resultsFound < config.MaxResultsToDisplay; page++ {
		fmt.Printf("Fetching page %d...\n", page)

		resp, err := apiCall(page)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch page %d: %v\n", page, err)
			continue
		}

		for _, movie := range resp.Results {
			if resultsFound >= config.MaxResultsToDisplay {
				break
			}

			// Apply rating and vote filters
			if !filters.MeetsRatingCriteria(movie.VoteAverage, movie.VoteCount, mp.config.MinRating, mp.config.MinVotes) {
				continue
			}

			// Apply genre filter
			if mp.config.GenreFilter != "" && !filters.FilterByGenre(&movie, mp.config.GenreFilter, mp.config.GenreMap) {
				continue
			}

			// Check streaming availability. If no client is provided (e.g. in tests),
			// assume availability so tests can focus on filtering logic.
			var availableProviders []string
			var isAvailable bool
			if mp.client == nil {
				availableProviders = []string{}
				isAvailable = true
			} else {
				providerData, err := mp.client.GetWatchProviders(movie.ID, mp.config.Region)
				if err != nil {
					continue
				}
				availableProviders, isAvailable = filters.CheckAvailability(providerData, mp.config.DesiredProviders)
				if !isAvailable {
					continue
				}
			}

			// Movie passed all filters
			resultsFound++
			genreNames := filters.GetGenreNames(movie.GenreIDs, mp.config.GenreList)

			// Call the processing function with filtered results
			if err := processFunc(&movie, availableProviders, genreNames); err != nil {
				continue
			}
		}

		if page >= resp.TotalPages {
			break
		}
	}

	return nil
}
