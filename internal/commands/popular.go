package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/display"
	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/sebastianneubert/tmdb/internal/models"
)

var (
	popularProviders string
	popularRegion    string
	popularMinRating float64
	popularMinVotes  int
	popularTimeout   int
	popularGenre     string
)

var popularCmd = &cobra.Command{
	Use:   "popular",
	Short: "Find popular movies available on your streaming providers.",
	Long:  "Queries TMDb's Popular Movies list and checks streaming availability.",
	Run:   runPopular,
}

func init() {
	popularCmd.Flags().StringVarP(&popularProviders, "providers", "p", config.DefaultProviders, "Comma-separated providers")
	popularCmd.Flags().StringVarP(&popularRegion, "region", "r", config.DefaultRegion, "Watch region")
	popularCmd.Flags().Float64Var(&popularMinRating, "min-rating", config.DefaultMinRating, "Minimum rating")
	popularCmd.Flags().IntVar(&popularMinVotes, "min-votes", config.DefaultMinVotes, "Minimum votes")
	popularCmd.Flags().IntVarP(&popularTimeout, "timeout", "T", config.DefaultTimeout, "Timeout in seconds")
	popularCmd.Flags().StringVar(&popularGenre, "genre", "", "Filter by genre (name or ID)")
}

func runPopular(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	finalRegion := cfg.Region
	if cmd.Flags().Changed("region") {
		finalRegion = popularRegion
	}

	finalProviders := cfg.Providers
	if cmd.Flags().Changed("providers") {
		finalProviders = popularProviders
	}

	finalMinRating := cfg.MinRating
	if cmd.Flags().Changed("min-rating") {
		finalMinRating = popularMinRating
	}

	finalMinVotes := cfg.MinVotes
	if cmd.Flags().Changed("min-votes") {
		finalMinVotes = popularMinVotes
	}

	finalTimeout := cfg.Timeout
	if cmd.Flags().Changed("timeout") {
		finalTimeout = popularTimeout
	}

	client, err := api.NewClient(cfg.APIKey, finalTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	desiredProviders := filters.ParseProviders(finalProviders)

	var genreList []models.Genre
	var genreMap map[string]int

	genreResp, err := client.GetGenres("de-DE")
	if err == nil {
		genreList = genreResp.Genres
		genreMap = filters.BuildGenreMap(genreList)
	}

	fmt.Printf("Searching TMDb's Popular Movies...\n")
	fmt.Printf("Criteria: Min Rating: %.1f | Min Votes: %d\n", finalMinRating, finalMinVotes)
	fmt.Printf("Filtering for [%s] in region [%s]\n\n", finalProviders, strings.ToUpper(finalRegion))

	resultsFound := 0
	for page := 1; page <= config.MaxPagesToSearch && resultsFound < config.MaxResultsToDisplay; page++ {
		fmt.Printf("Fetching page %d...\n", page)

		resp, err := client.GetPopularMovies(page, finalRegion)
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

			if popularGenre != "" && !filters.FilterByGenre(&movie, popularGenre, genreMap) {
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

			// Get region-specific title (e.g., DE -> de-DE)
			languageCode := strings.ToLower(finalRegion) + "-" + strings.ToUpper(finalRegion)
			regionalTitle, _ := client.GetRegionalTitle(movie.ID, languageCode)
			if regionalTitle == "" {
				regionalTitle = movie.Title
			}

			genreNames := filters.GetGenreNames(movie.GenreIDs, genreList)

			display.DisplayMovie(display.MovieDisplay{
				Number:       resultsFound,
				Title:        regionalTitle,
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

		if page >= resp.TotalPages {
			break
		}
	}

	display.DisplaySeparator()
	if resultsFound == 0 {
		fmt.Println("No movies found matching criteria.")
	} else {
		fmt.Printf("Displayed %d popular movies.\n", resultsFound)
	}
}
