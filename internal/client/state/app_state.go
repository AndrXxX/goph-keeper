package state

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type AppState struct {
	User       *entities.User
	DBProvider dbProvider
	Storages   *Storages
	AS         authSetup
}

func (as *AppState) Auth() error {
	if as.User.Login != "" {
		// TODO: вызывает ошибку attempt to write to readonly database
		//err := as.DBProvider.RemoveDB()
		//if err != nil {
		//	return err
		//}
		created, err := as.Storages.User.Create(as.User)
		if err != nil {
			return fmt.Errorf("error saving user: %w", err)
		}
		as.User = created
		as.AS(as.User)
		return nil
	}
	list := as.Storages.User.FindAll(nil)
	for i := range list {
		if list[i].MasterPassword == as.User.MasterPassword {
			as.User = &list[i]
			as.AS(as.User)
			return nil
		}
	}
	return fmt.Errorf("пользователь с такими данными не найден")
}
