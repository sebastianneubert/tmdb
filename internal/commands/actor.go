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

var (
	actorProviders string
	actorRegion    string
	actorMinRating float64
	actorMinVotes  int
	actorTimeout   int
)

var actorCmd = &cobra.Command{
	Use:   "actor [name]",
	Short: "Find an actor's filmography with streaming availability.",
	Long:  "Search for an actor and display their movies with ratings and availability.",
	Args:  cobra.MinimumNArgs(1),
	Run:   runActor,
}

func init() {
	actorCmd.Flags().StringVarP(&actorProviders, "providers", "p", config.DefaultProviders, "Comma-separated providers")
	actorCmd.Flags().StringVarP(&actorRegion, "region", "r", config.DefaultRegion, "Watch region")
	actorCmd.Flags().Float64Var(&actorMinRating, "min-rating", config.DefaultMinRating, "Minimum rating")
	actorCmd.Flags().IntVar(&actorMinVotes, "min-votes", config.DefaultMinVotes, "Minimum votes")
	actorCmd.Flags().IntVarP(&actorTimeout, "timeout", "T", config.DefaultTimeout, "Timeout in seconds")
}

func runActor(cmd *cobra.Command, args []string) {
	actorName := strings.Join(args, " ")
	cfg := config.Get()

	finalRegion := cfg.Region
	if cmd.Flags().Changed("region") {
		finalRegion = actorRegion
	}

	finalProviders := cfg.Providers
	if cmd.Flags().Changed("providers") {
		finalProviders = actorProviders
	}

	finalMinRating := cfg.MinRating
	if cmd.Flags().Changed("min-rating") {
		finalMinRating = actorMinRating
	}

	finalMinVotes := cfg.MinVotes
	if cmd.Flags().Changed("min-votes") {
		finalMinVotes = actorMinVotes
	}

	finalTimeout := cfg.Timeout
	if cmd.Flags().Changed("timeout") {
		finalTimeout = actorTimeout
	}

	client, err := api.NewClient(cfg.APIKey, finalTimeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	desiredProviders := filters.ParseProviders(finalProviders)

	fmt.Printf("Searching for actor: %s\n\n", actorName)

	actorResults, err := client.SearchActor(actorName, finalRegion)
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
		return
	}

	if len(actorResults.Results) == 0 {
		fmt.Printf("No actors found matching '%s'\n", actorName)
		return
	}

	actor := actorResults.Results[0]
	fmt.Printf("Found: %s (TMDb ID: %d)\n", display.TitleStyle.Render(actor.Name), actor.ID)
	fmt.Printf("Fetching filmography...\n\n")

	credits, err := client.GetActorCredits(actor.ID, finalRegion)
	if err != nil {
		fmt.Printf("Error fetching filmography: %v\n", err)
		return
	}

	if len(credits.Cast) == 0 {
		fmt.Printf("No movie credits found.\n")
		return
	}

	fmt.Printf("Filtering with Min Rating: %.1f | Min Votes: %d\n", finalMinRating, finalMinVotes)
	fmt.Printf("Checking [%s] in region [%s]\n\n", finalProviders, strings.ToUpper(finalRegion))

	resultsFound := 0
	moviesChecked := 0

	for _, movie := range credits.Cast {
		if !filters.MeetsRatingCriteria(movie.VoteAverage, movie.VoteCount, actorMinRating, actorMinVotes) {
			continue
		}

		moviesChecked++

		providerData, err := client.GetWatchProviders(movie.ID, actorRegion)
		if err != nil {
		  if cfg.DEBUG {
		    fmt.Printf("Error fetching providers for %s: %v\n", movie.Title, err)
		  }
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
			Number:       resultsFound,
			Title:        movie.Title,
			EnglishTitle: englishTitle,
			Year:         movie.GetYear(),
			Rating:       movie.VoteAverage,
			Votes:        movie.VoteCount,
			Providers:    availableProviders,
			TmdbID:       movie.ID,
			ImdbID:       externalIDs.ImdbID,
			Overview:     movie.Overview,
			Character:    movie.Character,
		})

		if resultsFound >= config.MaxResultsToDisplay {
			break
		}
	}

	display.DisplaySeparator()
	if resultsFound == 0 {
		fmt.Printf("No movies found for %s.\n", actor.Name)
		fmt.Printf("(Checked %d movies meeting criteria)\n", moviesChecked)
	} else {
		fmt.Printf("Found %d movies starring %s.\n", resultsFound, actor.Name)
	}
}
