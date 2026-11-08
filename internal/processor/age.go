package processor

import (
	"fmt"
	"strings"

	"github.com/sebastianneubert/tmdb/internal/api"
	"github.com/sebastianneubert/tmdb/internal/filters"
)

// AgeLookupFunc returns the age recommendation of a movie. The bool is false
// if the movie has no age recommendation.
type AgeLookupFunc func(movieID int) (age int, known bool, err error)

// AgeFilter hides movies whose age recommendation in a region is above MaxAge.
// All methods are safe to call on a nil *AgeFilter, which means "no age filter".
type AgeFilter struct {
	MaxAge         int
	IncludeUnrated bool
	Region         string
	Lookup         AgeLookupFunc

	ages          map[int]int // ages found by Allows, used for display
	hiddenTooOld  int
	hiddenUnrated int
}

// NewAgeFilter creates an age filter that looks up certifications via the TMDb API.
// It returns nil if maxAge <= 0, i.e. if the filter is switched off.
func NewAgeFilter(client *api.Client, region string, maxAge int, includeUnrated bool) *AgeFilter {
	if maxAge <= 0 {
		return nil
	}

	f := &AgeFilter{MaxAge: maxAge, IncludeUnrated: includeUnrated, Region: region}
	if client != nil {
		f.Lookup = func(movieID int) (int, bool, error) {
			releases, err := client.GetReleaseDates(movieID)
			if err != nil {
				return 0, false, err
			}
			age, known := filters.CertificationAgeForRegion(releases, region)
			return age, known, nil
		}
	}
	return f
}

// Enabled reports whether the filter does anything.
func (f *AgeFilter) Enabled() bool {
	return f != nil && f.MaxAge > 0 && f.Lookup != nil
}

// Allows reports whether the movie may be shown. Failed lookups count as
// "unrated", so they never let a movie through unnoticed.
func (f *AgeFilter) Allows(movieID int) bool {
	if !f.Enabled() {
		return true
	}

	age, known, err := f.Lookup(movieID)
	if err != nil {
		known = false
	}
	if known {
		if f.ages == nil {
			f.ages = make(map[int]int)
		}
		f.ages[movieID] = age
	}

	if filters.MeetsAgeLimit(age, known, f.MaxAge, f.IncludeUnrated) {
		return true
	}

	if known {
		f.hiddenTooOld++
	} else {
		f.hiddenUnrated++
	}
	return false
}

// AgeLabel returns the age recommendation of a movie as shown in the output (e.g. "12+").
// It is "" if the filter is off or the age of the movie is not known, so callers can
// simply skip empty labels. Ages are only known for movies that went through Allows.
func (f *AgeFilter) AgeLabel(movieID int) string {
	if f == nil {
		return ""
	}
	age, ok := f.ages[movieID]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%d+", age)
}

// Summary describes how many movies the filter hid, or "" if it is off.
func (f *AgeFilter) Summary() string {
	if !f.Enabled() {
		return ""
	}

	summary := fmt.Sprintf("Age filter (max. %d, region %s): %d hidden as too old", f.MaxAge, strings.ToUpper(f.Region), f.hiddenTooOld)
	if !f.IncludeUnrated {
		summary += fmt.Sprintf(", %d hidden without age rating (use --include-unrated to show them)", f.hiddenUnrated)
	}
	return summary
}
