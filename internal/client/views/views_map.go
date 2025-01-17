package views

import (
	"github.com/AndrXxX/goph-keeper/internal/client/views/names"
)

type Map map[names.ViewName]View

func AuthViewsMap(f *Factory) Map {
	mf := f.MenusFactory()
	return Map{
		names.AuthMenu: mf.AuthMenu(),
	}
}

func MainViewsMap(f *Factory) Map {
	mf := f.MenusFactory()
	return Map{
		names.MainMenu:     mf.MainMenu(),
		names.PasswordList: mf.PasswordList(),
		names.NotesList:    mf.NoteList(),
		names.BankCardList: mf.BankCardList(),
		names.FileList:     mf.FileList(),
	}
}
