package dbprovider

import (
	"fmt"
	"os"

	sqliteEncrypt "github.com/hinha/gorm-sqlite-cipher"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const path = "./client.db"

type DBProvider struct {
}

func (p *DBProvider) IsDBExist() bool {
	_, e := os.Stat(path)
	return e == nil || !os.IsNotExist(e)
}

func (p *DBProvider) DB(key string) (*gorm.DB, error) {
	dsn := path + fmt.Sprintf("?_pragma_key=%s&_pragma_cipher_page_size=4096", key)
	return gorm.Open(sqliteEncrypt.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func (p *DBProvider) RemoveDB() error {
	return os.RemoveAll(path)
}
