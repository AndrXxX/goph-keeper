package lists

import (
	"time"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/views/contract"
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
	"github.com/AndrXxX/goph-keeper/internal/client/views/helpers"
)

const refreshListInterval = 2 * time.Second

type Factory struct {
	FF *forms.Factory
	SM contract.SyncManager
}

func (f *Factory) AuthMenu() *authMenu {
	m := newAuthMenu()
	m.f = f.FF
	return m
}

func (f *Factory) MainMenu() *mainMenu {
	return newMainMenu()
}

func (f *Factory) PasswordList() *passwordList {
	return newPasswordList()
}

func (f *Factory) NoteList() *noteList {
	l := newNoteList()
	l.sm = f.SM
	l.lr = &helpers.ListRefresher[entities.NoteItem]{
		S:    f.FF.AppState.Storages.Note,
		List: &l.list,
	}
	return l
}

func (f *Factory) BankCardList() *bankCardList {
	l := newBankCardList()
	l.sm = f.SM
	l.lr = &helpers.ListRefresher[entities.BankCardItem]{
		S:    f.FF.AppState.Storages.BankCard,
		List: &l.list,
	}
	return l
}

func (f *Factory) FileList() *fileList {
	return newFileList()
}
