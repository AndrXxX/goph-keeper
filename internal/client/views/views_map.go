package views

import (
	"github.com/AndrXxX/goph-keeper/internal/client/views/lists"
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

type Map map[names.ViewName]View

func NewMap(f Factory) Map {
	return Map{
		names.AuthMenu:     f.MenusFactory().AuthMenu(),
		names.MainMenu:     lists.NewMainMenu(),
		names.PasswordList: lists.NewPasswordList(),
		names.NotesList:    lists.NewNoteList(),
		names.BankCardList: lists.NewBankCardList(),
		names.FileList:     lists.NewFileList(),
	}
}
