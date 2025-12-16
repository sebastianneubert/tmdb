package commands

import (
	"fmt"

	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/display"
	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/sebastianneubert/tmdb/internal/models"
	"github.com/sebastianneubert/tmdb/internal/processor"
	"github.com/spf13/cobra"
)

var topFlags = MovieCommandFlags{}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Find top-rated movies available on your streaming providers.",
	Long:  "Queries TMDb's Top Rated Movies list and checks streaming availability.",
	Run:   runTop,
}

func init() {
	topFlags.Register(topCmd, true)
}

func runTop(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	finalRegion, finalProviders, finalMinRating, finalMinVotes, finalTimeout, topGenre := topFlags.Resolve(cmd, cfg)

	client, err := api.NewClient(cfg.APIKey, finalTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	desiredProviders := filters.ParseProviders(finalProviders)
	genreList, genreMap := LoadGenres(client)

	display.PrintSearchStartMessage("Top Rated Movies", finalMinRating, finalMinVotes, finalProviders, finalRegion)

	processor := processor.NewMovieProcessor(client, processor.FilterConfig{
		MinRating:        finalMinRating,
		MinVotes:         finalMinVotes,
		Region:           finalRegion,
		GenreFilter:      topGenre,
		DesiredProviders: desiredProviders,
		GenreList:        genreList,
		GenreMap:         genreMap,
	})

	fetcher := display.NewDetailsFetcher(client, finalRegion, genreList)
	resultsFound := 0

	err = processor.Process(
		func(page int) (*models.DiscoverResponse, error) {
			return client.GetTopRatedMovies(page, finalRegion)
		},
		func(movie *models.Movie, providers []string, genres []string) error {
			resultsFound++
			movieDisplay := fetcher.BuildMovieDisplay(resultsFound, movie, providers, genres)
			display.DisplayMovie(movieDisplay)
			return nil
		},
	)

	if err != nil {
		fmt.Printf("Error processing movies: %v\n", err)
		return
	}

	display.PrintSearchResultsSummary("top-rated movies", resultsFound)
}
