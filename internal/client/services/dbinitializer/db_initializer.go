package dbinitializer

import (
	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/state"
)

type Initializer struct {
	Provider dbProvider
	State    *state.AppState
}

func (i *Initializer) Init(masterPass string, recreate bool) (*gorm.DB, error) {
	if recreate {
		err := i.Provider.RemoveDB()
		if err != nil {
			return nil, err
		}
	}
	actDB, err := i.Provider.DB(masterPass)
	if err != nil {
		return nil, err
	}
	i.State.DB = actDB
	i.State.MasterPass = masterPass
	return actDB, nil
}
