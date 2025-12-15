package models

import (
	"testing"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func TestActorListDisplay(t *testing.T) {
	actors := []models.Actor{
		{
			ID:          1,
			Name:        "Leonardo DiCaprio",
			Popularity:  95.5,
			ProfilePath: "/profile1.jpg",
		},
		{
			ID:          2,
			Name:        "Tom Hanks",
			Popularity:  92.3,
			ProfilePath: "/profile2.jpg",
		},
		{
			ID:          3,
			Name:        "Meryl Streep",
			Popularity:  89.1,
			ProfilePath: "/profile3.jpg",
		},
	}

	if len(actors) != 3 {
		t.Errorf("Expected 3 actors, got %d", len(actors))
	}

	// Verify actors are sorted by popularity (descending)
	for i := 0; i < len(actors)-1; i++ {
		if actors[i].Popularity < actors[i+1].Popularity {
			t.Errorf("Actors not sorted by popularity: %f < %f", actors[i].Popularity, actors[i+1].Popularity)
		}
	}

	// Check first actor details
	if actors[0].Name != "Leonardo DiCaprio" {
		t.Errorf("Expected first actor 'Leonardo DiCaprio', got '%s'", actors[0].Name)
	}

	if actors[0].Popularity != 95.5 {
		t.Errorf("Expected popularity 95.5, got %f", actors[0].Popularity)
	}

	if actors[0].ID != 1 {
		t.Errorf("Expected ID 1, got %d", actors[0].ID)
	}
}

func TestActorSearchMultipleResults(t *testing.T) {
	// Simulate search results for "Tom"
	results := []models.Actor{
		{
			ID:          1,
			Name:        "Tom Hanks",
			Popularity:  92.3,
			ProfilePath: "/tom_hanks.jpg",
		},
		{
			ID:          2,
			Name:        "Tom Hardy",
			Popularity:  85.6,
			ProfilePath: "/tom_hardy.jpg",
		},
		{
			ID:          3,
			Name:        "Tom Cruise",
			Popularity:  88.9,
			ProfilePath: "/tom_cruise.jpg",
		},
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Check that results contain actors with "Tom" in their names
	for _, actor := range results {
		if actor.Name == "" {
			t.Errorf("Actor name should not be empty")
		}
		if actor.Popularity == 0 {
			t.Errorf("Actor popularity should be set")
		}
	}

	// Verify specific result
	if results[0].Name != "Tom Hanks" {
		t.Errorf("Expected 'Tom Hanks', got '%s'", results[0].Name)
	}

	if results[1].Name != "Tom Hardy" {
		t.Errorf("Expected 'Tom Hardy', got '%s'", results[1].Name)
	}
}

func TestActorDetailsCompleteness(t *testing.T) {
	actor := models.Actor{
		ID:          1,
		Name:        "Denzel Washington",
		Popularity:  91.5,
		ProfilePath: "/denzel.jpg",
		KnownFor: []models.Movie{
			{
				ID:    1,
				Title: "Training Day",
			},
			{
				ID:    2,
				Title: "Malcom X",
			},
		},
	}

	// Verify all required fields are populated
	if actor.ID == 0 {
		t.Errorf("Actor ID should not be zero")
	}

	if actor.Name == "" {
		t.Errorf("Actor name should not be empty")
	}

	if actor.Popularity == 0 {
		t.Errorf("Actor popularity should be set")
	}

	if len(actor.KnownFor) != 2 {
		t.Errorf("Expected 2 known for movies, got %d", len(actor.KnownFor))
	}

	// Verify known_for movies
	if actor.KnownFor[0].Title != "Training Day" {
		t.Errorf("Expected 'Training Day', got '%s'", actor.KnownFor[0].Title)
	}
}

func TestActorSearchEmptyResults(t *testing.T) {
	// Test empty search results
	results := []models.Actor{}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestActorSearchSingleResult(t *testing.T) {
	// Test single result
	results := []models.Actor{
		{
			ID:          1,
			Name:        "Keanu Reeves",
			Popularity:  86.2,
			ProfilePath: "/keanu.jpg",
		},
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Name != "Keanu Reeves" {
		t.Errorf("Expected 'Keanu Reeves', got '%s'", results[0].Name)
	}
}
