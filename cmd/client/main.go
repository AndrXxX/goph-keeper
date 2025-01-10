package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/app"
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
	var db *gorm.DB
	sp := ormstorages.Factory()
	sa := storageadapters.Factory{}
	rs := requestsender.New(&http.Client{})
	ua := &useraccessor.Accessor{
		User: &entities.User{},
		US:   sa.ORMUserAdapter(sp.User(ctx, db)),
		ST: func(token string) {
			*rs = *requestsender.New(&http.Client{}, requestsender.WithToken(token))
		},
		SDB: func(masterPass string) error {
			actDB, err := dbProvider.DB(masterPass)
			if err != nil {
				return err
			}
			// TODO: отрефакторить
			*db = *actDB
			return nil
		},
		HG: func(key string) useraccessor.HashGenerator {
			return hashgenerator.Factory().SHA256(key)
		},
	}
	appState := &state.AppState{}
	sFactory := synchronize.Factory{RS: rs, UB: ub, Storages: &synchronize.Storages{
		Password: sa.ORMPasswordsAdapter(sp.Password(ctx, db)),
		Note:     sa.ORMNotesAdapter(sp.Note(ctx, db)),
		BankCard: sa.ORMBankCardAdapter(sp.BankCard(ctx, db)),
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
			Password: sa.ORMPasswordsAdapter(sp.Password(ctx, db)),
			Note:     sa.ORMNotesAdapter(sp.Note(ctx, db)),
			BankCard: sa.ORMBankCardAdapter(sp.BankCard(ctx, db)),
		},
	}
	application := app.App{
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
	}
	if err := application.Run(ctx); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
