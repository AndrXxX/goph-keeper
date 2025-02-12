package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Focused = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	Blurred = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	Cursor  = Focused
	Empty   = lipgloss.NewStyle()
	Help    = Blurred.Margin(1, 0, 0)
	Error   = lipgloss.NewStyle().
		Background(lipgloss.Color("#f12929")).
		Foreground(lipgloss.Color("#ffffff")).
		Margin(1, 0, 1).
		Padding(1)
	Info = lipgloss.NewStyle().
		Background(lipgloss.Color("#8162e4")).
		Foreground(lipgloss.Color("#ffffff")).
		Margin(1, 0, 1).
		Padding(1)

	Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)
)
