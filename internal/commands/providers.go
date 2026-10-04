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
	providersRegion string
	providersSearch string
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List streaming providers known to TMDb for a region.",
	Long: `Display the streaming providers TMDb knows for a region, exactly as TMDb names them.
Providers matched by your current --providers setting are marked with a check mark.
Entries of your setting that match no provider at all are reported as well.

Examples:
  tmdb providers
  tmdb providers --region AT
  tmdb providers --search disney`,
	Run: runProviders,
}

func init() {
	providersCmd.Flags().StringVarP(&providersRegion, "region", "r", config.DefaultRegion, "Watch region")
	providersCmd.Flags().StringVarP(&providersSearch, "search", "s", "", "Only show providers whose name contains this text")
}

func runProviders(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	region := cfg.Region
	if cmd.Flags().Changed("region") {
		region = providersRegion
	}

	client, err := api.NewClient(cfg.APIKey, cfg.Timeout)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("📺 Fetching streaming providers for %s...\n\n", region)

	resp, err := client.GetWatchProviderList(region)
	if err != nil {
		fmt.Printf("Error fetching providers: %v\n", err)
		return
	}

	if len(resp.Results) == 0 {
		fmt.Println("No providers found.")
		return
	}

	all := resp.Results
	sort.Slice(all, func(i, j int) bool {
		return strings.ToLower(all[i].Name) < strings.ToLower(all[j].Name)
	})

	desired := filters.ParseProviders(cfg.Providers)

	shown := all
	if providersSearch != "" {
		needle := strings.ToLower(providersSearch)
		shown = nil
		for _, p := range all {
			if strings.Contains(strings.ToLower(p.Name), needle) {
				shown = append(shown, p)
			}
		}
		if len(shown) == 0 {
			fmt.Printf("No providers matching %q in %s.\n", providersSearch, region)
			return
		}
	}

	fmt.Println(display.SeparatorStyle.Render(strings.Repeat("=", 60)))
	fmt.Printf("Streaming Providers in %s (%d shown, %d total)\n", region, len(shown), len(all))
	fmt.Println(display.SeparatorStyle.Render(strings.Repeat("=", 60)))

	for _, p := range shown {
		marker := " "
		if matchesProvider(p.Name, desired) {
			marker = "✓"
		}
		line := fmt.Sprintf("%s %-35s %-6d", marker, p.Name, p.ID)
		fmt.Println(display.ProviderStyle.Render(line))
	}

	fmt.Println(display.SeparatorStyle.Render(strings.Repeat("=", 60)))
	fmt.Printf("\nCurrent providers setting: %s\n", cfg.Providers)

	// Report entries of the current setting that match no provider in this region.
	var unmatched []string
	for entry := range desired {
		single := map[string]bool{entry: true}
		found := false
		for _, p := range all {
			if matchesProvider(p.Name, single) {
				found = true
				break
			}
		}
		if !found {
			unmatched = append(unmatched, entry)
		}
	}
	sort.Strings(unmatched)

	if len(unmatched) > 0 {
		fmt.Println("\n⚠️  These entries match no TMDb provider name and will never find anything:")
		for _, entry := range unmatched {
			fmt.Printf("   %s\n", entry)
		}
		fmt.Println("   Use the names from the list above, e.g. tmdb top --providers \"Netflix,Disney Plus\"")
	}
}

// matchesProvider reuses the matching logic of the search commands for a single provider name.
func matchesProvider(name string, desired map[string]bool) bool {
	_, ok := filters.CheckAvailability(models.RegionProviders{
		Flatrate: []models.Provider{{ProviderName: name}},
	}, desired)
	return ok
}
