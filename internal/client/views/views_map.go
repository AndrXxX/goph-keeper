package views

import "github.com/AndrXxX/goph-keeper/internal/client/views/names"

type Map map[names.ViewName]View

func NewMap() Map {
	return Map{
		names.AuthMenu:     NewAuthMenu(),
		names.LoginForm:    NewLoginForm(),
		names.MainMenu:     NewMainMenu(),
		names.PasswordList: NewPasswordList(),
		names.NotesList:    NewNoteList(),
		names.BankCardList: NewBankCardList(),
		names.FileList:     NewFileList(),
	}
}
