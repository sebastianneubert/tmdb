package display

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1).
		Inline(true)

	OriginalTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true).
		Inline(true)

	RatingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FCE043")).
		Bold(true)

	SeparatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555"))

	ProviderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF")).
		Bold(true)

	ActorNameStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Inline(true)

	PopularityStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Inline(true)
)