package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/client/app"
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/services/auth"
	"github.com/AndrXxX/goph-keeper/internal/client/services/dbprovider"
	"github.com/AndrXxX/goph-keeper/internal/client/services/ormstorages"
	"github.com/AndrXxX/goph-keeper/internal/client/services/storageadapters"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize"
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views"
	vContract "github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/queue"
	"github.com/AndrXxX/goph-keeper/pkg/requestsender"
	"github.com/AndrXxX/goph-keeper/pkg/urlbuilder"
	_ "github.com/AndrXxX/goph-keeper/pkg/validators"
)

const msgTimeout = 2 * time.Second

func main() {
	_ = logger.Initialize("debug", []string{"./client.log"})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	ub := urlbuilder.New("http://localhost:8081")
	ap := &auth.Provider{Sender: requestsender.New(&http.Client{}), UB: ub}
	dbProvider := &dbprovider.DBProvider{}
	db, err := dbProvider.DB("testKey")
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	sp := ormstorages.Factory()
	sa := storageadapters.Factory{}
	rs := requestsender.New(&http.Client{})
	appState := &state.AppState{
		User:       &entities.User{},
		DBProvider: dbProvider,
		Storages:   &state.Storages{User: sa.ORMUserAdapter(sp.User(ctx, db))},
		AS: func(u *entities.User) {
			*rs = *requestsender.New(&http.Client{}, requestsender.WithToken(u.Token))
		},
	}
	sFactory := synchronize.Factory{RS: rs, UB: ub, Storages: &synchronize.Storages{
		Password: sa.ORMPasswordsAdapter(sp.Password(ctx, db)),
		Note:     sa.ORMNotesAdapter(sp.Note(ctx, db)),
		BankCard: sa.ORMBankCardAdapter(sp.BankCard(ctx, db)),
	}}
	sm := &synchronize.SyncManager{Synchronizers: sFactory.Map(), TR: func() {
		token, err := ap.Login(appState.User)
		if err != nil {
			logger.Log.Error("failed to refresh token", zap.Error(err))
			return
		}
		appState.User.Token = token
		_ = appState.Storages.User.Update(appState.User)
		*rs = *requestsender.New(&http.Client{}, requestsender.WithToken(token))
	}}
	qr := queue.NewRunner(1 * time.Second).SetWorkersCount(5)
	viewsFactory := views.Factory{
		Loginer:    ap,
		Registerer: ap,
		S: &vContract.Storages{
			User:     sa.ORMUserAdapter(sp.User(ctx, db)),
			Password: sa.ORMPasswordsAdapter(sp.Password(ctx, db)),
			Note:     sa.ORMNotesAdapter(sp.Note(ctx, db)),
			BankCard: sa.ORMBankCardAdapter(sp.BankCard(ctx, db)),
		},
	}
	application := app.App{
		TUI: tea.NewProgram(viewsFactory.Container(
			views.WithShowMessage(msgTimeout),
			views.WithShowError(msgTimeout),
			views.WithUpdateUser(appState),
			views.WithAuth(appState),
			views.WithUploadItemUpdates(sm, qr),
		), tea.WithAltScreen()),
		State: appState,
		Sync:  sm,
		QR:    qr,
	}
	if err := application.Run(ctx); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
