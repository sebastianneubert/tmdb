package models

import (
	"testing"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func TestMovieGetYear(t *testing.T) {
	tests := []struct {
		name        string
		movie       models.Movie
		expected    string
		description string
	}{
		{
			name: "Valid ReleaseDate",
			movie: models.Movie{
				ID:          1,
				Title:       "Test Movie",
				ReleaseDate: "2023-05-15",
			},
			expected:    "(2023)",
			description: "Should extract year from release date",
		},
		{
			name: "Empty ReleaseDate with FirstAirDate",
			movie: models.Movie{
				ID:           1,
				Title:        "Test Show",
				ReleaseDate:  "",
				FirstAirDate: "2022-01-10",
			},
			expected:    "(2022)",
			description: "Should extract year from first air date when release date is empty",
		},
		{
			name: "Both empty dates",
			movie: models.Movie{
				ID:           1,
				Title:        "No Date Movie",
				ReleaseDate:  "",
				FirstAirDate: "",
			},
			expected:    "",
			description: "Should return empty string when both dates are empty",
		},
		{
			name: "ShortDate",
			movie: models.Movie{
				ID:          1,
				Title:       "Short Date",
				ReleaseDate: "99",
			},
			expected:    "",
			description: "Should return empty string when date is shorter than 4 chars",
		},
		{
			name: "Exactly 4 character date",
			movie: models.Movie{
				ID:          1,
				Title:       "Four Char Date",
				ReleaseDate: "2020",
			},
			expected:    "(2020)",
			description: "Should handle exactly 4 character dates",
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

func TestMovieGetTitle(t *testing.T) {
	tests := []struct {
		name        string
		movie       models.Movie
		expected    string
		description string
	}{
		{
			name: "Title provided",
			movie: models.Movie{
				ID:    1,
				Title: "Movie Title",
				Name:  "Show Name",
			},
			expected:    "Movie Title",
			description: "Should return Title when available",
		},
		{
			name: "Only Name provided",
			movie: models.Movie{
				ID:    1,
				Title: "",
				Name:  "Show Name",
			},
			expected:    "Show Name",
			description: "Should return Name when Title is empty",
		},
		{
			name: "Both empty",
			movie: models.Movie{
				ID:    1,
				Title: "",
				Name:  "",
			},
			expected:    "",
			description: "Should return empty string when both are empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.movie.GetTitle()
			if result != tt.expected {
				t.Errorf("%s: expected '%s', got '%s'", tt.description, tt.expected, result)
			}
		})
	}
}

func TestMovieGetGenreNames(t *testing.T) {
	tests := []struct {
		name        string
		movie       models.Movie
		expected    []string
		description string
	}{
		{
			name: "Multiple genres",
			movie: models.Movie{
				ID: 1,
				Genres: []models.Genre{
					{ID: 1, Name: "Action"},
					{ID: 2, Name: "Thriller"},
					{ID: 3, Name: "Drama"},
				},
			},
			expected:    []string{"Action", "Thriller", "Drama"},
			description: "Should extract all genre names",
		},
		{
			name: "Single genre",
			movie: models.Movie{
				ID: 1,
				Genres: []models.Genre{
					{ID: 1, Name: "Comedy"},
				},
			},
			expected:    []string{"Comedy"},
			description: "Should handle single genre",
		},
		{
			name: "No genres",
			movie: models.Movie{
				ID:     1,
				Genres: []models.Genre{},
			},
			expected:    []string{},
			description: "Should return empty slice when no genres",
		},
		{
			name: "Nil genres",
			movie: models.Movie{
				ID:     1,
				Genres: nil,
			},
			expected:    []string{},
			description: "Should handle nil genres slice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.movie.GetGenreNames()
			if len(result) != len(tt.expected) {
				t.Errorf("%s: expected %d genres, got %d", tt.description, len(tt.expected), len(result))
				return
			}
			for i, name := range result {
				if name != tt.expected[i] {
					t.Errorf("%s: expected genre[%d] '%s', got '%s'", tt.description, i, tt.expected[i], name)
				}
			}
		})
	}
}

func TestMovieCharacterField(t *testing.T) {
	movie := models.Movie{
		ID:        1,
		Title:     "Toy Story",
		Character: "Woody",
	}

	if movie.Character != "Woody" {
		t.Errorf("Expected Character 'Woody', got '%s'", movie.Character)
	}
}

func TestMovieStructFields(t *testing.T) {
	movie := models.Movie{
		ID:            1,
		Title:         "The Matrix",
		Name:          "",
		OriginalTitle: "The Matrix",
		Overview:      "A hacker learns about the true nature of reality",
		ReleaseDate:   "1999-03-31",
		FirstAirDate:  "",
		VoteAverage:   8.7,
		VoteCount:     1000000,
		GenreIDs:      []int{28, 878},
		Genres:        []models.Genre{{ID: 28, Name: "Action"}},
		Character:     "",
	}

	if movie.ID != 1 {
		t.Errorf("Expected ID 1, got %d", movie.ID)
	}

	if movie.Title != "The Matrix" {
		t.Errorf("Expected Title 'The Matrix', got '%s'", movie.Title)
	}

	if movie.VoteAverage != 8.7 {
		t.Errorf("Expected VoteAverage 8.7, got %f", movie.VoteAverage)
	}

	if movie.VoteCount != 1000000 {
		t.Errorf("Expected VoteCount 1000000, got %d", movie.VoteCount)
	}

	if len(movie.GenreIDs) != 2 {
		t.Errorf("Expected 2 GenreIDs, got %d", len(movie.GenreIDs))
	}
}
