package commands

import (
	"fmt"
	"time"

	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/display"
	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/sebastianneubert/tmdb/internal/models"
	"github.com/sebastianneubert/tmdb/internal/processor"
	"github.com/spf13/cobra"
)

var (
	topFlags = MovieCommandFlags{}
	topYear  int
)

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Find top-rated movies available on your streaming providers.",
	Long:  "Queries TMDb's Top Rated Movies list and checks streaming availability.",
	Run:   runTop,
}

func init() {
	topFlags.Register(topCmd, true)
	topFlags.RegisterAge(topCmd)
	topCmd.Flags().IntVarP(&topYear, "year", "y", 0, "Only movies released in this year (e.g. 2024)")
}

func runTop(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	if topYear != 0 && (topYear < 1888 || topYear > time.Now().Year()+1) {
		fmt.Printf("Error: invalid year %d\n", topYear)
		return
	}

	finalRegion, finalProviders, finalMinRating, finalMinVotes, finalTimeout, topGenre := topFlags.Resolve(cmd, cfg)
	maxAge, includeUnrated := topFlags.ResolveAge(cmd, cfg)

	client, err := api.NewClient(cfg.APIKey, finalTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	desiredProviders := filters.ParseProviders(finalProviders)
	genreList, genreMap := LoadGenres(client)

	title := "Top Rated Movies"
	if topYear > 0 {
		title = fmt.Sprintf("Top Rated Movies of %d", topYear)
	}
	display.PrintSearchStartMessage(title, finalMinRating, finalMinVotes, finalProviders, finalRegion)
	printAgeFilterNotice(maxAge, includeUnrated, finalRegion)

	ageFilter := processor.NewAgeFilter(client, finalRegion, maxAge, includeUnrated)

	processor := processor.NewMovieProcessor(client, processor.FilterConfig{
		MinRating:        finalMinRating,
		MinVotes:         finalMinVotes,
		Region:           finalRegion,
		GenreFilter:      topGenre,
		DesiredProviders: desiredProviders,
		GenreList:        genreList,
		GenreMap:         genreMap,
		AgeFilter:        ageFilter,
	})

	fetcher := display.NewDetailsFetcher(client, finalRegion, genreList)
	resultsFound := 0

	err = processor.Process(
		func(page int) (*models.DiscoverResponse, error) {
			if topYear > 0 {
				return client.GetTopRatedMoviesByYear(page, finalRegion, topYear, finalMinVotes)
			}
			return client.GetTopRatedMovies(page, finalRegion)
		},
		func(movie *models.Movie, providers []string, genres []string) error {
			resultsFound++
			movieDisplay := fetcher.BuildMovieDisplay(resultsFound, movie, providers, genres)
			movieDisplay.AgeRating = ageFilter.AgeLabel(movie.ID)
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
