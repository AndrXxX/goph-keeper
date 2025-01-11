package postgressql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	db := ksqlProvider{}
	f := Factory{DB: &db}
	assert.Equal(t, &usersStorage{&db}, f.UsersStorage())
	assert.Equal(t, &storedItemsStorage{&db}, f.StoredItemsStorage())
}
