package commands

import (
	"fmt"
	"sort"
	"strconv"
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
	Use:   "actor [name] [index]",
	Short: "Find an actor's filmography with streaming availability.",
	Args:  cobra.MaximumNArgs(2),
	Run:   runActor,
	Long: `Search for an actor and display their movies with ratings and availability.
If no name is provided with --list flag, shows popular actors.
If search results have multiple matches, displays a list to choose from.
You can specify an index to select a specific actor from multiple results.
Movies are displayed with genres and filtered by region-specific titles.

Examples:
  tmdb actor
  tmdb actor "Leonardo DiCaprio"
  tmdb actor "Tom" 1
  tmdb actor "Megan Fox" 2
  tmdb actor "Tom Hanks" --genre Action`,
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
	var actorIndex int = -1

	if len(args) > 0 {
		actorName = args[0]
	}

	// Check if second argument is provided (actor index)
	if len(args) > 1 {
		var err error
		actorIndex, err = strconv.Atoi(args[1])
		if err != nil || actorIndex < 1 {
			fmt.Printf("Invalid actor index: %s. Please provide a positive number (1, 2, 3, ...)\n", args[1])
			return
		}
		// Convert to 0-based index
		actorIndex--
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

	// Fetch genres for filtering
	var genreList []models.Genre
	var genreMap map[string]int

	genreResp, err := client.GetGenres(finalRegion)
	if err == nil {
		genreList = genreResp.Genres
		genreMap = filters.BuildGenreMap(genreList)
	}

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

	// Sort all results by popularity (descending) so index matches displayed order
	sort.Slice(actorResults.Results, func(i, j int) bool {
		return actorResults.Results[i].Popularity > actorResults.Results[j].Popularity
	})

	// If actor index is provided, use it to select from sorted results
	if actorIndex >= 0 {
		if actorIndex >= len(actorResults.Results) {
			fmt.Printf("Invalid actor index: %d. Found only %d actors matching '%s' (use 1-%d)\n", 
				actorIndex+1, len(actorResults.Results), actorName, len(actorResults.Results))
			displayActorMatches(actorResults.Results)
			return
		}
		actor := actorResults.Results[actorIndex]
		displayActorFilmography(client, actor, finalRegion, finalProviders, finalMinRating, finalMinVotes, desiredProviders, actorGenre, genreList, genreMap)
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
		fmt.Printf("\nTo view filmography, use:\n  tmdb actor \"%s\" 1\n\n", actorName)
		return
	}

	// Proceed with single match
	actor := actorResults.Results[0]
	displayActorFilmography(client, actor, finalRegion, finalProviders, finalMinRating, finalMinVotes, desiredProviders, actorGenre, genreList, genreMap)
}

func displayActorMatches(actors []models.Actor) {
	display.DisplaySeparator()

	// Display top actors (already sorted by popularity)
	displayCount := 0
	maxDisplay := 15

	for i := 0; i < len(actors) && displayCount < maxDisplay; i++ {
		actor := actors[i]
		displayCount++
		display.DisplayActor(display.ActorDisplay{
			Number:     displayCount,
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

func displayActorFilmography(client *api.Client, actor models.Actor, finalRegion, finalProviders string, finalMinRating float64, finalMinVotes int, desiredProviders map[string]bool, genreFilter string, genreList []models.Genre, genreMap map[string]int) {
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

		// Apply genre filter
		if genreFilter != "" && !filters.FilterByGenre(&movie, genreFilter, genreMap) {
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
			Character:    movie.Character,
			Genres:       genreNames,
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
