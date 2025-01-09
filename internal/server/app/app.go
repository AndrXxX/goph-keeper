package app

import (
	"context"
	"errors"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vingarcia/ksql"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	"github.com/AndrXxX/goph-keeper/internal/server/api/middlewares"
	"github.com/AndrXxX/goph-keeper/internal/server/config"
	"github.com/AndrXxX/goph-keeper/internal/server/controllers"
	entities2 "github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/services/entityconvertors"
	"github.com/AndrXxX/goph-keeper/internal/server/services/valueconvertors"
	"github.com/AndrXxX/goph-keeper/pkg/hashgenerator"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/requestjsonentity"
	"github.com/AndrXxX/goph-keeper/pkg/token"
)

const shutdownTimeout = 5 * time.Second

type app struct {
	config  appConfig
	storage Storage
}

func New(c *config.Config, s Storage) *app {
	return &app{
		config:  appConfig{c},
		storage: s,
	}
}

func (a *app) Run(commonCtx context.Context) error {
	srv := a.runServer(a.getRouter())

	<-commonCtx.Done()
	return a.shutdown(srv)
}

func (a *app) getRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	a.registerAPI(r)
	return r
}

func (a *app) runServer(r *chi.Mux) *http.Server {
	srv := &http.Server{Addr: a.config.c.Host, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("HTTP server ListenAndServe", zap.Error(err))
		}
	}()

	logger.Log.Info("listening", zap.String("host", a.config.c.Host))
	return srv
}

func (a *app) shutdown(srv *http.Server) error {
	logger.Log.Info("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	shutdown := make(chan struct{}, 1)
	go func() {
		if db, ok := a.storage.DB.(ksql.DB); ok {
			_ = db.Close()
		}
		shutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("server shutdown: %w", shutdownCtx.Err())
	case <-shutdown:
		log.Println("finished")
	}

	return nil
}

func (a *app) registerAPI(r *chi.Mux) {
	hg := hashgenerator.Factory().SHA256(a.config.c.PasswordKey)
	ts := token.New(a.config.c.AuthKey, time.Duration(a.config.c.AuthKeyExpired)*time.Second)

	r.Use(middlewares.RequestLogger().Handler)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.CompressGzip().Handler)
		ac := controllers.AuthController{US: a.storage.US, HG: hg, TS: ts, UF: &requestjsonentity.Fetcher[entities2.User]{}}
		r.Post("/api/user/register", ac.Register)
		r.Post("/api/user/login", ac.Login)
	})

	r.Group(func(r chi.Router) {
		ecf := entityconvertors.Factory{}
		vcf := valueconvertors.Factory{}

		r.Use(middlewares.IsAuthorized(ts).Handler)
		r.Use(middlewares.CompressGzip().Handler)
		lpc := controllers.ItemsController[entities2.PasswordItem]{
			Type:      datatypes.Passwords,
			Fetcher:   &requestjsonentity.Fetcher[entities2.PasswordItem]{},
			Storage:   a.storage.IS,
			Convertor: ecf.Password(vcf.Password()),
		}
		r.Post(fmt.Sprintf("/api/updates/%s", datatypes.Passwords), lpc.StoreUpdates)
		r.Get(fmt.Sprintf("/api/updates/%s", datatypes.Passwords), lpc.FetchUpdates)

		tc := controllers.ItemsController[entities2.NoteItem]{
			Type:      datatypes.Notes,
			Fetcher:   &requestjsonentity.Fetcher[entities2.NoteItem]{},
			Storage:   a.storage.IS,
			Convertor: ecf.Note(vcf.Note()),
		}
		r.Post(fmt.Sprintf("/api/updates/%s", datatypes.Notes), tc.StoreUpdates)
		r.Get(fmt.Sprintf("/api/updates/%s", datatypes.Notes), tc.FetchUpdates)

		bcc := controllers.ItemsController[entities2.BankCardItem]{
			Type:      datatypes.BankCards,
			Fetcher:   &requestjsonentity.Fetcher[entities2.BankCardItem]{},
			Storage:   a.storage.IS,
			Convertor: ecf.BankCard(vcf.BankCard()),
		}
		r.Post(fmt.Sprintf("/api/updates/%s", datatypes.BankCards), bcc.StoreUpdates)
		r.Get(fmt.Sprintf("/api/updates/%s", datatypes.BankCards), bcc.FetchUpdates)

		bc := controllers.ItemsController[entities2.FileItem]{
			Type:      datatypes.Files,
			Fetcher:   &requestjsonentity.Fetcher[entities2.FileItem]{},
			Storage:   a.storage.IS,
			Convertor: ecf.File(vcf.File()),
		}
		r.Post(fmt.Sprintf("/api/updates/%s", datatypes.Files), bc.StoreUpdates)
		r.Get(fmt.Sprintf("/api/updates/%s", datatypes.Files), bc.FetchUpdates)
	})
}
