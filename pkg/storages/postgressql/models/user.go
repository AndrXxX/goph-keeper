package models

import (
	"time"
)

type User struct {
	ID        uint      `igor:"primary_key"`
	CreatedAt time.Time `sql:"default:(now() at time zone 'utc')"`
	UpdatedAt time.Time `sql:"default:(now() at time zone 'utc')"`
	Login     string
	Password  string
}

func (u User) TableName() string {
	return "users"
}
