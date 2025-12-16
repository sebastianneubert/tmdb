package display

import (
	"fmt"
	"strings"
)

type MovieDisplay struct {
	Number       int
	Title        string
	EnglishTitle string
	Year         string
	Rating       float64
	Votes        int
	Providers    []string
	TmdbID       int
	ImdbID       string
	Overview     string
	Character    string
	Genres       []string
}

func DisplayMovie(m MovieDisplay) {
	fmt.Println(SeparatorStyle.Render(strings.Repeat("=", 60)))

	englishTitleDisplay := ""
	if m.Title != m.EnglishTitle && m.EnglishTitle != "" {
		englishTitleDisplay = OriginalTitleStyle.Render(" (" + m.EnglishTitle + ")")
	}

	fmt.Printf("%d. %s%s %s\n", m.Number, TitleStyle.Render(m.Title), englishTitleDisplay, m.Year)
	fmt.Printf("   Rating: %s/10 (Votes: %d)\n", RatingStyle.Render(fmt.Sprintf("%.1f", m.Rating)), m.Votes)

	if len(m.Genres) > 0 {
		genreStr := strings.Join(m.Genres, ", ")
		fmt.Printf("   Genres: %s\n", OriginalTitleStyle.Render(genreStr))
	}

	if m.Character != "" {
		fmt.Printf("   Character: %s\n", m.Character)
	}

	styledProviders := make([]string, len(m.Providers))
	for i, p := range m.Providers {
		styledProviders[i] = ProviderStyle.Render(p)
	}
	fmt.Printf("   STREAMING on: %s\n", strings.Join(styledProviders, ", "))

	fmt.Printf("   TMDb Details: https://www.themoviedb.org/movie/%d\n", m.TmdbID)
	if m.ImdbID != "" {
		fmt.Printf("   IMDb Details: https://www.imdb.com/title/%s/\n", m.ImdbID)
	}

	fmt.Printf("   Overview: %s\n", truncateString(m.Overview, 100))
}

func DisplaySeparator() {
	fmt.Println(SeparatorStyle.Render(strings.Repeat("=", 60)))
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return strings.TrimSpace(s[:maxLen]) + "..."
	}
	return s
}

type ShowDisplay struct {
	Number       int
	Title        string
	EnglishTitle string
	Year         string
	Rating       float64
	Votes        int
	Providers    []string
	TmdbID       int
	ImdbID       string
	TvdbID       int
	Overview     string
}

type ActorDisplay struct {
	Number      int
	Name        string
	Popularity  float64
	TmdbID      int
	ProfilePath string
}

func DisplayActor(a ActorDisplay) {
	fmt.Println(SeparatorStyle.Render(strings.Repeat("-", 60)))
	fmt.Printf("%d. %s\n", a.Number, ActorNameStyle.Render(a.Name))
	fmt.Printf("   Popularity: %s\n", PopularityStyle.Render(fmt.Sprintf("%.1f", a.Popularity)))
	fmt.Printf("   TMDb Profile: https://www.themoviedb.org/person/%d\n", a.TmdbID)
}

func DisplayShow(s ShowDisplay) {
	fmt.Println(SeparatorStyle.Render(strings.Repeat("=", 60)))

	englishTitleDisplay := ""
	if s.Title != s.EnglishTitle && s.EnglishTitle != "" {
		englishTitleDisplay = OriginalTitleStyle.Render(" (" + s.EnglishTitle + ")")
	}

	fmt.Printf("%d. %s%s %s\n", s.Number, TitleStyle.Render(s.Title), englishTitleDisplay, s.Year)
	fmt.Printf("   Rating: %s/10 (Votes: %d)\n", RatingStyle.Render(fmt.Sprintf("%.1f", s.Rating)), s.Votes)

	styledProviders := make([]string, len(s.Providers))
	for i, p := range s.Providers {
		styledProviders[i] = ProviderStyle.Render(p)
	}
	fmt.Printf("   STREAMING on: %s\n", strings.Join(styledProviders, ", "))

	fmt.Printf("   TMDb Details: https://www.themoviedb.org/tv/%d\n", s.TmdbID)
	if s.ImdbID != "" {
		fmt.Printf("   IMDb Details: https://www.imdb.com/title/%s/\n", s.ImdbID)
	}
	if s.TvdbID > 0 {
		fmt.Printf("   TVDB Details: https://thetvdb.com/?tab=series&id=%d\n", s.TvdbID)
	}

	fmt.Printf("   Overview: %s\n", truncateString(s.Overview, 100))
}
