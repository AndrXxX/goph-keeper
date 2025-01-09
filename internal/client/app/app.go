package app

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

const shutdownTimeout = 5 * time.Second

type App struct {
	TUI   *tea.Program
	State *state.AppState
	Sync  syncManager
	QR    queueRunner
	//crypto *CryptoManager
}

func (a *App) Run(ctx context.Context) error {
	go a.runQueue(ctx)
	go a.runTIU()

	<-ctx.Done()
	return a.shutdown()
}

func (a *App) runQueue(ctx context.Context) {
	err := a.QR.Run(ctx)
	if err != nil {
		logger.Log.Error("start queue runner: %w", zap.Error(err))
	}
}

func (a *App) runTIU() {
	_, err := a.TUI.Run()
	if err != nil {
		logger.Log.Error("start TIU: %w", zap.Error(err))
	}
}

func (a *App) shutdown() error {
	logger.Log.Info("shutting down client gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	shutdown := make(chan struct{}, 1)
	go func() {
		err := a.QR.Stop(shutdownCtx)
		if err != nil {
			logger.Log.Error("shutdown queue", zap.Error(err))
		}
		a.TUI.Quit()
		shutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("client shutdown: %w", shutdownCtx.Err())
	case <-shutdown:
		logger.Log.Info("finished")
	}
	return nil
}
