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
	"github.com/AndrXxX/goph-keeper/internal/client/jobs"
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
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	"github.com/AndrXxX/goph-keeper/pkg/filestorage"
	"github.com/AndrXxX/goph-keeper/pkg/hashgenerator"
	"github.com/AndrXxX/goph-keeper/pkg/httpclient"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/queue"
	"github.com/AndrXxX/goph-keeper/pkg/requestsender"
	"github.com/AndrXxX/goph-keeper/pkg/tlsconfig"
	"github.com/AndrXxX/goph-keeper/pkg/urlbuilder"
)

const queueTimeout = 1 * time.Second
const msgTimeout = 2 * time.Second
const shutdownTimeout = 5 * time.Second
const syncInterval = 10 * time.Second

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
	ap := &auth.Provider{
		Sender: requestsender.New(&http.Client{}), UB: ub,
		KS: &auth.KeySaver{KeyPath: c.ServerKeyPath},
	}
	dbProvider := &dbprovider.DBProvider{Path: c.DBPath}
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
			app.State.MasterPass = masterPass
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
			a.TUI.Kill()
			stop()
		}),
	), tea.WithAltScreen(), tea.WithContext(ctx))
	a.ua.AfterAuth = func() {
		a.TUI.Kill()
		if err := a.runFull(ctx); err != nil {
			logger.Log.Error(err.Error())
		}
		stop()
	}
	go a.runTIU()
	<-ctx.Done()
	return nil
}

func (a *App) runFull(ctx context.Context) error {
	ctx, stop := context.WithCancel(ctx)
	cp := httpclient.Provider{ConfProvider: tlsconfig.NewProvider(a.c.ServerKeyPath)}
	client, err := cp.Fetch()
	if err != nil {
		stop()
		return err
	}
	rs := requestsender.New(client, requestsender.WithToken(a.ua.GetUser().Token))
	ub := urlbuilder.New(a.c.Host)
	sa := storageadapters.Factory{}
	sp := ormstorages.Factory()

	us := sa.ORMUserAdapter(sp.User(ctx, a.State.DB))
	ps := sa.ORMPasswordsAdapter(sp.Password(ctx, a.State.DB))
	ns := sa.ORMNotesAdapter(sp.Note(ctx, a.State.DB))
	bs := sa.ORMBankCardAdapter(sp.BankCard(ctx, a.State.DB))
	fs := sa.ORMFileAdapter(sp.File(ctx, a.State.DB))
	dfs, err := filestorage.New(a.c.FileStoragePath, hashgenerator.Factory().SHA256(a.State.MasterPass))
	if err != nil {
		stop()
		return fmt.Errorf("initialize file storage: %w", err)
	}

	sFactory := synchronize.Factory{
		RS:       rs,
		UB:       ub,
		Storages: &synchronize.Storages{Password: ps, Note: ns, BankCard: bs, File: fs, FS: dfs},
	}
	a.Sync = &synchronize.SyncManager{Synchronizers: sFactory.Map(), TR: func() {
		token, err := a.vf.Loginer.Login(a.ua.GetUser())
		if err != nil {
			logger.Log.Error("failed to refresh token", zap.Error(err))
			return
		}
		a.ua.SetToken(token)
		_ = us.Update(a.ua.GetUser())
		*rs = *requestsender.New(&http.Client{}, requestsender.WithToken(token))
	}}
	a.vf.S = &vContract.Storages{Password: ps, Note: ns, BankCard: bs, File: fs}

	a.TUI = tea.NewProgram(a.vf.Container(
		views.WithBuildInfo(a.c),
		views.WithStartView(names.MainMenu),
		views.WithMap(views.NewMainMap(a.vf)),
		views.WithShowMessage(msgTimeout),
		views.WithShowError(msgTimeout),
		views.WithUploadItemUpdates(a.Sync, a.QR),
		views.WithRepeatableJob(a.QR, syncInterval, &jobs.SyncJob{Type: datatypes.Passwords, SyncManager: a.Sync}),
		views.WithRepeatableJob(a.QR, syncInterval, &jobs.SyncJob{Type: datatypes.Notes, SyncManager: a.Sync}),
		views.WithRepeatableJob(a.QR, syncInterval, &jobs.SyncJob{Type: datatypes.BankCards, SyncManager: a.Sync}),
		views.WithRepeatableJob(a.QR, syncInterval, &jobs.SyncJob{Type: datatypes.Files, SyncManager: a.Sync}),
		views.WithDownloadFile(dfs),
		views.WithQuit(func() {
			stop()
			a.TUI.Kill()
		}),
	), tea.WithAltScreen())

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
		a.TUI.Kill()
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
