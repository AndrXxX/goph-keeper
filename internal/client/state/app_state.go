package state

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type AppState struct {
	User       *entities.User
	DBProvider dbProvider
	Storages   *Storages
}

func (as AppState) Auth() error {
	if as.User.Login != "" {
		err := as.DBProvider.RemoveDB()
		if err != nil {
			return err
		}
		created, err := as.Storages.User.Create(as.User)
		if err != nil {
			return fmt.Errorf("error saving user: %w", err)
		}
		as.User = created
		return nil
	}
	exist := as.Storages.User.Find(&entities.User{Login: as.User.Login})
	if exist != nil {
		as.User = exist
		return nil
	}
	return nil
}
