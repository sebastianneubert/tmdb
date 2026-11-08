package models

// ReleaseDateEntry is a single release of a movie in one country,
// including the age certification (e.g. "12" for FSK 12) of that release.
type ReleaseDateEntry struct {
	Certification string `json:"certification"`
	Type          int    `json:"type"`
}

// CountryReleaseDates holds all releases of a movie in one country.
type CountryReleaseDates struct {
	ISO31661     string             `json:"iso_3166_1"`
	ReleaseDates []ReleaseDateEntry `json:"release_dates"`
}

// ReleaseDatesResponse is the response of /movie/{id}/release_dates.
type ReleaseDatesResponse struct {
	ID      int                   `json:"id"`
	Results []CountryReleaseDates `json:"results"`
}
