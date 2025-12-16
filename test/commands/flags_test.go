package commands_test

import (
	"testing"

	"github.com/sebastianneubert/tmdb/internal/commands"
	"github.com/sebastianneubert/tmdb/internal/config"
	"github.com/spf13/cobra"
)

func TestMovieCommandFlagsRegisterAndResolve(t *testing.T) {
	cmd := &cobra.Command{}
	flags := commands.MovieCommandFlags{}
	flags.Register(cmd, true)

	// Simulate setting flags
	cmd.Flags().Set("region", "US")
	cmd.Flags().Set("providers", "Netflix")
	cmd.Flags().Set("min-rating", "7.5")
	cmd.Flags().Set("min-votes", "100")
	cmd.Flags().Set("timeout", "5")
	cmd.Flags().Set("genre", "Action")

	cfg := config.Config{APIKey: "", Region: "DE"}
	region, providers, minRating, minVotes, timeout, genre := flags.Resolve(cmd, cfg)

	if region != "US" {
		t.Errorf("Expected region US, got %s", region)
	}
	if providers != "Netflix" {
		t.Errorf("Expected providers Netflix, got %s", providers)
	}
	if genre != "Action" {
		t.Errorf("Expected genre Action, got %s", genre)
	}
	if minRating != 7.5 {
		t.Errorf("Expected minRating 7.5, got %f", minRating)
	}
	if minVotes != 100 {
		t.Errorf("Expected minVotes 100, got %d", minVotes)
	}
	if timeout == 0 {
		t.Errorf("Expected non-zero timeout")
	}
}
