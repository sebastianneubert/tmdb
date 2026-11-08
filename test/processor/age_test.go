package processor_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/sebastianneubert/tmdb/internal/models"
	"github.com/sebastianneubert/tmdb/internal/processor"
)

// ageLookupFromMap builds a lookup that returns the given ages. Movies that are
// not in the map have no age rating; movies mapped to -1 make the lookup fail.
func ageLookupFromMap(ages map[int]int) processor.AgeLookupFunc {
	return func(movieID int) (int, bool, error) {
		age, ok := ages[movieID]
		if !ok {
			return 0, false, nil
		}
		if age < 0 {
			return 0, false, errors.New("lookup failed")
		}
		return age, true, nil
	}
}

func runWithAgeFilter(t *testing.T, filter *processor.AgeFilter) []int {
	t.Helper()

	mp := processor.NewMovieProcessor(nil, processor.FilterConfig{
		MinRating: 7.0,
		MinVotes:  100,
		Region:    "DE",
		AgeFilter: filter,
	})

	var shown []int
	err := mp.Process(
		func(page int) (*models.DiscoverResponse, error) {
			return &models.DiscoverResponse{
				Page:       page,
				TotalPages: 1,
				Results: []models.Movie{
					{ID: 1, Title: "FSK 6", VoteAverage: 8, VoteCount: 1000},
					{ID: 2, Title: "FSK 12", VoteAverage: 8, VoteCount: 1000},
					{ID: 3, Title: "FSK 16", VoteAverage: 8, VoteCount: 1000},
					{ID: 4, Title: "Unrated", VoteAverage: 8, VoteCount: 1000},
					{ID: 5, Title: "Lookup error", VoteAverage: 8, VoteCount: 1000},
					{ID: 6, Title: "Low rating, never looked up", VoteAverage: 5, VoteCount: 1000},
				},
			}, nil
		},
		func(movie *models.Movie, providers []string, genres []string) error {
			shown = append(shown, movie.ID)
			return nil
		},
	)
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}
	return shown
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestProcessorWithoutAgeFilterShowsAllMovies(t *testing.T) {
	shown := runWithAgeFilter(t, nil)
	want := []int{1, 2, 3, 4, 5}
	if !equalInts(shown, want) {
		t.Errorf("shown = %v, want %v", shown, want)
	}
}

func TestProcessorAgeFilterHidesOlderAndUnrated(t *testing.T) {
	filter := &processor.AgeFilter{
		MaxAge: 12,
		Region: "DE",
		Lookup: ageLookupFromMap(map[int]int{1: 6, 2: 12, 3: 16, 5: -1}),
	}

	shown := runWithAgeFilter(t, filter)
	want := []int{1, 2}
	if !equalInts(shown, want) {
		t.Errorf("shown = %v, want %v", shown, want)
	}

	summary := filter.Summary()
	if !strings.Contains(summary, "1 hidden as too old") {
		t.Errorf("summary %q should report 1 movie as too old", summary)
	}
	if !strings.Contains(summary, "2 hidden without age rating") {
		t.Errorf("summary %q should report 2 movies without age rating (unrated + failed lookup)", summary)
	}
}

func TestProcessorAgeFilterIncludeUnrated(t *testing.T) {
	filter := &processor.AgeFilter{
		MaxAge:         12,
		IncludeUnrated: true,
		Region:         "DE",
		Lookup:         ageLookupFromMap(map[int]int{1: 6, 2: 12, 3: 16, 5: -1}),
	}

	shown := runWithAgeFilter(t, filter)
	want := []int{1, 2, 4, 5}
	if !equalInts(shown, want) {
		t.Errorf("shown = %v, want %v", shown, want)
	}

	if summary := filter.Summary(); strings.Contains(summary, "without age rating") {
		t.Errorf("summary %q should not mention unrated movies when they are included", summary)
	}
}

func TestAgeFilterAgeLabel(t *testing.T) {
	filter := &processor.AgeFilter{
		MaxAge:         12,
		IncludeUnrated: true,
		Region:         "DE",
		Lookup:         ageLookupFromMap(map[int]int{1: 0, 2: 12, 3: 16, 5: -1}),
	}

	if label := filter.AgeLabel(2); label != "" {
		t.Errorf("label before lookup = %q, want empty", label)
	}

	for _, id := range []int{1, 2, 3, 4, 5} {
		filter.Allows(id)
	}

	tests := []struct {
		movieID int
		want    string
	}{
		{1, "0+"},
		{2, "12+"},
		{3, "16+"}, // too old to be shown, but the age is known
		{4, ""},    // no age rating
		{5, ""},    // failed lookup
		{99, ""},   // never looked up
	}
	for _, tt := range tests {
		if got := filter.AgeLabel(tt.movieID); got != tt.want {
			t.Errorf("AgeLabel(%d) = %q, want %q", tt.movieID, got, tt.want)
		}
	}
}

func TestProcessorPassesAgeToCallback(t *testing.T) {
	filter := &processor.AgeFilter{
		MaxAge:         12,
		IncludeUnrated: true,
		Region:         "DE",
		Lookup:         ageLookupFromMap(map[int]int{1: 6, 2: 12, 3: 16}),
	}

	mp := processor.NewMovieProcessor(nil, processor.FilterConfig{
		MinRating: 7.0,
		MinVotes:  100,
		Region:    "DE",
		AgeFilter: filter,
	})

	labels := map[int]string{}
	err := mp.Process(
		func(page int) (*models.DiscoverResponse, error) {
			return &models.DiscoverResponse{
				Page:       page,
				TotalPages: 1,
				Results: []models.Movie{
					{ID: 1, VoteAverage: 8, VoteCount: 1000},
					{ID: 2, VoteAverage: 8, VoteCount: 1000},
					{ID: 3, VoteAverage: 8, VoteCount: 1000},
					{ID: 4, VoteAverage: 8, VoteCount: 1000},
				},
			}, nil
		},
		func(movie *models.Movie, providers []string, genres []string) error {
			labels[movie.ID] = filter.AgeLabel(movie.ID)
			return nil
		},
	)
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}

	want := map[int]string{1: "6+", 2: "12+", 4: ""}
	if len(labels) != len(want) {
		t.Fatalf("labels = %v, want %v", labels, want)
	}
	for id, label := range want {
		if labels[id] != label {
			t.Errorf("label of movie %d = %q, want %q", id, labels[id], label)
		}
	}
}

func TestAgeFilterNilIsSafe(t *testing.T) {
	var filter *processor.AgeFilter

	if filter.AgeLabel(1) != "" {
		t.Error("nil filter must have an empty age label")
	}

	if filter.Enabled() {
		t.Error("nil filter must not be enabled")
	}
	if !filter.Allows(1) {
		t.Error("nil filter must allow every movie")
	}
	if filter.Summary() != "" {
		t.Error("nil filter must have an empty summary")
	}
}

func TestNewAgeFilterIsNilWhenOff(t *testing.T) {
	if f := processor.NewAgeFilter(nil, "DE", 0, false); f != nil {
		t.Errorf("expected nil filter for maxAge 0, got %+v", f)
	}
}
