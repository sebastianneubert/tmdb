package commands

import (
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/spf13/cobra"
)

// MovieCommandFlags holds all the common flag values for movie-related commands
type MovieCommandFlags struct {
	Providers string
	Region    string
	MinRating float64
	MinVotes  int
	Timeout   int
	Genre     string
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
