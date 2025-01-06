package views

import (
	tea "github.com/charmbracelet/bubbletea"

	kb "github.com/AndrXxX/goph-keeper/internal/client/keyboard"
)

type View interface {
	tea.Model
}

type keyboardKeys interface {
	kbKeys() kb.KeyMap
}
