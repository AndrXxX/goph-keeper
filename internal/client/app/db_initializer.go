package app

import (
	"gorm.io/gorm"

	"github.com/AndrXxX/goph-keeper/internal/client/state"
)

type dbInitializer struct {
	dbProvider dbProvider
	state      *state.AppState
}

func (i *dbInitializer) Init(masterPass string, recreate bool) (*gorm.DB, error) {
	if recreate {
		err := i.dbProvider.RemoveDB()
		if err != nil {
			return nil, err
		}
	}
	actDB, err := i.dbProvider.DB(masterPass)
	if err != nil {
		return nil, err
	}
	i.state.DB = actDB
	i.state.MasterPass = masterPass
	return actDB, nil
}
