package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/server/app"
	"github.com/AndrXxX/goph-keeper/internal/server/config"
	"github.com/AndrXxX/goph-keeper/internal/server/services/dbprovider"
	"github.com/AndrXxX/goph-keeper/pkg/buildformatter"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	c, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}
	buildFormatter := buildformatter.BuildFormatter{
		Labels: []string{"Build version", "Build date", "Build commit"},
		Values: []string{buildVersion, buildDate, buildCommit},
	}
	for _, bInfo := range buildFormatter.Format() {
		logger.Log.Info(bInfo)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s, err := initStorage(c)
	if err != nil {
		logger.Log.Error("init storage", zap.Error(err))
		return
	}

	if err := app.New(c, *s).Run(ctx); err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func initConfig() (*config.Config, error) {
	c := config.NewConfig()
	if err := logger.Initialize(c.LogLevel); err != nil {
		return nil, err
	}
	//parseFlags(c) TODO: parseFlags
	//parseEnv(c) TODO: parseEnv
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return nil, err
	}
	return c, nil
}

func initStorage(c *config.Config) (*app.Storage, error) {
	p := &dbprovider.DBProvider{DSN: c.DatabaseURI}
	db, err := p.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	sf := postgressql.Factory{DB: db}
	return &app.Storage{
		DB: db,
		US: sf.UsersStorage(),
	}, nil
}
