package models

import (
	"time"
)

type User struct {
	ID        uint      `ksql:"id"`
	CreatedAt time.Time `ksql:"created_at"`
	UpdatedAt time.Time `ksql:"updated_at"`
	Login     string    `ksql:"login"`
	Password  string    `ksql:"password"`
}

func (u User) TableName() string {
	return "users"
}
