package commands

import (
	"fmt"
	"strings"

	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/display"
	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/sebastianneubert/tmdb/internal/models"
	"github.com/spf13/cobra"
)

var popularFlags = MovieCommandFlags{}

var popularCmd = &cobra.Command{
	Use:   "popular",
	Short: "Find popular movies available on your streaming providers.",
	Long:  "Queries TMDb's Popular Movies list and checks streaming availability.",
	Run:   runPopular,
}

func init() {
	popularFlags.Register(popularCmd, true)
}

func runPopular(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	finalRegion, finalProviders, finalMinRating, finalMinVotes, finalTimeout, popularGenre := popularFlags.Resolve(cmd, cfg)

	client, err := api.NewClient(cfg.APIKey, finalTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	desiredProviders := filters.ParseProviders(finalProviders)

	genreList, genreMap := LoadGenres(client)

	fmt.Printf("Searching TMDb's Popular Movies...\n")
	fmt.Printf("Criteria: Min Rating: %.1f | Min Votes: %d\n", finalMinRating, finalMinVotes)
	fmt.Printf("Filtering for [%s] in region [%s]\n\n", finalProviders, strings.ToUpper(finalRegion))

	processor := NewMovieProcessor(client, MovieFilterConfig{
		MinRating:        finalMinRating,
		MinVotes:         finalMinVotes,
		Region:           finalRegion,
		GenreFilter:      popularGenre,
		DesiredProviders: desiredProviders,
		GenreList:        genreList,
		GenreMap:         genreMap,
	})

	resultsFound := 0
	err = processor.Process(
		func(page int) (*models.DiscoverResponse, error) {
			return client.GetPopularMovies(page, finalRegion)
		},
		func(movie *models.Movie, providers []string, genres []string) error {
			resultsFound++
			externalIDs, _ := client.GetExternalIDs(movie.ID)
			englishTitle, _ := client.GetEnglishTitle(movie.ID)
			if englishTitle == "" {
				englishTitle = movie.OriginalTitle
			}

			// Get region-specific title (e.g., DE -> de-DE)
			languageCode := strings.ToLower(finalRegion) + "-" + strings.ToUpper(finalRegion)
			regionalTitle, _ := client.GetRegionalTitle(movie.ID, languageCode)
			if regionalTitle == "" {
				regionalTitle = movie.Title
			}

			display.DisplayMovie(display.MovieDisplay{
				Number:       resultsFound,
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
			})
			return nil
		},
	)

	if err != nil {
		fmt.Printf("Error processing movies: %v\n", err)
		return
	}

	display.DisplaySeparator()
	if resultsFound == 0 {
		fmt.Println("No movies found matching criteria.")
	} else {
		fmt.Printf("Displayed %d popular movies.\n", resultsFound)
	}
}
