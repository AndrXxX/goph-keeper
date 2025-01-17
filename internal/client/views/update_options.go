package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateOption func(msg tea.Msg) (tea.Model, tea.Cmd)
