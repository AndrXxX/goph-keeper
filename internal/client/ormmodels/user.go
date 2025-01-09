package ormmodels

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             uint `gorm:"primarykey"`
	Login          string
	Password       string
	Token          string
	MasterPassword string
}
