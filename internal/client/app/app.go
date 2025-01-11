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
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
	"github.com/AndrXxX/goph-keeper/pkg/hashgenerator"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/queue"
	"github.com/AndrXxX/goph-keeper/pkg/requestsender"
	"github.com/AndrXxX/goph-keeper/pkg/urlbuilder"
)

const queueTimeout = 1 * time.Second
const msgTimeout = 2 * time.Second
const shutdownTimeout = 5 * time.Second

type App struct {
	TUI   *tea.Program
	State *state.AppState
	Sync  syncManager
	QR    queueRunner
	//crypto *CryptoManager
	c  *config.Config
	vf *views.Factory
	ua *useraccessor.Accessor
}

func NewApp(c *config.Config) *App {
	app := App{
		State: &state.AppState{},
		c:     c,
		QR:    queue.NewRunner(queueTimeout).SetWorkersCount(5),
	}
	ub := urlbuilder.New(c.Host)
	ap := &auth.Provider{Sender: requestsender.New(&http.Client{}), UB: ub}
	dbProvider := &dbprovider.DBProvider{}
	sp := ormstorages.Factory()
	sa := storageadapters.Factory{}
	app.ua = &useraccessor.Accessor{
		User: &entities.User{},
		SP: func(db *gorm.DB) useraccessor.Storage[entities.User] {
			return sa.ORMUserAdapter(sp.User(context.Background(), db))
		},
		SDB: func(masterPass string, recreate bool) (*gorm.DB, error) {
			if recreate {
				err := dbProvider.RemoveDB()
				if err != nil {
					return nil, err
				}
			}
			actDB, err := dbProvider.DB(masterPass)
			if err != nil {
				return nil, err
			}
			app.State.DB = actDB
			return actDB, nil
		},
		HG: func(key string) useraccessor.HashGenerator {
			return hashgenerator.Factory().SHA256(key)
		},
	}
	app.vf = &views.Factory{Loginer: ap, Registerer: ap}
	return &app
}

func (a *App) Run(ctx context.Context) error {
	ctx, stop := context.WithCancel(ctx)
	a.TUI = tea.NewProgram(a.vf.Container(
		views.WithBuildInfo(a.c),
		views.WithStartView(names.AuthMenu),
		views.WithMap(views.AuthMap(a.vf)),
		views.WithShowMessage(msgTimeout),
		views.WithShowError(msgTimeout),
		views.WithUpdateUser(a.ua),
		views.WithAuth(a.ua),
		views.WithQuit(func() {
			stop()
		}),
	), tea.WithAltScreen())
	a.ua.AfterAuth = func() {
		a.TUI.Kill()
		if err := a.runFull(ctx); err != nil {
			logger.Log.Error(err.Error())
		}
	}
	go a.runTIU()
	<-ctx.Done()
	return a.shutdown()
}

func (a *App) runFull(ctx context.Context) error {
	ctx, stop := context.WithCancel(ctx)

	rs := requestsender.New(&http.Client{}, requestsender.WithToken(a.ua.GetUser().Token))
	ub := urlbuilder.New(a.c.Host)
	sa := storageadapters.Factory{}
	sp := ormstorages.Factory()

	us := sa.ORMUserAdapter(sp.User(ctx, a.State.DB))
	ps := sa.ORMPasswordsAdapter(sp.Password(ctx, a.State.DB))
	ns := sa.ORMNotesAdapter(sp.Note(ctx, a.State.DB))
	bs := sa.ORMBankCardAdapter(sp.BankCard(ctx, a.State.DB))

	sFactory := synchronize.Factory{
		RS:       rs,
		UB:       ub,
		Storages: &synchronize.Storages{Password: ps, Note: ns, BankCard: bs},
	}
	sm := &synchronize.SyncManager{Synchronizers: sFactory.Map(), TR: func() {
		token, err := a.vf.Loginer.Login(a.ua.GetUser())
		if err != nil {
			logger.Log.Error("failed to refresh token", zap.Error(err))
			return
		}
		a.ua.SetToken(token)
		_ = us.Update(a.ua.GetUser())
		*rs = *requestsender.New(&http.Client{}, requestsender.WithToken(token))
	}}
	a.vf.S = &vContract.Storages{Password: ps, Note: ns, BankCard: bs}

	a.TUI = tea.NewProgram(a.vf.Container(
		views.WithBuildInfo(a.c),
		views.WithStartView(names.MainMenu),
		views.WithMap(views.NewMainMap(a.vf)),
		views.WithShowMessage(msgTimeout),
		views.WithShowError(msgTimeout),
		views.WithUploadItemUpdates(sm, a.QR),
		views.WithQuit(func() {
			stop()
		}),
	), tea.WithAltScreen())

	go a.runQueue(ctx)
	go a.runTIU()
	return nil
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
