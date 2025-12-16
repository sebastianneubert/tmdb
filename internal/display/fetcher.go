package display

import (
	"strings"

	"github.com/sebastianneubert/tmdb/internal/models"
)

// DetailsFetcher handles fetching and assembling movie details for display
// APIClient is the subset of api.Client methods used by DetailsFetcher
type APIClient interface {
	GetExternalIDs(movieID int) (models.ExternalIDs, error)
	GetEnglishTitle(movieID int) (string, error)
	GetRegionalTitle(movieID int, language string) (string, error)
}

type DetailsFetcher struct {
	client    APIClient
	region    string
	genreList []models.Genre
}

// NewDetailsFetcher creates a new fetcher for movie details
func NewDetailsFetcher(client APIClient, region string, genreList []models.Genre) *DetailsFetcher {
	return &DetailsFetcher{
		client:    client,
		region:    region,
		genreList: genreList,
	}
}

// BuildMovieDisplay fetches all necessary details and returns a complete MovieDisplay struct
// with region-specific title and English title
func (df *DetailsFetcher) BuildMovieDisplay(number int, movie *models.Movie, providers []string, genres []string) MovieDisplay {
	externalIDs, _ := df.client.GetExternalIDs(movie.ID)

	englishTitle, _ := df.client.GetEnglishTitle(movie.ID)
	if englishTitle == "" {
		englishTitle = movie.OriginalTitle
	}

	languageCode := strings.ToLower(df.region) + "-" + strings.ToUpper(df.region)
	regionalTitle, _ := df.client.GetRegionalTitle(movie.ID, languageCode)
	if regionalTitle == "" {
		regionalTitle = movie.Title
	}

	return MovieDisplay{
		Number:       number,
		Title:        regionalTitle,
		EnglishTitle: englishTitle,
		Year:         movie.GetYear(),
		Rating:       movie.VoteAverage,
		Votes:        movie.VoteCount,
		Providers:    providers,
		TmdbID:       movie.ID,
		ImdbID:       externalIDs.ImdbID,
		Overview:     movie.Overview,
		Genres:       genres,
	}
}

// BuildMovieDisplaySimple is a simpler version that doesn't fetch region-specific titles
// (useful for commands like 'search' that might not need regional titles)
func (df *DetailsFetcher) BuildMovieDisplaySimple(number int, movie *models.Movie, providers []string, genres []string) MovieDisplay {
	externalIDs, _ := df.client.GetExternalIDs(movie.ID)

	englishTitle, _ := df.client.GetEnglishTitle(movie.ID)
	if englishTitle == "" {
		englishTitle = movie.OriginalTitle
	}

	return MovieDisplay{
		Number:       number,
		Title:        movie.GetTitle(),
		EnglishTitle: englishTitle,
		Year:         movie.GetYear(),
		Rating:       movie.VoteAverage,
		Votes:        movie.VoteCount,
		Providers:    providers,
		TmdbID:       movie.ID,
		ImdbID:       externalIDs.ImdbID,
		Overview:     movie.Overview,
		Genres:       genres,
	}
}
