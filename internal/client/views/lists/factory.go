package lists

import (
	"github.com/AndrXxX/goph-keeper/internal/client/views/forms"
)

type Factory struct {
	FF *forms.Factory
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
	return newNoteList()
}

func (f *Factory) BankCardList() *bankCardList {
	return newBankCardList()
}

func (f *Factory) FileList() *fileList {
	return newFileList()
}
