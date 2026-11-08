package commands

import (
	"fmt"
	"strings"

	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/spf13/cobra"
)

// MovieCommandFlags holds all the common flag values for movie-related commands
type MovieCommandFlags struct {
	Providers      string
	Region         string
	MinRating      float64
	MinVotes       int
	Timeout        int
	Genre          string
	MaxAge         int
	IncludeUnrated bool
}

// registerAgeFlags registers --max-age and --include-unrated on a command.
func registerAgeFlags(cmd *cobra.Command, maxAge *int, includeUnrated *bool) {
	cmd.Flags().IntVar(maxAge, "max-age", config.DefaultMaxAge,
		"Hide movies with an age rating above this age in the watch region, e.g. 12 (0 = off)")
	cmd.Flags().BoolVar(includeUnrated, "include-unrated", false,
		"With --max-age: also show movies that have no age rating in the watch region")
}

// resolveAgeFlags combines config values with command-line overrides.
func resolveAgeFlags(cmd *cobra.Command, cfg config.Config, maxAgeFlag int, includeUnratedFlag bool) (maxAge int, includeUnrated bool) {
	maxAge = cfg.MaxAge
	if cmd.Flags().Changed("max-age") {
		maxAge = maxAgeFlag
	}

	includeUnrated = cfg.IncludeUnrated
	if cmd.Flags().Changed("include-unrated") {
		includeUnrated = includeUnratedFlag
	}

	return
}

// printAgeFilterNotice tells the user that the age filter is active.
func printAgeFilterNotice(maxAge int, includeUnrated bool, region string) {
	if maxAge <= 0 {
		return
	}

	notice := fmt.Sprintf("Age filter: max. %d (age ratings of region %s)", maxAge, strings.ToUpper(region))
	if includeUnrated {
		notice += ", unrated movies included"
	}
	fmt.Printf("%s\n\n", notice)
}

// RegisterAge registers the age filter flags (--max-age, --include-unrated).
// Call it for movie commands that apply the age filter.
func (f *MovieCommandFlags) RegisterAge(cmd *cobra.Command) {
	registerAgeFlags(cmd, &f.MaxAge, &f.IncludeUnrated)
}

// ResolveAge returns the final age filter settings (config defaults plus CLI overrides).
func (f *MovieCommandFlags) ResolveAge(cmd *cobra.Command, cfg config.Config) (maxAge int, includeUnrated bool) {
	return resolveAgeFlags(cmd, cfg, f.MaxAge, f.IncludeUnrated)
}

// Register registers all flags with the given command
// If includeGenre is false, the genre flag will not be registered
func (f *MovieCommandFlags) Register(cmd *cobra.Command, includeGenre bool) {
	cmd.Flags().StringVarP(&f.Providers, "providers", "p", config.DefaultProviders, "Comma-separated providers")
	cmd.Flags().StringVarP(&f.Region, "region", "r", config.DefaultRegion, "Watch region")
	cmd.Flags().Float64Var(&f.MinRating, "min-rating", config.DefaultMinRating, "Minimum rating")
	cmd.Flags().IntVar(&f.MinVotes, "min-votes", config.DefaultMinVotes, "Minimum votes")
	cmd.Flags().IntVarP(&f.Timeout, "timeout", "T", config.DefaultTimeout, "Timeout in seconds")
	if includeGenre {
		cmd.Flags().StringVar(&f.Genre, "genre", "", "Filter by genre (name or ID)")
	}
}

// Resolve returns the final values by combining config defaults with any command-line overrides
func (f *MovieCommandFlags) Resolve(cmd *cobra.Command, cfg config.Config) (region, providers string, minRating float64, minVotes, timeout int, genre string) {
	region = cfg.Region
	if cmd.Flags().Changed("region") {
		region = f.Region
	}

	providers = cfg.Providers
	if cmd.Flags().Changed("providers") {
		providers = f.Providers
	}

	minRating = cfg.MinRating
	if cmd.Flags().Changed("min-rating") {
		minRating = f.MinRating
	}

	minVotes = cfg.MinVotes
	if cmd.Flags().Changed("min-votes") {
		minVotes = f.MinVotes
	}

	timeout = cfg.Timeout
	if cmd.Flags().Changed("timeout") {
		timeout = f.Timeout
	}

	genre = f.Genre

	return
}
