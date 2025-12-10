package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/display"
	"github.com/sebastianneubert/tmdb/internal/filters"
)

var (
	topProviders string
	topRegion    string
	topMinRating float64
	topMinVotes  int
	topTimeout   int
)

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Find top-rated movies available on your streaming providers.",
	Long:  "Queries TMDb's Top Rated Movies list and checks streaming availability.",
	Run:   runTop,
}

func init() {
	topCmd.Flags().StringVarP(&topProviders, "providers", "p", config.DefaultProviders, "Comma-separated providers")
	topCmd.Flags().StringVarP(&topRegion, "region", "r", config.DefaultRegion, "Watch region")
	topCmd.Flags().Float64Var(&topMinRating, "min-rating", config.DefaultMinRating, "Minimum rating")
	topCmd.Flags().IntVar(&topMinVotes, "min-votes", config.DefaultMinVotes, "Minimum votes")
	topCmd.Flags().IntVarP(&topTimeout, "timeout", "T", config.DefaultTimeout, "Timeout in seconds")
}

func runTop(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	finalRegion := cfg.Region
	if cmd.Flags().Changed("region") {
		finalRegion = topRegion
	}

	finalProviders := cfg.Providers
	if cmd.Flags().Changed("providers") {
		finalProviders = topProviders
	}

	finalMinRating := cfg.MinRating
	if cmd.Flags().Changed("min-rating") {
		finalMinRating = topMinRating
	}

	finalMinVotes := cfg.MinVotes
	if cmd.Flags().Changed("min-votes") {
		finalMinVotes = topMinVotes
	}

	finalTimeout := cfg.Timeout
	if cmd.Flags().Changed("timeout") {
		finalTimeout = topTimeout
	}

	client, err := api.NewClient(cfg.APIKey, finalTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	desiredProviders := filters.ParseProviders(finalProviders)

	fmt.Printf("Searching TMDb's Top Rated Movies...\n")
	fmt.Printf("Criteria: Min Rating: %.1f | Min Votes: %d\n", finalMinRating, finalMinVotes)
	fmt.Printf("Filtering for [%s] in region [%s]\n\n", finalProviders, strings.ToUpper(finalRegion))

	resultsFound := 0
	for page := 1; page <= config.MaxPagesToSearch && resultsFound < config.MaxResultsToDisplay; page++ {
		fmt.Printf("Fetching page %d...\n", page)

		resp, err := client.GetTopRatedMovies(page, finalRegion)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch page %d: %v\n", page, err)
			continue
		}

		for _, movie := range resp.Results {
			if resultsFound >= config.MaxResultsToDisplay {
				break
			}

			if !filters.MeetsRatingCriteria(movie.VoteAverage, movie.VoteCount, finalMinRating, finalMinVotes) {
				continue
			}

			providerData, err := client.GetWatchProviders(movie.ID, finalRegion)
			if err != nil {
				continue
			}

			availableProviders, isAvailable := filters.CheckAvailability(providerData, desiredProviders)
			if !isAvailable {
				continue
			}

			resultsFound++
			externalIDs, _ := client.GetExternalIDs(movie.ID)
			englishTitle, _ := client.GetEnglishTitle(movie.ID)
			if englishTitle == "" {
				englishTitle = movie.OriginalTitle
			}

			display.DisplayMovie(display.MovieDisplay{
				Number: resultsFound,
				Title: movie.GetTitle(),
				EnglishTitle: englishTitle,
				Year: movie.GetYear(),
				Rating: movie.VoteAverage,
				Votes: movie.VoteCount,
				Providers: availableProviders,
				TmdbID: movie.ID,
				ImdbID: externalIDs.ImdbID,
				Overview: movie.Overview,
			})
		}

		if page >= resp.TotalPages {
			break
		}
	}

	display.DisplaySeparator()
	if resultsFound == 0 {
		fmt.Println("No movies found matching criteria.")
	} else {
		fmt.Printf("Displayed %d top-rated movies.\n", resultsFound)
	}
}