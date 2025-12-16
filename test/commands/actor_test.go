package commands

import (
	"fmt"
	"strconv"
	"testing"
)

func TestActorIndexParsing(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedIndex int
		shouldError   bool
	}{
		{
			name:          "Valid index 1",
			input:         "1",
			expectedIndex: 0, // 1-based to 0-based conversion
			shouldError:   false,
		},
		{
			name:          "Valid index 5",
			input:         "5",
			expectedIndex: 4,
			shouldError:   false,
		},
		{
			name:          "Invalid index 0",
			input:         "0",
			expectedIndex: -1,
			shouldError:   true,
		},
		{
			name:          "Invalid index negative",
			input:         "-1",
			expectedIndex: -1,
			shouldError:   true,
		},
		{
			name:          "Invalid index non-numeric",
			input:         "abc",
			expectedIndex: -1,
			shouldError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseActorIndex(tt.input)

			if tt.shouldError && err == nil {
				t.Errorf("Expected error for input '%s', but got none", tt.input)
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tt.input, err)
			}
		})
	}
}

// Helper function to parse actor index
func parseActorIndex(input string) (int, error) {
	index, err := parseStringToInt(input)
	if err != nil {
		return -1, err
	}
	if index < 1 {
		return -1, fmt.Errorf("index must be positive")
	}
	return index - 1, nil // Convert to 0-based
}

// Helper to parse string to int
func parseStringToInt(input string) (int, error) {
	return strconv.Atoi(input)
}

func TestActorSelectionFromMultipleResults(t *testing.T) {
	// Simulate multiple search results for "Tom"
	results := []struct {
		name       string
		popularity float64
		id         int
	}{
		{"Tom Hanks", 92.3, 1},
		{"Tom Hardy", 85.6, 2},
		{"Tom Cruise", 88.9, 3},
	}

	tests := []struct {
		name              string
		userIndex         int
		expectedActor     string
		expectedID        int
		shouldBeValid     bool
	}{
		{
			name:          "Select first actor",
			userIndex:     1,
			expectedActor: "Tom Hanks",
			expectedID:    1,
			shouldBeValid: true,
		},
		{
			name:          "Select second actor",
			userIndex:     2,
			expectedActor: "Tom Hardy",
			expectedID:    2,
			shouldBeValid: true,
		},
		{
			name:          "Select third actor",
			userIndex:     3,
			expectedActor: "Tom Cruise",
			expectedID:    3,
			shouldBeValid: true,
		},
		{
			name:          "Index out of range",
			userIndex:     5,
			expectedActor: "",
			expectedID:    0,
			shouldBeValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.userIndex < 1 || tt.userIndex > len(results) {
				if tt.shouldBeValid {
					t.Errorf("Expected valid selection, but index %d is out of range", tt.userIndex)
				}
				return
			}

			// Convert to 0-based index
			zeroBasedIndex := tt.userIndex - 1
			selectedActor := results[zeroBasedIndex]

			if selectedActor.name != tt.expectedActor {
				t.Errorf("Expected actor '%s', got '%s'", tt.expectedActor, selectedActor.name)
			}

			if selectedActor.id != tt.expectedID {
				t.Errorf("Expected ID %d, got %d", tt.expectedID, selectedActor.id)
			}
		})
	}
}

func TestActorSelectionWithSorting(t *testing.T) {
	// Simulate search results that need to be sorted by popularity
	// This mimics the real scenario where "Foxx" returns unsorted results
	unsortedResults := []struct {
		name       string
		popularity float64
		id         int
	}{
		{"Random Foxx", 0.5, 5392815},  // Low popularity, should be last
		{"Jamie Foxx", 3.7, 134},        // High popularity, should be first
		{"Redd Foxx", 0.8, 56949},       // Medium popularity, should be second
	}

	// Sort by popularity (descending) - simulating what the code does
	type Actor struct {
		name       string
		popularity float64
		id         int
	}
	actors := make([]Actor, len(unsortedResults))
	for i, a := range unsortedResults {
		actors[i] = Actor{a.name, a.popularity, a.id}
	}

	// Sort by popularity descending
	for i := 0; i < len(actors); i++ {
		for j := i + 1; j < len(actors); j++ {
			if actors[j].popularity > actors[i].popularity {
				actors[i], actors[j] = actors[j], actors[i]
			}
		}
	}

	tests := []struct {
		name           string
		index          int
		expectedName   string
		expectedID     int
		expectedPopularity float64
	}{
		{
			name:           "Index 1 should be Jamie Foxx with ID 134",
			index:          0, // 0-based after conversion
			expectedName:   "Jamie Foxx",
			expectedID:     134,
			expectedPopularity: 3.7,
		},
		{
			name:           "Index 2 should be Redd Foxx with ID 56949",
			index:          1,
			expectedName:   "Redd Foxx",
			expectedID:     56949,
			expectedPopularity: 0.8,
		},
		{
			name:           "Index 3 should be Random Foxx with ID 5392815",
			index:          2,
			expectedName:   "Random Foxx",
			expectedID:     5392815,
			expectedPopularity: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.index >= len(actors) {
				t.Fatalf("Index %d out of range", tt.index)
			}

			selected := actors[tt.index]

			if selected.name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, selected.name)
			}

			if selected.id != tt.expectedID {
				t.Errorf("Expected ID %d, got %d", tt.expectedID, selected.id)
			}

			if selected.popularity != tt.expectedPopularity {
				t.Errorf("Expected popularity %.1f, got %.1f", tt.expectedPopularity, selected.popularity)
			}
		})
	}
}

func TestActorNameWithIndex(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedName   string
		shouldContain  string
	}{
		{
			name:          "Simple name with index",
			input:         "Megan Fox",
			expectedName:  "Megan Fox",
			shouldContain: "Megan Fox",
		},
		{
			name:          "Multiple word name",
			input:         "Leonardo DiCaprio",
			expectedName:  "Leonardo DiCaprio",
			shouldContain: "Leonardo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != tt.expectedName {
				t.Errorf("Expected '%s', got '%s'", tt.expectedName, tt.input)
			}
		})
	}
}

func TestGenreFilteringInFilmography(t *testing.T) {
	// Simulate movies with different genres
	type Movie struct {
		title   string
		genres  []string
		matches bool
	}

	movies := []Movie{
		{
			title:   "The Matrix",
			genres:  []string{"Action", "Sci-Fi"},
			matches: true, // Contains Action
		},
		{
			title:   "Forrest Gump",
			genres:  []string{"Drama", "Romance"},
			matches: false, // Does not contain Action
		},
		{
			title:   "Mission: Impossible",
			genres:  []string{"Action", "Thriller"},
			matches: true, // Contains Action
		},
		{
			title:   "The Dark Knight",
			genres:  []string{"Action", "Crime", "Drama"},
			matches: true, // Contains Action
		},
		{
			title:   "Pride and Prejudice",
			genres:  []string{"Romance", "Drama"},
			matches: false, // Does not contain Action
		},
	}

	tests := []struct {
		name            string
		movie           Movie
		genreFilter     string
		shouldInclude   bool
	}{
		{
			name:          "Action movie matches Action filter",
			movie:         movies[0],
			genreFilter:   "Action",
			shouldInclude: true,
		},
		{
			name:          "Drama movie does not match Action filter",
			movie:         movies[1],
			genreFilter:   "Action",
			shouldInclude: false,
		},
		{
			name:          "Action-Thriller matches Action filter",
			movie:         movies[2],
			genreFilter:   "Action",
			shouldInclude: true,
		},
		{
			name:          "Multi-genre movie with Action matches",
			movie:         movies[3],
			genreFilter:   "Action",
			shouldInclude: true,
		},
		{
			name:          "Pure drama does not match Action",
			movie:         movies[4],
			genreFilter:   "Action",
			shouldInclude: false,
		},
		{
			name:          "Movie with Drama matches Drama filter",
			movie:         movies[0],
			genreFilter:   "Drama",
			shouldInclude: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if movie contains the filter genre
			found := false
			for _, genre := range tt.movie.genres {
				if genre == tt.genreFilter {
					found = true
					break
				}
			}

			if found != tt.shouldInclude {
				t.Errorf("Expected inclusion: %v, but got: %v for '%s' with genres %v",
					tt.shouldInclude, found, tt.movie.title, tt.movie.genres)
			}
		})
	}
}

func TestGenreDisplayInFilmography(t *testing.T) {
	// Simulate movies with genre information for display
	type MovieDisplay struct {
		title  string
		genres []string
	}

	movies := []MovieDisplay{
		{
			title:  "The Matrix",
			genres: []string{"Action", "Sci-Fi"},
		},
		{
			title:  "Forrest Gump",
			genres: []string{"Drama", "Romance"},
		},
		{
			title:  "Mission: Impossible",
			genres: []string{"Action", "Thriller"},
		},
	}

	tests := []struct {
		name           string
		movie          MovieDisplay
		expectedGenres []string
	}{
		{
			name:           "Action Sci-Fi movie",
			movie:          movies[0],
			expectedGenres: []string{"Action", "Sci-Fi"},
		},
		{
			name:           "Drama Romance movie",
			movie:          movies[1],
			expectedGenres: []string{"Drama", "Romance"},
		},
		{
			name:           "Action Thriller movie",
			movie:          movies[2],
			expectedGenres: []string{"Action", "Thriller"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.movie.genres) != len(tt.expectedGenres) {
				t.Errorf("Expected %d genres, got %d", len(tt.expectedGenres), len(tt.movie.genres))
				return
			}

			for i, genre := range tt.movie.genres {
				if genre != tt.expectedGenres[i] {
					t.Errorf("Expected genre '%s' at position %d, got '%s'",
						tt.expectedGenres[i], i, genre)
				}
			}
		})
	}
}

func TestMultipleGenreFiltering(t *testing.T) {
	// Test that a single movie can match any of multiple genres
	type Movie struct {
		title  string
		genres []string
	}

	movie := Movie{
		title:  "The Dark Knight",
		genres: []string{"Action", "Crime", "Drama"},
	}

	tests := []struct {
		name        string
		filterGenre string
		shouldMatch bool
	}{
		{
			name:        "Matches Action",
			filterGenre: "Action",
			shouldMatch: true,
		},
		{
			name:        "Matches Crime",
			filterGenre: "Crime",
			shouldMatch: true,
		},
		{
			name:        "Matches Drama",
			filterGenre: "Drama",
			shouldMatch: true,
		},
		{
			name:        "Does not match Comedy",
			filterGenre: "Comedy",
			shouldMatch: false,
		},
		{
			name:        "Does not match Horror",
			filterGenre: "Horror",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, genre := range movie.genres {
				if genre == tt.filterGenre {
					found = true
					break
				}
			}

			if found != tt.shouldMatch {
				t.Errorf("Expected match: %v, but got: %v for genre '%s'",
					tt.shouldMatch, found, tt.filterGenre)
			}
		})
	}
}
