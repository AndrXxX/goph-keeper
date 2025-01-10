package useraccessor

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type Accessor struct {
	User *entities.User
	US   Storage[entities.User]
	AS   authSetup
}

func (u *Accessor) GetUser() *entities.User {
	return u.User
}

func (u *Accessor) SetUser(user *entities.User) {
	u.User = user
}

func (u *Accessor) SetToken(t string) {
	u.User.Token = t
}

func (u *Accessor) SetMasterPass(mp string) {
	u.User.MasterPassword = mp
}

func (u *Accessor) Auth() error {
	if u.User.Login != "" {
		// TODO: вызывает ошибку attempt to write to readonly database
		//err := u.DBProvider.RemoveDB()
		//if err != nil {
		//	return err
		//}
		created, err := u.US.Create(u.User)
		if err != nil {
			return fmt.Errorf("error saving user: %w", err)
		}
		u.User = created
		u.AS(u.User)
		return nil
	}
	list := u.US.FindAll(nil)
	for i := range list {
		if list[i].MasterPassword == u.User.MasterPassword {
			u.User = &list[i]
			u.AS(u.User)
			return nil
		}
	}
	return fmt.Errorf("пользователь с такими данными не найден")
}
