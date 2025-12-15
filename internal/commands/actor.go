package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/display"
	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/sebastianneubert/tmdb/internal/models"
	"github.com/spf13/cobra"
)

var (
	actorProviders string
	actorRegion    string
	actorMinRating float64
	actorMinVotes  int
	actorTimeout   int
	actorGenre     string
	actorList      bool
)

var actorCmd = &cobra.Command{
	Use:   "actor [name]",
	Short: "Find an actor's filmography with streaming availability.",
	Args:  cobra.MaximumNArgs(1),
	Run:   runActor,
	Long: `Search for an actor and display their movies with ratings and availability.
If no name is provided with --list flag, shows popular actors.
If search results have multiple matches, displays a list to choose from.

Examples:
  tmdb actor
  tmdb actor "Leonardo DiCaprio"
  tmdb actor "Tom"`,
}

func init() {
	actorCmd.Flags().StringVarP(&actorProviders, "providers", "p", config.DefaultProviders, "Comma-separated providers")
	actorCmd.Flags().StringVarP(&actorRegion, "region", "r", config.DefaultRegion, "Watch region")
	actorCmd.Flags().Float64Var(&actorMinRating, "min-rating", config.DefaultMinRating, "Minimum rating")
	actorCmd.Flags().IntVar(&actorMinVotes, "min-votes", config.DefaultMinVotes, "Minimum votes")
	actorCmd.Flags().IntVarP(&actorTimeout, "timeout", "T", config.DefaultTimeout, "Timeout in seconds")
	actorCmd.Flags().StringVar(&actorGenre, "genre", "", "Filter by genre (name or ID)")
	actorCmd.Flags().BoolVar(&actorList, "list", false, "List actors instead of fetching filmography")
}

func runActor(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	var actorName string
	if len(args) > 0 {
		actorName = strings.Join(args, " ")
	}

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

	// If --list flag is set and no actor name provided, show popular actors
	if actorList || actorName == "" {
		displayPopularActors(client, finalRegion)
		return
	}

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

	// If multiple results and --list flag is set, show the list
	if len(actorResults.Results) > 1 && actorList {
		displayActorMatches(actorResults.Results)
		return
	}

	// If multiple results and no --list flag, show matches prompt
	if len(actorResults.Results) > 1 {
		fmt.Printf("Found %d actors matching '%s'. Did you mean one of these?\n\n", len(actorResults.Results), actorName)
		displayActorMatches(actorResults.Results)
		fmt.Printf("\nTo view filmography, use:\n  tmdb actor \"%s\"\n\n", actorResults.Results[0].Name)
		return
	}

	// Proceed with single match
	actor := actorResults.Results[0]
	displayActorFilmography(client, actor, finalRegion, finalProviders, finalMinRating, finalMinVotes, desiredProviders)
}

func displayActorMatches(actors []models.Actor) {
	display.DisplaySeparator()

	// Sort actors by popularity (descending)
	sortedActors := actors
	sort.Slice(sortedActors, func(i, j int) bool {
		return sortedActors[i].Popularity > sortedActors[j].Popularity
	})

	// Display top popularactors
	displayCount := 0
	maxDisplay := 15

	for i := 0; i < len(sortedActors) && displayCount < maxDisplay; i++ {
		actor := sortedActors[i]
		displayCount++
		display.DisplayActor(display.ActorDisplay{
			Number:     i + 1,
			Name:       actor.Name,
			Popularity: actor.Popularity,
			TmdbID:     actor.ID,
		})
	}
	display.DisplaySeparator()
}

func displayPopularActors(client *api.Client, language string) {
	fmt.Println("Fetching popular actors...")
	fmt.Println("Popular actors:")

	results, err := client.GetPopularActors(language, 1)
	if err != nil {
		fmt.Printf("Error fetching popular actors: %v\n", err)
		return
	}

	if len(results.Results) == 0 {
		fmt.Println("No popular actors found.")
		return
	}

	// Sort actors by popularity (descending)
	sortedActors := results.Results
	sort.Slice(sortedActors, func(i, j int) bool {
		return sortedActors[i].Popularity > sortedActors[j].Popularity
	})

	display.DisplaySeparator()

	// Display top 10 actors
	displayCount := 0
	maxDisplay := 15

	for i := 0; i < len(sortedActors) && displayCount < maxDisplay; i++ {
		actor := sortedActors[i]
		displayCount++
		display.DisplayActor(display.ActorDisplay{
			Number:     displayCount,
			Name:       actor.Name,
			Popularity: actor.Popularity,
			TmdbID:     actor.ID,
		})
	}
	display.DisplaySeparator()
	fmt.Printf("Showing top %d popular actors\n", displayCount)
}

func displayActorFilmography(client *api.Client, actor models.Actor, finalRegion, finalProviders string, finalMinRating float64, finalMinVotes int, desiredProviders map[string]bool) {
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
		if !filters.MeetsRatingCriteria(movie.VoteAverage, movie.VoteCount, finalMinRating, finalMinVotes) {
			continue
		}

		moviesChecked++

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
