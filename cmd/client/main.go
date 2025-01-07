package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AndrXxX/goph-keeper/internal/client/app"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
	"github.com/AndrXxX/goph-keeper/internal/client/services/auth"
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
	viewsFactory := views.Factory{
		AppState: &state.AppState{
			User: &ormmodels.User{},
		},
		Loginer:    ap,
		Registerer: ap,
	}
	if err := app.New(views.NewMap(viewsFactory)).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
