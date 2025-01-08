package main

import (
	"context"
	"net/http"
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
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/queue"
	"github.com/AndrXxX/goph-keeper/pkg/requestsender"
	"github.com/AndrXxX/goph-keeper/pkg/urlbuilder"
)

func main() {
	_ = logger.Initialize("debug", []string{"./client.log"})
	ub := urlbuilder.New("http://localhost:8081")
	ap := &auth.Provider{Sender: requestsender.New(&http.Client{}), UB: ub}
	ctx := context.Background()
	dbProvider := &dbprovider.DBProvider{}
	db, err := dbProvider.DB()
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	sp := ormstorages.Factory()
	sa := storageadapters.Factory{}
	rs := requestsender.New(&http.Client{})
	appState := &state.AppState{
		User:       &entities.User{},
		DBProvider: dbProvider,
		Storages: &state.Storages{
			User:     sa.ORMUserAdapter(sp.User(ctx, db)),
			Password: sa.ORMPasswordsAdapter(sp.Password(ctx, db)),
			Note:     sa.ORMNotesAdapter(sp.Note(ctx, db)),
			BankCard: sa.ORMBankCardAdapter(sp.BankCard(ctx, db)),
		},
		AS: func(u *entities.User) {
			*rs = *requestsender.New(&http.Client{}, requestsender.WithToken(u.Token))
		},
	}
	sFactory := synchronize.Factory{RS: rs, UB: ub, Storages: (*synchronize.Storages)(appState.Storages)}
	sm := &synchronize.SyncManager{Synchronizers: sFactory.Map()}
	viewsFactory := views.Factory{
		AppState:   appState,
		Loginer:    ap,
		Registerer: ap,
		SM:         sm,
	}
	application := app.App{
		TUI:   tea.NewProgram(views.NewContainer(views.NewMap(viewsFactory)), tea.WithAltScreen()),
		State: appState,
		Sync:  sm,
		QR:    queue.NewRunner(1 * time.Second),
	}
	if err := application.Run(); err != nil {
		logger.Log.Fatal("failed to start application", zap.Error(err))
	}
}
