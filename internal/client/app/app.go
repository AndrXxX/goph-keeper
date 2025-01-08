package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views"
)

type App struct {
	TUI   *tea.Program
	State *state.AppState
	//sync   *SyncManager
	//crypto *CryptoManager
}

func New(l views.Map, s *state.AppState) *App {
	p := tea.NewProgram(views.NewContainer(l), tea.WithAltScreen())
	return &App{
		TUI:   p,
		State: s,
	}
}

func (a App) Run() error {
	_, err := a.TUI.Run()
	return err
}
