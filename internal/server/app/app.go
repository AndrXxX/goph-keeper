package app

import (
	"context"
	"errors"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vingarcia/ksql"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/server/config"
	"github.com/AndrXxX/goph-keeper/internal/server/router"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/tlsconfig"
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
	r := router.New(a.config.c, router.Storage{
		DB: a.storage.DB,
		US: a.storage.US,
		IS: a.storage.IS,
	})
	srv, err := a.runServer(r.RegisterApi())
	if err != nil {
		return err
	}

	<-commonCtx.Done()
	return a.shutdown(srv)
}

func (a *app) runServer(r *chi.Mux) (*http.Server, error) {
	tlsConfig, err := tlsconfig.NewProvider(a.config.c.PrivateCryptoKey).ForPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tls config: %w", err)
	}
	srv := &http.Server{Addr: a.config.c.Host, Handler: r, TLSConfig: tlsConfig}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("HTTP server ListenAndServe", zap.Error(err))
		}
	}()

	logger.Log.Info("listening", zap.String("host", a.config.c.Host))
	return srv, nil
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
