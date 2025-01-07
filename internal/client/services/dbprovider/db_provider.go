package dbprovider

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	path   = "./app.db"
	params = "?_auth&_auth_user=admin&_auth_pass="
)

type DBProvider struct {
}

func (p *DBProvider) IsDBExist() bool {
	_, e := os.Stat(path)
	if e != nil && os.IsNotExist(e) {
		return false
	}
	return true
}

func (p *DBProvider) DB(masterPass string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path+params+masterPass), &gorm.Config{})
}
