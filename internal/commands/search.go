package commands

import (
	"fmt"
	"strings"

	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/display"
	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/spf13/cobra"
)

var searchFlags = MovieCommandFlags{}

var (
	searchMaxResults int
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for movies by title.",
	Long: `Search for movies by title and show ratings and streaming availability.

Examples:
  tmdb search "star"
  tmdb search "star wars" --min-rating 7.0
  tmdb search "matrix" --region US --providers Netflix,Amazon --genre Action`,
	Args: cobra.MinimumNArgs(1),
	Run:  runSearch,
}

func init() {
	searchFlags.Register(searchCmd, true)
	searchCmd.Flags().IntVar(&searchMaxResults, "max", 20, "Maximum results to display")
}

func runSearch(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	query := strings.Join(args, " ")

	finalRegion, finalProviders, finalMinRating, finalMinVotes, finalTimeout, searchGenre := searchFlags.Resolve(cmd, cfg)

	client, err := api.NewClient(cfg.APIKey, finalTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	desiredProviders := filters.ParseProviders(finalProviders)
	genreList, genreMap := LoadGenres(client)

	fmt.Printf("🔍 Searching for: \"%s\"\n", query)
	fmt.Printf("Criteria: Min Rating: %.1f | Min Votes: %d\n", finalMinRating, finalMinVotes)
	fmt.Printf("Filtering for [%s] in region [%s]\n\n", finalProviders, strings.ToUpper(finalRegion))

	searchResp, err := client.SearchMovie(query, "de-DE", finalRegion)
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
		return
	}

	if len(searchResp.Results) == 0 {
		fmt.Printf("No movies found for \"%s\"\n", query)
		return
	}

	fmt.Printf("Found %d movies, filtering...\n\n", len(searchResp.Results))

	fetcher := display.NewDetailsFetcher(client, finalRegion, genreList)
	resultsFound := 0

	for _, movie := range searchResp.Results {
		if resultsFound >= searchMaxResults {
			break
		}

		// Use processor pattern for filtering
		m := movie

		// Apply rating and vote filters
		if !filters.MeetsRatingCriteria(m.VoteAverage, m.VoteCount, finalMinRating, finalMinVotes) {
			continue
		}

		// Check streaming availability
		providerData, err := client.GetWatchProviders(m.ID, finalRegion)
		if err != nil {
			continue
		}

		availableProviders, isAvailable := filters.CheckAvailability(providerData, desiredProviders)
		if !isAvailable {
			continue
		}

		if searchGenre != "" && !filters.FilterByGenre(&m, searchGenre, genreMap) {
			continue
		}

		// Movie matches all criteria
		resultsFound++

		genreNames := filters.GetGenreNames(m.GenreIDs, genreList)
		movieDisplay := fetcher.BuildMovieDisplaySimple(resultsFound, &m, availableProviders, genreNames)
		display.DisplayMovie(movieDisplay)
	}

	if resultsFound == 0 {
		display.PrintSearchNoResults(query, len(searchResp.Results), finalMinRating, finalMinVotes)
	} else {
		display.PrintSearchCompleteMessage(resultsFound, len(searchResp.Results))
	}
}
