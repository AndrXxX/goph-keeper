package views

import "github.com/AndrXxX/goph-keeper/internal/client/views/names"

type Map map[names.ViewName]View

func NewMap() Map {
	return Map{
		names.AuthMenu:     NewAuthMenu(),
		names.AuthForm:     NewAuthForm(),
		names.PasswordList: NewPasswordList(),
	}
}
