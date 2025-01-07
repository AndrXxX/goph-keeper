package state

import "github.com/AndrXxX/goph-keeper/internal/client/ormmodels"

type AppState struct {
	User *ormmodels.User
}
