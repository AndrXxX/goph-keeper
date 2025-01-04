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
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
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
		ac := controllers.AuthController{US: a.storage.US, HG: hg, TS: ts, UF: &requestjsonentity.Fetcher[entities.User]{}}
		r.Post("/api/user/register", ac.Register)
		r.Post("/api/user/login", ac.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.IsAuthorized(ts).Handler)
		r.Use(middlewares.CompressGzip().Handler)
		lpc := controllers.ItemsController[entities.LoginPassItem]{
			Type: datatypes.Passwords,
			IF:   &requestjsonentity.Fetcher[entities.LoginPassItem]{},
			IS:   a.storage.IS,
			IC:   entityconvertors.Factory{}.LoginPass(valueconvertors.Factory{}.LoginPass()),
		}
		r.Post("/api/update/login-pass", lpc.Update)
		r.Get("/api/updates/login-pass", lpc.Updates)

		tc := controllers.ItemsController[entities.TextItem]{
			Type: datatypes.Notes,
			IF:   &requestjsonentity.Fetcher[entities.TextItem]{},
			IS:   a.storage.IS,
			IC:   entityconvertors.Factory{}.Text(valueconvertors.Factory{}.Text()),
		}
		r.Post("/api/update/text", tc.Update)
		r.Get("/api/updates/text", tc.Updates)

		bcc := controllers.ItemsController[entities.BankCardItem]{
			Type: datatypes.BankCards,
			IF:   &requestjsonentity.Fetcher[entities.BankCardItem]{},
			IS:   a.storage.IS,
			IC:   entityconvertors.Factory{}.BankCard(valueconvertors.Factory{}.BankCard()),
		}
		r.Post("/api/update/bank-card", bcc.Update)
		r.Get("/api/updates/bank-card", bcc.Updates)

		bc := controllers.ItemsController[entities.BinaryItem]{
			Type: datatypes.Files,
			IF:   &requestjsonentity.Fetcher[entities.BinaryItem]{},
			IS:   a.storage.IS,
			IC:   entityconvertors.Factory{}.Binary(valueconvertors.Factory{}.Binary()),
		}
		r.Post("/api/update/binary", bc.Update)
		r.Get("/api/updates/binary", bc.Updates)
	})
}
