package lists

import (
	"time"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
	"github.com/AndrXxX/goph-keeper/internal/client/views/menuitems"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

const refreshListInterval = 2 * time.Second

type Factory struct {
	FF *forms.Factory
	S  *contract.Storages
}

func (f *Factory) AuthMenu() *authMenu {
	m := newAuthMenu(
		withAuthItem(menuitems.AuthItem{Name: "Register", Code: "register", Desc: "Create a new account"}),
		withAuthItem(menuitems.AuthItem{Name: "Login", Code: "login", Desc: "Enter an exist account"}),
		withAuthItem(menuitems.AuthItem{Name: "Enter", Code: "master_pass", Desc: "Enter a master password to access"}),
	)
	m.f = f.FF
	return m
}

func (f *Factory) MainMenu() *mainMenu {
	return newMainMenu(
		withMenuItem(menuitems.MainMenuItem{Name: "Passwords", Code: datatypes.Passwords, Desc: "Manage passwords"}),
		withMenuItem(menuitems.MainMenuItem{Name: "Notes", Code: datatypes.Notes, Desc: "Manage notes"}),
		withMenuItem(menuitems.MainMenuItem{Name: "Bank Cards", Code: datatypes.BankCards, Desc: "Manage bank cards"}),
		withMenuItem(menuitems.MainMenuItem{Name: "Files", Code: datatypes.Files, Desc: "Manage files"}),
	)
}

func (f *Factory) PasswordList() *passwordList {
	l := newPasswordList()
	l.lr = &helpers.ListRefresher[entities.PasswordItem]{
		S:    f.S.Password,
		List: &l.list,
	}
	return l
}

func (f *Factory) NoteList() *noteList {
	l := newNoteList()
	l.lr = &helpers.ListRefresher[entities.NoteItem]{
		S:    f.S.Note,
		List: &l.list,
	}
	return l
}

func (f *Factory) BankCardList() *bankCardList {
	l := newBankCardList()
	l.lr = &helpers.ListRefresher[entities.BankCardItem]{
		S:    f.S.BankCard,
		List: &l.list,
	}
	return l
}

func (f *Factory) FileList() *fileList {
	l := newFileList()
	l.lr = &helpers.ListRefresher[entities.FileItem]{S: f.S.File, List: &l.list}
	return newFileList()
}
