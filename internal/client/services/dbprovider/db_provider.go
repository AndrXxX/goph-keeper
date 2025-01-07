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
	return e == nil || !os.IsNotExist(e)
}

func (p *DBProvider) DB(masterPass string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path+params+masterPass), &gorm.Config{})
}

func (p *DBProvider) RemoveDB() error {
	return os.RemoveAll(path)
}
