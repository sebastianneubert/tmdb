package processor_test

import (
	"errors"
	"testing"

	"github.com/sebastianneubert/tmdb/internal/models"
	"github.com/sebastianneubert/tmdb/internal/processor"
)

// MockAPIClient is a mock implementation of the API client for testing
type MockAPIClient struct {
	Pages              int
	MoviesPerPage      int
	FailOnPage         int
	RatingsToReturn    []float64
	VotesToReturn      []int
	GenreIDsToReturn   [][]int
	WatchProvidersMap  map[int][]*models.Provider
	FilterByGenre      bool
	GenreIDFilterValue int
}

func TestMovieProcessorWithSuccessfulFetch(t *testing.T) {
	mockClient := &MockAPIClient{
		Pages:           2,
		MoviesPerPage:   2,
		RatingsToReturn: []float64{8.5, 7.5, 8.0, 7.0},
		VotesToReturn:   []int{1000, 500, 2000, 100},
		GenreIDsToReturn: [][]int{
			{1, 2},
			{1},
			{2, 3},
			{3},
		},
	}

	mp := processor.NewMovieProcessor(nil, processor.FilterConfig{
		MinRating: 7.0,
		MinVotes:  200,
		Region:    "US",
	})

	fetchCallCount := 0
	processCallCount := 0

	fetchFunc := func(page int) (*models.DiscoverResponse, error) {
		fetchCallCount++
		if page > mockClient.Pages {
			return &models.DiscoverResponse{
				Results:   []models.Movie{},
				TotalPages: mockClient.Pages,
			}, nil
		}

		startIdx := (page - 1) * mockClient.MoviesPerPage
		endIdx := startIdx + mockClient.MoviesPerPage
		if endIdx > len(mockClient.RatingsToReturn) {
			endIdx = len(mockClient.RatingsToReturn)
		}

		var movies []models.Movie
		for i := startIdx; i < endIdx; i++ {
			movies = append(movies, models.Movie{
				ID:          i + 1,
				Title:       "Movie " + string(rune(i+1)),
				VoteAverage: mockClient.RatingsToReturn[i],
				VoteCount:   mockClient.VotesToReturn[i],
				GenreIDs:    mockClient.GenreIDsToReturn[i],
				ReleaseDate: "2023-01-01",
				Overview:    "Overview for movie " + string(rune(i+1)),
			})
		}

		return &models.DiscoverResponse{
			Results:    movies,
			TotalPages: mockClient.Pages,
		}, nil
	}

	processFunc := func(movie *models.Movie, providers []string, genres []string) error {
		processCallCount++
		// Verify movie meets criteria
		if movie.VoteAverage < 7.0 || movie.VoteCount < 200 {
			t.Errorf("Movie passed to processFunc doesn't meet criteria: rating=%.1f, votes=%d", movie.VoteAverage, movie.VoteCount)
		}
		return nil
	}

	err := mp.Process(fetchFunc, processFunc)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if fetchCallCount == 0 {
		t.Error("fetchFunc was never called")
	}

	if processCallCount == 0 {
		t.Error("processFunc was never called for any movies")
	}

	// We expect 3 movies to pass filtering (8.5/1000, 8.0/2000, 7.5/500)
	expectedProcessCalls := 3
	if processCallCount != expectedProcessCalls {
		t.Errorf("Expected %d processFunc calls, got %d", expectedProcessCalls, processCallCount)
	}
}

func TestMovieProcessorWithFetchError(t *testing.T) {
	mp := processor.NewMovieProcessor(nil, processor.FilterConfig{
		MinRating: 7.0,
		MinVotes:  200,
		Region:    "US",
	})

	processCallCount := 0
	processFunc := func(movie *models.Movie, providers []string, genres []string) error {
		processCallCount++
		return nil
	}

	fetchCallCount := 0
	fetchFunc := func(page int) (*models.DiscoverResponse, error) {
		fetchCallCount++
		if page == 1 {
			return &models.DiscoverResponse{
				Results: []models.Movie{},
			}, errors.New("API error")
		}
		return &models.DiscoverResponse{}, nil
	}

	// Process continues even on fetch error (logs warning and moves to next page)
	err := mp.Process(fetchFunc, processFunc)
	if err != nil {
		t.Fatalf("Process should not return error: %v", err)
	}

	if processCallCount != 0 {
		t.Errorf("processFunc should not be called when fetch fails, but was called %d times", processCallCount)
	}
}

func TestMovieProcessorProcessFuncError(t *testing.T) {
	mp := processor.NewMovieProcessor(nil, processor.FilterConfig{
		MinRating: 7.0,
		MinVotes:  200,
		Region:    "US",
	})

	processCallCount := 0
	processFunc := func(movie *models.Movie, providers []string, genres []string) error {
		processCallCount++
		return errors.New("processing error")
	}

	fetchFunc := func(page int) (*models.DiscoverResponse, error) {
		return &models.DiscoverResponse{
			Results: []models.Movie{
				{
					ID:          1,
					Title:       "Test Movie",
					VoteAverage: 8.0,
					VoteCount:   1000,
				},
			},
			TotalPages: 1,
		}, nil
	}

	// Process continues even if processFunc returns error (continues to next movie)
	err := mp.Process(fetchFunc, processFunc)
	if err != nil {
		t.Fatalf("Process should not return error: %v", err)
	}

	if processCallCount != 1 {
		t.Errorf("processFunc should be called once, but was called %d times", processCallCount)
	}
}

func TestMovieProcessorRatingFilter(t *testing.T) {
	tests := []struct {
		minRating     float64
		movieRating   float64
		shouldProcess bool
	}{
		{7.0, 8.0, true},
		{7.0, 7.0, true},
		{7.0, 6.9, false},
		{8.5, 8.5, true},
		{8.5, 8.4, false},
	}

	for _, tt := range tests {
		mp := processor.NewMovieProcessor(nil, processor.FilterConfig{
			MinRating: tt.minRating,
			MinVotes:  0,
			Region:    "US",
		})

		processCallCount := 0
		processFunc := func(movie *models.Movie, providers []string, genres []string) error {
			processCallCount++
			return nil
		}

		fetchFunc := func(page int) (*models.DiscoverResponse, error) {
			return &models.DiscoverResponse{
				Results: []models.Movie{
					{
						ID:          1,
						Title:       "Test Movie",
						VoteAverage: tt.movieRating,
						VoteCount:   1000,
					},
				},
				TotalPages: 1,
			}, nil
		}

		mp.Process(fetchFunc, processFunc)

		if tt.shouldProcess && processCallCount != 1 {
			t.Errorf("MinRating=%.1f, MovieRating=%.1f: expected processFunc to be called, but wasn't", tt.minRating, tt.movieRating)
		}
		if !tt.shouldProcess && processCallCount != 0 {
			t.Errorf("MinRating=%.1f, MovieRating=%.1f: expected processFunc not to be called, but was", tt.minRating, tt.movieRating)
		}
	}
}

func TestMovieProcessorVotesFilter(t *testing.T) {
	tests := []struct {
		minVotes      int
		movieVotes    int
		shouldProcess bool
	}{
		{100, 200, true},
		{100, 100, true},
		{100, 99, false},
		{1000, 1000, true},
		{1000, 999, false},
	}

	for _, tt := range tests {
		mp := processor.NewMovieProcessor(nil, processor.FilterConfig{
			MinRating: 0,
			MinVotes:  tt.minVotes,
			Region:    "US",
		})

		processCallCount := 0
		processFunc := func(movie *models.Movie, providers []string, genres []string) error {
			processCallCount++
			return nil
		}

		fetchFunc := func(page int) (*models.DiscoverResponse, error) {
			return &models.DiscoverResponse{
				Results: []models.Movie{
					{
						ID:          1,
						Title:       "Test Movie",
						VoteAverage: 8.0,
						VoteCount:   tt.movieVotes,
					},
				},
				TotalPages: 1,
			}, nil
		}

		mp.Process(fetchFunc, processFunc)

		if tt.shouldProcess && processCallCount != 1 {
			t.Errorf("MinVotes=%d, MovieVotes=%d: expected processFunc to be called, but wasn't", tt.minVotes, tt.movieVotes)
		}
		if !tt.shouldProcess && processCallCount != 0 {
			t.Errorf("MinVotes=%d, MovieVotes=%d: expected processFunc not to be called, but was", tt.minVotes, tt.movieVotes)
		}
	}
}

func TestMovieProcessorEmptyResults(t *testing.T) {
	mp := processor.NewMovieProcessor(nil, processor.FilterConfig{
		MinRating: 7.0,
		MinVotes:  200,
		Region:    "US",
	})

	processCallCount := 0
	processFunc := func(movie *models.Movie, providers []string, genres []string) error {
		processCallCount++
		return nil
	}

	fetchFunc := func(page int) (*models.DiscoverResponse, error) {
		return &models.DiscoverResponse{
			Results:    []models.Movie{},
			TotalPages: 0,
		}, nil
	}

	err := mp.Process(fetchFunc, processFunc)
	if err != nil {
		t.Fatalf("Process should succeed with empty results: %v", err)
	}

	if processCallCount != 0 {
		t.Errorf("processFunc should not be called for empty results, but was called %d times", processCallCount)
	}
}
