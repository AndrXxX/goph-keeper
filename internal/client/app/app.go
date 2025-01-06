package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/views"
)

type App struct {
	TUI   *tea.Program
	State *State
	//sync   *SyncManager
	//crypto *CryptoManager
	Views map[string]views.View
}

func New(l views.Map) *App {
	p := tea.NewProgram(views.NewContainer(l), tea.WithAltScreen())
	return &App{
		TUI: p,
	}
}

func (a App) Run() error {
	_, err := a.TUI.Run()
	return err
}
