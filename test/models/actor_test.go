package models

import (
	"testing"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func TestActorSearchResponse(t *testing.T) {
	actor := models.Actor{
		ID:          1,
		Name:        "Tom Hanks",
		Popularity:  95.5,
		ProfilePath: "/path/to/profile.jpg",
	}

	if actor.ID != 1 {
		t.Errorf("Expected ID 1, got %d", actor.ID)
	}

	if actor.Name != "Tom Hanks" {
		t.Errorf("Expected Name 'Tom Hanks', got '%s'", actor.Name)
	}

	if actor.Popularity != 95.5 {
		t.Errorf("Expected Popularity 95.5, got %f", actor.Popularity)
	}
}

func TestActorCreditsResponse(t *testing.T) {
	movie1 := models.Movie{
		ID:          1,
		Title:       "Forrest Gump",
		ReleaseDate: "1994-07-06",
		VoteAverage: 8.8,
		VoteCount:   1000,
		Character:   "Forrest Gump",
	}

	movie2 := models.Movie{
		ID:          2,
		Title:       "Toy Story",
		ReleaseDate: "1995-11-22",
		VoteAverage: 8.3,
		VoteCount:   15000,
		Character:   "Woody",
	}

	response := models.ActorCreditsResponse{
		ID:   1,
		Cast: []models.Movie{movie1, movie2},
	}

	if len(response.Cast) != 2 {
		t.Errorf("Expected 2 movies, got %d", len(response.Cast))
	}

	if response.Cast[0].Title != "Forrest Gump" {
		t.Errorf("Expected first movie title 'Forrest Gump', got '%s'", response.Cast[0].Title)
	}

	if response.Cast[0].Character != "Forrest Gump" {
		t.Errorf("Expected character 'Forrest Gump', got '%s'", response.Cast[0].Character)
	}

	if response.Cast[1].Character != "Woody" {
		t.Errorf("Expected character 'Woody', got '%s'", response.Cast[1].Character)
	}
}

func TestActorWithKnownFor(t *testing.T) {
	movie := models.Movie{
		ID:    1,
		Title: "Famous Movie",
	}

	actor := models.Actor{
		ID:       2,
		Name:     "Great Actor",
		KnownFor: []models.Movie{movie},
	}

	if len(actor.KnownFor) != 1 {
		t.Errorf("Expected 1 known for movie, got %d", len(actor.KnownFor))
	}

	if actor.KnownFor[0].Title != "Famous Movie" {
		t.Errorf("Expected known for title 'Famous Movie', got '%s'", actor.KnownFor[0].Title)
	}
}
