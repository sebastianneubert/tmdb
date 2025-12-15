package models

import (
	"testing"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func TestShowGetYear(t *testing.T) {
	tests := []struct {
		name        string
		show        models.Show
		expected    string
		description string
	}{
		{
			name: "Valid FirstAirDate",
			show: models.Show{
				ID:           1,
				Name:         "Breaking Bad",
				FirstAirDate: "2008-01-20",
			},
			expected:    "(2008)",
			description: "Should extract year from first air date",
		},
		{
			name: "Empty FirstAirDate",
			show: models.Show{
				ID:           1,
				Name:         "Unknown Show",
				FirstAirDate: "",
			},
			expected:    "",
			description: "Should return empty string when first air date is empty",
		},
		{
			name: "Short date",
			show: models.Show{
				ID:           1,
				Name:         "Short Date Show",
				FirstAirDate: "99",
			},
			expected:    "",
			description: "Should return empty string when date is shorter than 4 chars",
		},
		{
			name: "Exactly 4 character date",
			show: models.Show{
				ID:           1,
				Name:         "Year Only Show",
				FirstAirDate: "2015",
			},
			expected:    "(2015)",
			description: "Should handle exactly 4 character dates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.show.GetYear()
			if result != tt.expected {
				t.Errorf("%s: expected '%s', got '%s'", tt.description, tt.expected, result)
			}
		})
	}
}

func TestShowGetTitle(t *testing.T) {
	tests := []struct {
		name        string
		show        models.Show
		expected    string
		description string
	}{
		{
			name: "With Name",
			show: models.Show{
				ID:   1,
				Name: "Game of Thrones",
			},
			expected:    "Game of Thrones",
			description: "Should return the Name field",
		},
		{
			name: "Empty Name",
			show: models.Show{
				ID:   1,
				Name: "",
			},
			expected:    "",
			description: "Should return empty string when Name is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.show.GetTitle()
			if result != tt.expected {
				t.Errorf("%s: expected '%s', got '%s'", tt.description, tt.expected, result)
			}
		})
	}
}

func TestShowStructFields(t *testing.T) {
	show := models.Show{
		ID:               1,
		Name:             "Stranger Things",
		OriginalName:     "Stranger Things",
		Overview:         "A mysterious incident causes a town to erupt in terror",
		FirstAirDate:     "2016-07-15",
		VoteAverage:      8.6,
		VoteCount:        500000,
		OriginalLanguage: "en",
	}

	if show.ID != 1 {
		t.Errorf("Expected ID 1, got %d", show.ID)
	}

	if show.Name != "Stranger Things" {
		t.Errorf("Expected Name 'Stranger Things', got '%s'", show.Name)
	}

	if show.OriginalName != "Stranger Things" {
		t.Errorf("Expected OriginalName 'Stranger Things', got '%s'", show.OriginalName)
	}

	if show.VoteAverage != 8.6 {
		t.Errorf("Expected VoteAverage 8.6, got %f", show.VoteAverage)
	}

	if show.VoteCount != 500000 {
		t.Errorf("Expected VoteCount 500000, got %d", show.VoteCount)
	}

	if show.OriginalLanguage != "en" {
		t.Errorf("Expected OriginalLanguage 'en', got '%s'", show.OriginalLanguage)
	}
}

func TestShowDiscoverResponse(t *testing.T) {
	show1 := models.Show{
		ID:           1,
		Name:         "Breaking Bad",
		FirstAirDate: "2008-01-20",
	}

	show2 := models.Show{
		ID:           2,
		Name:         "The Office",
		FirstAirDate: "2005-03-24",
	}

	response := models.ShowDiscoverResponse{
		Page:         1,
		Results:      []models.Show{show1, show2},
		TotalPages:   5,
		TotalResults: 100,
	}

	if response.Page != 1 {
		t.Errorf("Expected Page 1, got %d", response.Page)
	}

	if len(response.Results) != 2 {
		t.Errorf("Expected 2 shows, got %d", len(response.Results))
	}

	if response.TotalPages != 5 {
		t.Errorf("Expected TotalPages 5, got %d", response.TotalPages)
	}

	if response.TotalResults != 100 {
		t.Errorf("Expected TotalResults 100, got %d", response.TotalResults)
	}

	if response.Results[0].Name != "Breaking Bad" {
		t.Errorf("Expected first show 'Breaking Bad', got '%s'", response.Results[0].Name)
	}
}
