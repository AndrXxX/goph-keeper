package useraccessor

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type Accessor struct {
	User       *entities.User
	SP         storageProvider[entities.User]
	SDB        setupDb
	HG         hashGeneratorFetcher
	AfterAuth  func()
	masterPass string
}

func (a *Accessor) GetUser() *entities.User {
	return a.User
}

func (a *Accessor) GetToken() string {
	return a.User.Token
}

func (a *Accessor) SetUser(user *entities.User) {
	a.User = user
}

func (a *Accessor) SetToken(t string) {
	a.User.Token = t
}

func (a *Accessor) SetMasterPass(mp string) {
	a.masterPass = mp
	a.User.MasterPassword = a.HG(mp).Generate([]byte(mp))
}

func (a *Accessor) Auth() error {
	db, err := a.SDB(a.masterPass, a.User.Login != "")
	if err != nil {
		logger.Log.Info("неверный мастер пароль", zap.Error(err))
		return fmt.Errorf("неверный мастер пароль")
	}
	storage := a.SP(db)
	if a.User.Login != "" {
		created, err := storage.Create(a.User)
		if err != nil {
			return fmt.Errorf("error saving user: %w", err)
		}
		a.User = created
		a.AfterAuth()
		return nil
	}
	list := storage.FindAll(nil)
	for i := range list {
		if list[i].MasterPassword == a.User.MasterPassword {
			a.User = &list[i]
			a.AfterAuth()
			return nil
		}
	}
	return fmt.Errorf("пользователь с такими данными не найден")
}
