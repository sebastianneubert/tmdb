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
	showsProviders string
	showsRegion    string
	showsMinRating float64
	showsMinVotes  int
	showsTimeout   int
)

var showsCmd = &cobra.Command{
	Use:   "shows",
	Short: "Find top-rated TV shows available on your streaming providers.",
	Long:  "Queries TMDb's Top Rated TV Shows list and checks streaming availability.",
	Run:   runShows,
}

func init() {
	showsCmd.Flags().StringVarP(&showsProviders, "providers", "p", config.DefaultProviders, "Comma-separated providers")
	showsCmd.Flags().StringVarP(&showsRegion, "region", "r", config.DefaultRegion, "Watch region")
	showsCmd.Flags().Float64Var(&showsMinRating, "min-rating", config.DefaultMinRating, "Minimum rating")
	showsCmd.Flags().IntVar(&showsMinVotes, "min-votes", config.DefaultMinVotes, "Minimum votes")
	showsCmd.Flags().IntVarP(&showsTimeout, "timeout", "T", config.DefaultTimeout, "Timeout in seconds")
}

func runShows(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	finalRegion := cfg.Region
	if cmd.Flags().Changed("region") {
		finalRegion = showsRegion
	}

	finalProviders := cfg.Providers
	if cmd.Flags().Changed("providers") {
		finalProviders = showsProviders
	}

	finalMinRating := cfg.MinRating
	if cmd.Flags().Changed("min-rating") {
		finalMinRating = showsMinRating
	}

	finalMinVotes := cfg.MinVotes
	if cmd.Flags().Changed("min-votes") {
		finalMinVotes = showsMinVotes
	}

	finalTimeout := cfg.Timeout
	if cmd.Flags().Changed("timeout") {
		finalTimeout = showsTimeout
	}

	client, err := api.NewClient(cfg.APIKey, finalTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	desiredProviders := filters.ParseProviders(finalProviders)

	fmt.Printf("Searching TMDb's Top Rated TV Shows...\n")
	fmt.Printf("Criteria: Min Rating: %.1f | Min Votes: %d\n", finalMinRating, finalMinVotes)
	fmt.Printf("Filtering for [%s] in region [%s]\n\n", finalProviders, strings.ToUpper(finalRegion))

	resultsFound := 0
	for page := 1; page <= config.MaxPagesToSearch && resultsFound < config.MaxResultsToDisplay; page++ {
		fmt.Printf("Fetching page %d...\n", page)

		resp, err := client.GetTopRatedShows(page, "de-DE")
		if err != nil {
			fmt.Printf("Warning: Failed to fetch page %d: %v\n", page, err)
			continue
		}

		for _, show := range resp.Results {
			if resultsFound >= config.MaxResultsToDisplay {
				break
			}

			if !filters.MeetsRatingCriteria(show.VoteAverage, show.VoteCount, finalMinRating, finalMinVotes) {
				continue
			}

			providerData, err := client.GetShowWatchProviders(show.ID, finalRegion)
			if err != nil {
				continue
			}

			availableProviders, isAvailable := filters.CheckAvailability(providerData, desiredProviders)
			if !isAvailable {
				continue
			}

			resultsFound++
			externalIDs, _ := client.GetShowExternalIDs(show.ID)
			englishTitle, _ := client.GetShowEnglishTitle(show.ID)
			if englishTitle == "" {
				englishTitle = show.OriginalName
			}

			display.DisplayShow(display.ShowDisplay{
				Number:       resultsFound,
				Title:        show.GetTitle(),
				EnglishTitle: englishTitle,
				Year:         show.GetYear(),
				Rating:       show.VoteAverage,
				Votes:        show.VoteCount,
				Providers:    availableProviders,
				TmdbID:       show.ID,
				ImdbID:       externalIDs.ImdbID,
				TvdbID:       externalIDs.TvdbID,
				Overview:     show.Overview,
			})
		}

		if page >= resp.TotalPages {
			break
		}
	}

	display.DisplaySeparator()
	if resultsFound == 0 {
		fmt.Println("No TV shows found matching criteria.")
	} else {
		fmt.Printf("Displayed %d top-rated TV shows.\n", resultsFound)
	}
}