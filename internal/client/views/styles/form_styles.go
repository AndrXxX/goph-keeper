package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Focused = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	Blurred = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	Cursor  = Focused
	Empty   = lipgloss.NewStyle()
	Help    = Blurred
)
