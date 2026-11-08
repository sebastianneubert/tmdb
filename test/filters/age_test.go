package filters_test

import (
	"testing"

	"github.com/sebastianneubert/tmdb/internal/filters"
	"github.com/sebastianneubert/tmdb/internal/models"
)

func TestParseCertificationAge(t *testing.T) {
	tests := []struct {
		input   string
		wantAge int
		wantOK  bool
	}{
		// Germany (FSK)
		{"0", 0, true},
		{"6", 6, true},
		{"12", 12, true},
		{"16", 16, true},
		{"18", 18, true},
		// UK (BBFC)
		{"U", 0, true},
		{"PG", 10, true},
		{"12A", 12, true},
		{"15", 15, true},
		{"R18", 18, true},
		// US (MPA)
		{"G", 0, true},
		{"PG-13", 13, true},
		{"R", 17, true},
		{"NC-17", 18, true},
		// Variants
		{"16+", 16, true},
		{"FSK 12", 12, true},
		{"  12  ", 12, true},
		{"pg-13", 13, true},
		// Unknown
		{"", 0, false},
		{"   ", 0, false},
		{"NR", 0, false},
		{"Not Rated", 0, false},
		{"A", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			age, ok := filters.ParseCertificationAge(tt.input)
			if ok != tt.wantOK || age != tt.wantAge {
				t.Errorf("ParseCertificationAge(%q) = (%d, %v), want (%d, %v)", tt.input, age, ok, tt.wantAge, tt.wantOK)
			}
		})
	}
}

func TestCertificationAgeForRegion(t *testing.T) {
	data := &models.ReleaseDatesResponse{
		Results: []models.CountryReleaseDates{
			{
				ISO31661: "US",
				ReleaseDates: []models.ReleaseDateEntry{
					{Certification: "R"},
				},
			},
			{
				ISO31661: "DE",
				ReleaseDates: []models.ReleaseDateEntry{
					{Certification: ""},
					{Certification: "12", Type: 3},
					{Certification: "16", Type: 5},
				},
			},
			{
				ISO31661: "FR",
				ReleaseDates: []models.ReleaseDateEntry{
					{Certification: ""},
				},
			},
		},
	}

	t.Run("uses strictest certification of the region", func(t *testing.T) {
		age, ok := filters.CertificationAgeForRegion(data, "DE")
		if !ok || age != 16 {
			t.Errorf("got (%d, %v), want (16, true)", age, ok)
		}
	})

	t.Run("region is case-insensitive", func(t *testing.T) {
		age, ok := filters.CertificationAgeForRegion(data, "us")
		if !ok || age != 17 {
			t.Errorf("got (%d, %v), want (17, true)", age, ok)
		}
	})

	t.Run("region with only empty certifications is unknown", func(t *testing.T) {
		if _, ok := filters.CertificationAgeForRegion(data, "FR"); ok {
			t.Error("expected unknown age for FR")
		}
	})

	t.Run("missing region is unknown", func(t *testing.T) {
		if _, ok := filters.CertificationAgeForRegion(data, "JP"); ok {
			t.Error("expected unknown age for JP")
		}
	})

	t.Run("nil data is unknown", func(t *testing.T) {
		if _, ok := filters.CertificationAgeForRegion(nil, "DE"); ok {
			t.Error("expected unknown age for nil data")
		}
	})
}

func TestMeetsAgeLimit(t *testing.T) {
	tests := []struct {
		name           string
		age            int
		known          bool
		maxAge         int
		includeUnrated bool
		want           bool
	}{
		{"filter off shows everything", 18, true, 0, false, true},
		{"filter off shows unrated", 0, false, 0, false, true},
		{"age below limit", 6, true, 12, false, true},
		{"age equal to limit", 12, true, 12, false, true},
		{"age above limit", 16, true, 12, false, false},
		{"unrated hidden by default", 0, false, 12, false, false},
		{"unrated shown with include-unrated", 0, false, 12, true, true},
		{"include-unrated does not show too old movies", 18, true, 12, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filters.MeetsAgeLimit(tt.age, tt.known, tt.maxAge, tt.includeUnrated)
			if got != tt.want {
				t.Errorf("MeetsAgeLimit(%d, %v, %d, %v) = %v, want %v", tt.age, tt.known, tt.maxAge, tt.includeUnrated, got, tt.want)
			}
		})
	}
}
