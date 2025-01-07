package views

import "github.com/AndrXxX/goph-keeper/internal/client/views/names"

type Map map[names.ViewName]View

func NewMap(f Factory) Map {
	return Map{
		names.AuthMenu:     f.AuthMenu(),
		names.RegisterForm: f.RegisterForm(),
		names.MainMenu:     NewMainMenu(),
		names.PasswordList: NewPasswordList(),
		names.NotesList:    NewNoteList(),
		names.BankCardList: NewBankCardList(),
		names.FileList:     NewFileList(),
	}
}
