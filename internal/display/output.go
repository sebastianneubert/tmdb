package display

import (
	"fmt"
	"strings"
)

// PrintSearchStartMessage prints the initial search/query message
func PrintSearchStartMessage(searchType string, minRating float64, minVotes int, providers, region string) {
	fmt.Printf("Searching TMDb's %s...\n", searchType)
	fmt.Printf("Criteria: Min Rating: %.1f | Min Votes: %d\n", minRating, minVotes)
	fmt.Printf("Filtering for [%s] in region [%s]\n\n", providers, strings.ToUpper(region))
}

// PrintSearchResultsSummary prints the final results summary
func PrintSearchResultsSummary(searchType string, resultsFound int) {
	DisplaySeparator()
	if resultsFound == 0 {
		fmt.Printf("No %s found matching criteria.\n", searchType)
	} else {
		fmt.Printf("Displayed %d %s.\n", resultsFound, searchType)
	}
}

// PrintSearchNoResults prints a detailed "no results" message for search queries
func PrintSearchNoResults(query string, moviesChecked int, minRating float64, minVotes int) {
	DisplaySeparator()
	fmt.Printf("No movies found for \"%s\" that meet criteria and are available on your providers.\n", query)
	fmt.Printf("(Checked %d movies from search results)\n", moviesChecked)
	fmt.Println("\nTry:")
	fmt.Printf("  - Lowering --min-rating (current: %.1f)\n", minRating)
	fmt.Printf("  - Lowering --min-votes (current: %d)\n", minVotes)
	fmt.Printf("  - Adding more --providers\n")
}

// PrintSearchCompleteMessage prints the completion message for search results
func PrintSearchCompleteMessage(resultsFound, moviesChecked int) {
	DisplaySeparator()
	fmt.Printf("Search complete: Displayed %d movies (out of %d checked).\n", resultsFound, moviesChecked)
}
