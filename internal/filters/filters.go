package filters

import (
	"strings"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func MeetsRatingCriteria(voteAverage float64, voteCount int, minRating float64, minVotes int) bool {
	return voteAverage >= minRating && voteCount >= minVotes
}

func ParseProviders(input string) map[string]bool {
	parts := strings.Split(input, ",")
	providerMap := make(map[string]bool)
	for _, part := range parts {
		key := strings.ToLower(strings.TrimSpace(part))
		providerMap[key] = true
	}
	return providerMap
}

func CheckAvailability(providerData models.RegionProviders, desiredProviders map[string]bool) ([]string, bool) {
	available := []string{}

	for _, p := range providerData.Flatrate {
		if isProviderMatched(p.ProviderName, desiredProviders) {
			available = append(available, p.ProviderName)
		}
	}

	return available, len(available) > 0
}

func isProviderMatched(providerName string, desiredProviders map[string]bool) bool {
	providerLower := strings.ToLower(providerName)

	for desired := range desiredProviders {
		if desired == "amazon" || desired == "amazonprime" {
			if strings.Contains(providerLower, "amazon prime video") {
				return true
			}
		} else if desired == providerLower {
			return true
		}
	}

	return false
}