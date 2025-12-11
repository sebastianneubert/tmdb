package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/sebastianneubert/tmdb/internal/display"
)

var (
	genresLanguage string
)

var genresCmd = &cobra.Command{
	Use:   "genres",
	Short: "List all available movie genres.",
	Long: `Display a list of all movie genres available on TMDb.
Use genre IDs or names with --genre flag in other commands.

Examples:
  tmdb genres
  tmdb genres --language en-US`,
	Run: runGenres,
}

func init() {
	genresCmd.Flags().StringVarP(&genresLanguage, "language", "l", "de-DE", "Language for genre names")
}

func runGenres(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	client, err := api.NewClient(cfg.APIKey, cfg.Timeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("🎭 Fetching movie genres...\n\n")

	genreResp, err := client.GetGenres(genresLanguage)
	if err != nil {
		fmt.Printf("Error fetching genres: %v\n", err)
		return
	}

	if len(genreResp.Genres) == 0 {
		fmt.Println("No genres found.")
		return
	}

	// Sort genres by name
	sort.Slice(genreResp.Genres, func(i, j int) bool {
		return genreResp.Genres[i].Name < genreResp.Genres[j].Name
	})

	fmt.Println(display.SeparatorStyle.Render(strings.Repeat("=", 60)))
	fmt.Printf("Available Movie Genres (%d total)\n", len(genreResp.Genres))
	fmt.Println(display.SeparatorStyle.Render(strings.Repeat("=", 60)))

	// Display in two columns
	half := (len(genreResp.Genres) + 1) / 2
	for i := 0; i < half; i++ {
		left := genreResp.Genres[i]
		leftText := fmt.Sprintf("%-20s %-5d", left.Name, left.ID)

		if i+half < len(genreResp.Genres) {
			right := genreResp.Genres[i+half]
			rightText := fmt.Sprintf("%-20s %-5d", right.Name, right.ID)
			fmt.Printf("%s | %s\n",
				display.ProviderStyle.Render(leftText),
				display.ProviderStyle.Render(rightText))
		} else {
			fmt.Printf("%s\n", display.ProviderStyle.Render(leftText))
		}
	}

	fmt.Println(display.SeparatorStyle.Render(strings.Repeat("=", 60)))
	fmt.Println("\n💡 Usage examples:")
	fmt.Println("   tmdb top --genre Action")
	fmt.Println("   tmdb search \"star\" --genre \"Science Fiction\"")
	fmt.Println("   tmdb actor \"Tom Hanks\" --genre Drama")
}