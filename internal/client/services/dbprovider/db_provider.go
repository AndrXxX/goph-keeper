package dbprovider

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const path = "./app.db"

type DBProvider struct {
}

func (p *DBProvider) IsDBExist() bool {
	_, e := os.Stat(path)
	return e == nil || !os.IsNotExist(e)
}

func (p *DBProvider) DB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path), &gorm.Config{})
}

func (p *DBProvider) RemoveDB() error {
	return os.RemoveAll(path)
}
