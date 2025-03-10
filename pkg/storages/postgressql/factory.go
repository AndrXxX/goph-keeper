package postgressql

import (
	"github.com/vingarcia/ksql"
)

// Factory фабрика для хранилищ моделей
type Factory struct {
	DB ksql.Provider
}

// UsersStorage хранилище для модели User
func (f *Factory) UsersStorage() *usersStorage {
	return &usersStorage{f.DB}
}

// StoredItemsStorage хранилище для модели StoredItem
func (f *Factory) StoredItemsStorage() *storedItemsStorage {
	return &storedItemsStorage{f.DB}
}
