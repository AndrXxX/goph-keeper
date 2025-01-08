package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AndrXxX/goph-keeper/internal/client/state"
)

type App struct {
	TUI   *tea.Program
	State *state.AppState
	Sync  syncManager
	QR    queueRunner
	//crypto *CryptoManager
}

func (a App) Run() error {
	_, err := a.TUI.Run()
	return err
}
