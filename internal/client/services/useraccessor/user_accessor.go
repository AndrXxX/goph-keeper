package useraccessor

import (
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
)

type Accessor struct {
	User *entities.User
	US   Storage[entities.User]
	AS   authSetup
	HG   hashGeneratorFetcher
}

func (a *Accessor) GetUser() *entities.User {
	return a.User
}

func (a *Accessor) SetUser(user *entities.User) {
	a.User = user
}

func (a *Accessor) SetToken(t string) {
	a.User.Token = t
}

func (a *Accessor) SetMasterPass(mp string) {
	a.User.MasterPassword = a.HG(mp).Generate([]byte(mp))
}

func (a *Accessor) Auth() error {
	if a.User.Login != "" {
		created, err := a.US.Create(a.User)
		if err != nil {
			return fmt.Errorf("error saving user: %w", err)
		}
		a.User = created
		a.AS(a.User)
		return nil
	}
	list := a.US.FindAll(nil)
	for i := range list {
		if list[i].MasterPassword == a.User.MasterPassword {
			a.User = &list[i]
			a.AS(a.User)
			return nil
		}
	}
	return fmt.Errorf("пользователь с такими данными не найден")
}
