package filters

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/sebastianneubert/tmdb/internal/models"
)

// namedCertificationAges maps certifications that are not plain numbers
// (US MPA, UK BBFC and a few others) to a minimum age.
var namedCertificationAges = map[string]int{
	"G":     0,
	"U":     0,
	"TP":    0,
	"AL":    0,
	"ALL":   0,
	"PG":    10,
	"PG-13": 13,
	"R":     17,
	"NC-17": 18,
	"X":     18,
	"R18":   18,
}

// numericCertification matches numeric ratings such as "12", "12A", "16+",
// "FSK 16" or "-12" and captures the age.
var numericCertification = regexp.MustCompile(`^(?:FSK|AB|R)?\s*-?\s*(\d{1,2})\s*[A+]?$`)

// ParseCertificationAge converts a TMDb certification string into the minimum age
// it stands for. The second return value is false if the string is empty or unknown.
func ParseCertificationAge(certification string) (int, bool) {
	c := strings.ToUpper(strings.TrimSpace(certification))
	if c == "" {
		return 0, false
	}

	if age, ok := namedCertificationAges[c]; ok {
		return age, true
	}

	match := numericCertification.FindStringSubmatch(c)
	if match == nil {
		return 0, false
	}

	age, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, false
	}
	return age, true
}

// CertificationAgeForRegion returns the age recommendation of a movie in a region.
// A movie can have several releases (cinema, home video, ...) with different
// certifications; the strictest (highest) one is used. The second return value is
// false if the movie has no usable certification in that region.
func CertificationAgeForRegion(data *models.ReleaseDatesResponse, region string) (int, bool) {
	if data == nil {
		return 0, false
	}

	region = strings.ToUpper(strings.TrimSpace(region))
	highest := 0
	found := false

	for _, country := range data.Results {
		if strings.ToUpper(country.ISO31661) != region {
			continue
		}
		for _, release := range country.ReleaseDates {
			age, ok := ParseCertificationAge(release.Certification)
			if !ok {
				continue
			}
			if !found || age > highest {
				highest = age
				found = true
			}
		}
	}

	return highest, found
}

// MeetsAgeLimit reports whether a movie may be shown for the given age limit.
// maxAge <= 0 disables the filter. Movies without a known age recommendation are
// only shown if includeUnrated is set.
func MeetsAgeLimit(age int, known bool, maxAge int, includeUnrated bool) bool {
	if maxAge <= 0 {
		return true
	}
	if !known {
		return includeUnrated
	}
	return age <= maxAge
}
