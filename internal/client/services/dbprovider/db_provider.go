package dbprovider

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const path = "./app.db?_auth&_auth_user=admin&_auth_pass="

type DBProvider struct {
}

func (p *DBProvider) DB(masterPass string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path+masterPass), &gorm.Config{})
}
