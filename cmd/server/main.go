package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/server/app"
	"github.com/AndrXxX/goph-keeper/internal/server/config"
	"github.com/AndrXxX/goph-keeper/internal/server/services/dbprovider"
	"github.com/AndrXxX/goph-keeper/internal/server/services/envparser"
	"github.com/AndrXxX/goph-keeper/internal/server/services/flagsparser"
	"github.com/AndrXxX/goph-keeper/pkg/buildformatter"
	"github.com/AndrXxX/goph-keeper/pkg/configprovider"
	"github.com/AndrXxX/goph-keeper/pkg/filestorage"
	"github.com/AndrXxX/goph-keeper/pkg/hashgenerator"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	c, icErr := initConfig()
	if icErr != nil {
		log.Fatal(icErr)
	}
	initBuildInfo()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s, err := initStorage(ctx, c)
	if err != nil {
		logger.Log.Error("init storage", zap.Error(err))
		return
	}

	defer stop()
	if err := app.New(c, *s).Run(ctx); err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func initConfig() (*config.Config, error) {
	cp := configprovider.New[config.Config](config.NewConfig(), flagsparser.Parser{}, envparser.Parser{})
	c, fcErr := cp.Fetch()
	if fcErr != nil {
		return nil, fcErr
	}
	return c, logger.Initialize(c.LogLevel, nil)
}

func initStorage(ctx context.Context, c *config.Config) (*app.Storage, error) {
	p := &dbprovider.DBProvider{DSN: c.DatabaseURI}
	db, err := p.DB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	sf := postgressql.Factory{DB: db}
	fs, err := filestorage.New(c.FileStoragePath, hashgenerator.Factory().SHA256(c.PasswordKey))
	if err != nil {
		return nil, fmt.Errorf("initialize file storage: %w", err)
	}
	return &app.Storage{DB: db, US: sf.UsersStorage(), IS: sf.StoredItemsStorage(), FS: fs}, nil
}

func initBuildInfo() {
	buildFormatter := buildformatter.BuildFormatter{
		Labels: []string{"Build version", "Build date", "Build commit"},
		Values: []string{buildVersion, buildDate, buildCommit},
	}
	for _, bInfo := range buildFormatter.Format() {
		logger.Log.Info(bInfo)
	}
}
