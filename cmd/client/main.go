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

func main() {
	cfg := config.NewConfig()
	cfg.Host = "http://localhost:8081"
	_ = logger.Initialize("debug", []string{cfg.LogPath})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	if err := app.NewApp(cfg).Run(ctx); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
