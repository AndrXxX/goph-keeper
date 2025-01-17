package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/AndrXxX/goph-keeper/internal/client/app"
	"github.com/AndrXxX/goph-keeper/internal/client/config"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	_ "github.com/AndrXxX/goph-keeper/pkg/validators"
)

var buildVersion string
var buildDate string
var serverHost string

func main() {
	cfg := config.NewConfig()
	cfg.Host = serverHost
	cfg.BuildVersion = buildVersion
	cfg.BuildDate = buildDate
	_ = logger.Initialize("debug", []string{cfg.LogPath})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	a, err := app.NewApp(cfg)
	if err != nil {
		logger.Log.Fatal(err.Error())
		return
	}
	if err := a.Run(ctx); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
