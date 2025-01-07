package main

import (
	"fmt"
	"os"

	"github.com/AndrXxX/goph-keeper/internal/client/app"
	"github.com/AndrXxX/goph-keeper/internal/client/views"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

func main() {
	_ = logger.Initialize("debug", []string{"./client.log"})
	if err := app.New(views.NewMap()).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
