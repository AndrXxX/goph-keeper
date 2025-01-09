package app

import (
	"fmt"

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

func (a *App) Run() error {
	err := a.QR.Run()
	if err != nil {
		return fmt.Errorf("start queue runner: %w", err)
	}
	_, err = a.TUI.Run()
	return err
}
