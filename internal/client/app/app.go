package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/config"
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/services/auth"
	"github.com/AndrXxX/goph-keeper/internal/client/services/dbprovider"
	"github.com/AndrXxX/goph-keeper/internal/client/services/ormstorages"
	"github.com/AndrXxX/goph-keeper/internal/client/services/storageadapters"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize"
	"github.com/AndrXxX/goph-keeper/internal/client/services/useraccessor"
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views"
	vContract "github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/pkg/hashgenerator"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/queue"
	"github.com/AndrXxX/goph-keeper/pkg/requestsender"
	"github.com/AndrXxX/goph-keeper/pkg/urlbuilder"
)

const msgTimeout = 2 * time.Second
const shutdownTimeout = 5 * time.Second

type App struct {
	TUI   *tea.Program
	State *state.AppState
	Sync  syncManager
	QR    queueRunner
	//crypto *CryptoManager
	c *config.Config
}

func NewApp(c *config.Config) *App {
	ctx, stop := context.WithCancel(context.Background())
	ub := urlbuilder.New(c.Host)
	ap := &auth.Provider{Sender: requestsender.New(&http.Client{}), UB: ub}
	dbProvider := &dbprovider.DBProvider{}
	appState := &state.AppState{}
	sp := ormstorages.Factory()
	sa := storageadapters.Factory{}
	rs := requestsender.New(&http.Client{})
	ua := &useraccessor.Accessor{
		User: &entities.User{},
		SP: func(db *gorm.DB) useraccessor.Storage[entities.User] {
			return sa.ORMUserAdapter(sp.User(ctx, db))
		},
		ST: func(token string) {
			*rs = *requestsender.New(&http.Client{}, requestsender.WithToken(token))
		},
		SDB: func(masterPass string) (*gorm.DB, error) {
			actDB, err := dbProvider.DB(masterPass)
			if err != nil {
				return nil, err
			}
			// TODO: отрефакторить
			*appState.DB = *actDB
			return actDB, nil
		},
		HG: func(key string) useraccessor.HashGenerator {
			return hashgenerator.Factory().SHA256(key)
		},
	}
	sFactory := synchronize.Factory{RS: rs, UB: ub, Storages: &synchronize.Storages{
		Password: sa.ORMPasswordsAdapter(sp.Password(ctx, appState.DB)),
		Note:     sa.ORMNotesAdapter(sp.Note(ctx, appState.DB)),
		BankCard: sa.ORMBankCardAdapter(sp.BankCard(ctx, appState.DB)),
	}}
	sm := &synchronize.SyncManager{Synchronizers: sFactory.Map(), TR: func() {
		token, err := ap.Login(ua.GetUser())
		if err != nil {
			logger.Log.Error("failed to refresh token", zap.Error(err))
			return
		}
		ua.SetToken(token)
		_ = ua.US.Update(ua.GetUser())
		*rs = *requestsender.New(&http.Client{}, requestsender.WithToken(token))
	}}
	qr := queue.NewRunner(1 * time.Second).SetWorkersCount(5)
	viewsFactory := views.Factory{
		Loginer:    ap,
		Registerer: ap,
		S: &vContract.Storages{
			Password: sa.ORMPasswordsAdapter(sp.Password(ctx, appState.DB)),
			Note:     sa.ORMNotesAdapter(sp.Note(ctx, appState.DB)),
			BankCard: sa.ORMBankCardAdapter(sp.BankCard(ctx, appState.DB)),
		},
	}
	application := App{
		TUI: tea.NewProgram(viewsFactory.Container(
			views.WithShowMessage(msgTimeout),
			views.WithShowError(msgTimeout),
			views.WithUpdateUser(ua),
			views.WithAuth(ua),
			views.WithUploadItemUpdates(sm, qr),
			views.WithQuit(func() {
				stop()
			}),
		), tea.WithAltScreen()),
		State: appState,
		Sync:  sm,
		QR:    qr,
		c:     c,
	}
	return &application
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
