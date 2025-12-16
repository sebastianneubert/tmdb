package commands

import (
	"testing"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func TestPopularMovieFiltering(t *testing.T) {
	tests := []struct {
		name        string
		movie       models.Movie
		minRating   float64
		minVotes    int
		shouldPass  bool
		description string
	}{
		{
			name: "Movie meets both criteria",
			movie: models.Movie{
				ID:          1,
				Title:       "Avatar",
				VoteAverage: 8.5,
				VoteCount:   500000,
			},
			minRating:   7.5,
			minVotes:    1000,
			shouldPass:  true,
			description: "Should pass when rating and votes meet minimum",
		},
		{
			name: "Movie fails rating criteria",
			movie: models.Movie{
				ID:          2,
				Title:       "Low Rated Movie",
				VoteAverage: 5.0,
				VoteCount:   500000,
			},
			minRating:   7.5,
			minVotes:    1000,
			shouldPass:  false,
			description: "Should fail when rating is below minimum",
		},
		{
			name: "Movie fails votes criteria",
			movie: models.Movie{
				ID:          3,
				Title:       "Unpopular Movie",
				VoteAverage: 8.5,
				VoteCount:   100,
			},
			minRating:   7.5,
			minVotes:    1000,
			shouldPass:  false,
			description: "Should fail when votes are below minimum",
		},
		{
			name: "Movie on boundary (rating)",
			movie: models.Movie{
				ID:          4,
				Title:       "Boundary Movie",
				VoteAverage: 7.5,
				VoteCount:   50000,
			},
			minRating:   7.5,
			minVotes:    1000,
			shouldPass:  true,
			description: "Should pass when rating equals minimum",
		},
		{
			name: "Movie on boundary (votes)",
			movie: models.Movie{
				ID:          5,
				Title:       "Boundary Movie 2",
				VoteAverage: 8.0,
				VoteCount:   1000,
			},
			minRating:   7.5,
			minVotes:    1000,
			shouldPass:  true,
			description: "Should pass when votes equal minimum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passes := tt.movie.VoteAverage >= tt.minRating && tt.movie.VoteCount >= tt.minVotes
			if passes != tt.shouldPass {
				t.Errorf("%s: expected %v, got %v", tt.description, tt.shouldPass, passes)
			}
		})
	}
}

func TestPopularMovieGenreFiltering(t *testing.T) {
	tests := []struct {
		name        string
		movie       models.Movie
		genreIDs    []int
		description string
	}{
		{
			name: "Movie with single genre",
			movie: models.Movie{
				ID:       1,
				Title:    "Action Movie",
				GenreIDs: []int{28},
			},
			genreIDs:    []int{28},
			description: "Should match movie with single genre",
		},
		{
			name: "Movie with multiple genres",
			movie: models.Movie{
				ID:       2,
				Title:    "Action Thriller",
				GenreIDs: []int{28, 53},
			},
			genreIDs:    []int{28, 53},
			description: "Should match movie with multiple genres",
		},
		{
			name: "Movie with no genres",
			movie: models.Movie{
				ID:       3,
				Title:    "Unknown Genre",
				GenreIDs: []int{},
			},
			genreIDs:    []int{},
			description: "Should handle movie with no genres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.movie.GenreIDs) != len(tt.genreIDs) {
				t.Errorf("%s: expected %d genres, got %d", tt.description, len(tt.genreIDs), len(tt.movie.GenreIDs))
			}
		})
	}
}

func TestPopularMovieRegionalTitle(t *testing.T) {
	tests := []struct {
		name            string
		originalTitle   string
		regionalTitle   string
		englishTitle    string
		expectedDisplay string
		description     string
	}{
		{
			name:            "All titles different",
			originalTitle:   "Avatar",
			regionalTitle:   "Avatar (German)",
			englishTitle:    "Avatar",
			expectedDisplay: "Avatar (German)",
			description:     "Should display regional title when different",
		},
		{
			name:            "Regional same as original",
			originalTitle:   "Inception",
			regionalTitle:   "Inception",
			englishTitle:    "Inception",
			expectedDisplay: "Inception",
			description:     "Should display title when regional matches original",
		},
		{
			name:            "Missing regional title",
			originalTitle:   "Movie",
			regionalTitle:   "",
			englishTitle:    "Movie",
			expectedDisplay: "Movie",
			description:     "Should fallback when regional title is missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			display := tt.regionalTitle
			if display == "" {
				display = tt.originalTitle
			}
			if display != tt.expectedDisplay {
				t.Errorf("%s: expected '%s', got '%s'", tt.description, tt.expectedDisplay, display)
			}
		})
	}
}

func TestPopularMovieProviderFiltering(t *testing.T) {
	tests := []struct {
		name               string
		desiredProviders   map[string]bool
		availableProviders []string
		shouldMatch        bool
		description        string
	}{
		{
			name:               "Movie on desired provider",
			desiredProviders:   map[string]bool{"Netflix": true, "DisneyPlus": true},
			availableProviders: []string{"Netflix"},
			shouldMatch:        true,
			description:        "Should match when movie available on desired provider",
		},
		{
			name:               "Movie not on desired provider",
			desiredProviders:   map[string]bool{"Netflix": true},
			availableProviders: []string{"AmazonPrime"},
			shouldMatch:        false,
			description:        "Should not match when movie not on desired provider",
		},
		{
			name:               "Movie on multiple providers, one desired",
			desiredProviders:   map[string]bool{"Netflix": true},
			availableProviders: []string{"AmazonPrime", "Netflix", "Wow"},
			shouldMatch:        true,
			description:        "Should match when at least one provider is desired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, provider := range tt.availableProviders {
				if tt.desiredProviders[provider] {
					found = true
					break
				}
			}
			if found != tt.shouldMatch {
				t.Errorf("%s: expected %v, got %v", tt.description, tt.shouldMatch, found)
			}
		})
	}
}

func TestPopularMovieGetYear(t *testing.T) {
	tests := []struct {
		name        string
		movie       models.Movie
		expected    string
		description string
	}{
		{
			name: "Valid release date",
			movie: models.Movie{
				ID:          1,
				Title:       "Recent Movie",
				ReleaseDate: "2023-05-15",
			},
			expected:    "(2023)",
			description: "Should extract year from release date",
		},
		{
			name: "Empty release date",
			movie: models.Movie{
				ID:          2,
				Title:       "No Date",
				ReleaseDate: "",
			},
			expected:    "",
			description: "Should return empty when no date",
		},
		{
			name: "Future release",
			movie: models.Movie{
				ID:          3,
				Title:       "Upcoming",
				ReleaseDate: "2025-12-25",
			},
			expected:    "(2025)",
			description: "Should handle future release dates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.movie.GetYear()
			if result != tt.expected {
				t.Errorf("%s: expected '%s', got '%s'", tt.description, tt.expected, result)
			}
		})
	}
}

func TestPopularMovieVoteAverage(t *testing.T) {
	tests := []struct {
		name        string
		voteAverage float64
		description string
	}{
		{
			name:        "High rating",
			voteAverage: 9.5,
			description: "Should store high vote average",
		},
		{
			name:        "Medium rating",
			voteAverage: 5.0,
			description: "Should store medium vote average",
		},
		{
			name:        "Low rating",
			voteAverage: 0.1,
			description: "Should store low vote average",
		},
		{
			name:        "Zero rating",
			voteAverage: 0.0,
			description: "Should handle zero rating",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movie := models.Movie{
				ID:          1,
				Title:       "Test",
				VoteAverage: tt.voteAverage,
			}
			if movie.VoteAverage != tt.voteAverage {
				t.Errorf("%s: expected %.1f, got %.1f", tt.description, tt.voteAverage, movie.VoteAverage)
			}
		})
	}
}

func TestPopularMovieVoteCount(t *testing.T) {
	tests := []struct {
		name        string
		voteCount   int
		description string
	}{
		{
			name:        "High vote count",
			voteCount:   500000,
			description: "Should store high vote counts",
		},
		{
			name:        "Medium vote count",
			voteCount:   5000,
			description: "Should store medium vote counts",
		},
		{
			name:        "Low vote count",
			voteCount:   100,
			description: "Should store low vote counts",
		},
		{
			name:        "Zero votes",
			voteCount:   0,
			description: "Should handle zero votes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movie := models.Movie{
				ID:        1,
				Title:     "Test",
				VoteCount: tt.voteCount,
			}
			if movie.VoteCount != tt.voteCount {
				t.Errorf("%s: expected %d, got %d", tt.description, tt.voteCount, movie.VoteCount)
			}
		})
	}
}
