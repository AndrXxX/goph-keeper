package dbprovider

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	path           = "./app.db"
	authParams     = "?_auth&_auth_user=admin&_auth_pass=%s"
	createDBParams = "?_auth&_auth_user=admin&_auth_pass=%s&_auth_crypt=sha256"
)

type DBProvider struct {
}

func (p *DBProvider) IsDBExist() bool {
	_, e := os.Stat(path)
	return e == nil || !os.IsNotExist(e)
}

// TODO: The SQLITE_USER_AUTHENTICATION extension is deprecated

func (p *DBProvider) DB(masterPass string) (*gorm.DB, error) {
	if p.IsDBExist() {
		return gorm.Open(sqlite.Open(path+fmt.Sprintf(authParams, masterPass)), &gorm.Config{})
	}
	return gorm.Open(sqlite.Open(path+fmt.Sprintf(createDBParams, masterPass)), &gorm.Config{})
}

func (p *DBProvider) RemoveDB() error {
	return os.RemoveAll(path)
}
