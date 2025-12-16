package display_test

import (
	"testing"

	"github.com/sebastianneubert/tmdb/internal/display"
	"github.com/sebastianneubert/tmdb/internal/models"
)

// MockFetcherAPIClient is a mock API client for DetailsFetcher testing
type MockFetcherAPIClient struct {
	ExternalIDsToReturn models.ExternalIDs
	EnglishTitleToReturn string
	ShouldFail          bool
}

func (m *MockFetcherAPIClient) GetExternalIDs(movieID int) (models.ExternalIDs, error) {
	if m.ShouldFail {
		return models.ExternalIDs{}, nil
	}
	if m.ExternalIDsToReturn.ImdbID != "" {
		return m.ExternalIDsToReturn, nil
	}
	return models.ExternalIDs{ImdbID: "tt0000001"}, nil
}

func (m *MockFetcherAPIClient) GetEnglishTitle(movieID int) (string, error) {
	if m.ShouldFail {
		return "", nil
	}
	return m.EnglishTitleToReturn, nil
}

func (m *MockFetcherAPIClient) GetRegionalTitle(movieID int, region string) (string, error) {
	if m.ShouldFail {
		return "", nil
	}
	return "Regional Title", nil
}

func (m *MockFetcherAPIClient) GetWatchProviders(movieID int, region string) (*models.WatchProviderResponse, error) {
	return &models.WatchProviderResponse{}, nil
}

func TestNewDetailsFetcher(t *testing.T) {
	mockClient := &MockFetcherAPIClient{}
	genres := []models.Genre{{ID: 1, Name: "Action"}}

	fetcher := display.NewDetailsFetcher(mockClient, "US", genres)

	if fetcher == nil {
		t.Fatal("NewDetailsFetcher should return a non-nil fetcher")
	}
}

func TestBuildMovieDisplaySimple(t *testing.T) {
	mockClient := &MockFetcherAPIClient{
		EnglishTitleToReturn: "The Matrix",
		ExternalIDsToReturn: models.ExternalIDs{
			ImdbID: "tt0133093",
		},
	}
	genres := []models.Genre{
		{ID: 1, Name: "Action"},
		{ID: 12, Name: "Adventure"},
	}

	fetcher := display.NewDetailsFetcher(mockClient, "US", genres)

	movie := &models.Movie{
		ID:          603,
		Title:       "The Matrix",
		VoteAverage: 8.7,
		VoteCount:   18500,
		ReleaseDate: "1999-03-31",
		Overview:    "A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers.",
	}

	providers := []string{"Netflix", "Amazon Prime"}
	genreNames := []string{"Action", "Adventure"}

	display := fetcher.BuildMovieDisplaySimple(1, movie, providers, genreNames)

	if display.Number != 1 {
		t.Errorf("Expected Number to be 1, got %d", display.Number)
	}
	if display.Title != "The Matrix" {
		t.Errorf("Expected Title to be 'The Matrix', got '%s'", display.Title)
	}
	if display.Rating != 8.7 {
		t.Errorf("Expected Rating to be 8.7, got %.1f", display.Rating)
	}
	if len(display.Providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(display.Providers))
	}
	if len(display.Genres) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(display.Genres))
	}
}

func TestBuildMovieDisplaySimpleWithEmptyDetails(t *testing.T) {
	mockClient := &MockFetcherAPIClient{
		EnglishTitleToReturn: "",
		ExternalIDsToReturn:  models.ExternalIDs{},
	}
	genres := []models.Genre{}

	fetcher := display.NewDetailsFetcher(mockClient, "US", genres)

	movie := &models.Movie{
		ID:          1,
		Title:       "Test Movie",
		VoteAverage: 7.0,
		VoteCount:   100,
		ReleaseDate: "2023-01-01",
	}

	display := fetcher.BuildMovieDisplaySimple(1, movie, []string{}, []string{})

	if display.Title != "Test Movie" {
		t.Errorf("Expected Title to be 'Test Movie', got '%s'", display.Title)
	}
	if display.Rating != 7.0 {
		t.Errorf("Expected Rating to be 7.0, got %.1f", display.Rating)
	}
}

func TestBuildMovieDisplaySimpleMultipleProviders(t *testing.T) {
	mockClient := &MockFetcherAPIClient{}
	genres := []models.Genre{}

	fetcher := display.NewDetailsFetcher(mockClient, "US", genres)

	movie := &models.Movie{
		ID:          1,
		Title:       "Test Movie",
		VoteAverage: 8.0,
		VoteCount:   1000,
	}

	providers := []string{"Netflix", "Amazon Prime", "Disney+", "Hulu", "HBO Max"}

	display := fetcher.BuildMovieDisplaySimple(1, movie, providers, []string{})

	if len(display.Providers) != 5 {
		t.Errorf("Expected 5 providers, got %d", len(display.Providers))
	}

	for i, p := range providers {
		if display.Providers[i] != p {
			t.Errorf("Expected provider[%d] to be '%s', got '%s'", i, p, display.Providers[i])
		}
	}
}

func TestBuildMovieDisplaySimpleGenreMapping(t *testing.T) {
	mockClient := &MockFetcherAPIClient{}
	genres := []models.Genre{
		{ID: 1, Name: "Action"},
		{ID: 12, Name: "Adventure"},
		{ID: 28, Name: "Animation"},
	}

	fetcher := display.NewDetailsFetcher(mockClient, "US", genres)

	movie := &models.Movie{
		ID:    1,
		Title: "Test Movie",
	}

	genreNames := []string{"Action", "Adventure"}
	display := fetcher.BuildMovieDisplaySimple(1, movie, []string{}, genreNames)

	if len(display.Genres) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(display.Genres))
	}

	expectedGenres := []string{"Action", "Adventure"}
	for i, g := range expectedGenres {
		if display.Genres[i] != g {
			t.Errorf("Expected genre[%d] to be '%s', got '%s'", i, g, display.Genres[i])
		}
	}
}

func TestBuildMovieDisplaySimpleReleaseYear(t *testing.T) {
	mockClient := &MockFetcherAPIClient{}
	genres := []models.Genre{}

	fetcher := display.NewDetailsFetcher(mockClient, "US", genres)

	movie := &models.Movie{
		ID:          1,
		Title:       "Test Movie",
		ReleaseDate: "2023-06-15",
	}

	display := fetcher.BuildMovieDisplaySimple(1, movie, []string{}, []string{})

	if display.Year != "(2023)" {
		t.Errorf("Expected Year to be '(2023)', got '%s'", display.Year)
	}
}

func TestBuildMovieDisplaySimpleWithVotes(t *testing.T) {
	mockClient := &MockFetcherAPIClient{}
	genres := []models.Genre{}

	fetcher := display.NewDetailsFetcher(mockClient, "US", genres)

	movie := &models.Movie{
		ID:          1,
		Title:       "Test Movie",
		VoteAverage: 8.5,
		VoteCount:   25000,
	}

	display := fetcher.BuildMovieDisplaySimple(1, movie, []string{}, []string{})

	if display.Rating != 8.5 {
		t.Errorf("Expected Rating to be 8.5, got %.1f", display.Rating)
	}
	if display.Votes != 25000 {
		t.Errorf("Expected Votes to be 25000, got %d", display.Votes)
	}
}
