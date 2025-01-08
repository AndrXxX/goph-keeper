package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/client/app"
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/services/auth"
	"github.com/AndrXxX/goph-keeper/internal/client/services/dbprovider"
	"github.com/AndrXxX/goph-keeper/internal/client/services/ormstorages"
	"github.com/AndrXxX/goph-keeper/internal/client/services/storageadapters"
	"github.com/AndrXxX/goph-keeper/internal/client/state"
	"github.com/AndrXxX/goph-keeper/internal/client/views"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/requestsender"
	"github.com/AndrXxX/goph-keeper/pkg/urlbuilder"
)

func main() {
	_ = logger.Initialize("debug", []string{"./client.log"})
	ap := &auth.Provider{
		Sender: requestsender.New(&http.Client{}),
		UB:     urlbuilder.New("http://localhost:8081"),
	}
	ctx := context.Background()
	dbProvider := &dbprovider.DBProvider{}
	db, err := dbProvider.DB()
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	sp := ormstorages.Factory()
	sa := storageadapters.Factory{}
	appState := &state.AppState{
		User:       &entities.User{},
		DBProvider: dbProvider,
		Storages: &state.Storages{
			User:     sa.ORMUserAdapter(sp.User(ctx, db)),
			Password: sa.ORMPasswordsAdapter(sp.Password(ctx, db)),
			Note:     sa.ORMNotesAdapter(sp.Note(ctx, db)),
			BankCard: sa.ORMBankCardAdapter(sp.BankCard(ctx, db)),
		},
	}
	viewsFactory := views.Factory{
		AppState:   appState,
		Loginer:    ap,
		Registerer: ap,
	}
	if err := app.New(views.NewMap(viewsFactory), appState).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
