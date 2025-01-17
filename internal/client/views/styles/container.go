package styles

import "github.com/charmbracelet/lipgloss"

var (
	Border = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
	Spinner = lipgloss.NewStyle().
		Padding(1, 0, 0)
)
