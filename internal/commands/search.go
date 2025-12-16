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

	resultsFound := 0
	moviesChecked := 0

	for _, movie := range searchResp.Results {
		if resultsFound >= searchMaxResults {
			break
		}

		moviesChecked++

		// Apply rating and vote filters
		if !filters.MeetsRatingCriteria(movie.VoteAverage, movie.VoteCount, finalMinRating, finalMinVotes) {
			continue
		}

		// Check streaming availability
		providerData, err := client.GetWatchProviders(movie.ID, finalRegion)
		if err != nil {
			continue
		}

		availableProviders, isAvailable := filters.CheckAvailability(providerData, desiredProviders)
		if !isAvailable {
			continue
		}

		if searchGenre != "" && !filters.FilterByGenre(&movie, searchGenre, genreMap) {
			continue
		}

		// Movie matches all criteria
		resultsFound++

		// Get additional details
		externalIDs, _ := client.GetExternalIDs(movie.ID)
		englishTitle, _ := client.GetEnglishTitle(movie.ID)
		if englishTitle == "" {
			englishTitle = movie.OriginalTitle
		}

		genreNames := filters.GetGenreNames(movie.GenreIDs, genreList)

		display.DisplayMovie(display.MovieDisplay{
			Number:       resultsFound,
			Title:        movie.GetTitle(),
			EnglishTitle: englishTitle,
			Year:         movie.GetYear(),
			Rating:       movie.VoteAverage,
			Votes:        movie.VoteCount,
			Providers:    availableProviders,
			TmdbID:       movie.ID,
			ImdbID:       externalIDs.ImdbID,
			Overview:     movie.Overview,
			Genres:       genreNames,
		})
	}

	display.DisplaySeparator()
	if resultsFound == 0 {
		fmt.Printf("No movies found for \"%s\" that meet criteria and are available on your providers.\n", query)
		fmt.Printf("(Checked %d movies from search results)\n", moviesChecked)
		fmt.Println("\nTry:")
		fmt.Printf("  - Lowering --min-rating (current: %.1f)\n", finalMinRating)
		fmt.Printf("  - Lowering --min-votes (current: %d)\n", finalMinVotes)
		fmt.Printf("  - Adding more --providers\n")
	} else {
		fmt.Printf("Search complete: Displayed %d movies (out of %d checked).\n", resultsFound, moviesChecked)
	}
}
