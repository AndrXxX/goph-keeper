package storageadapters

import (
	e "github.com/AndrXxX/goph-keeper/internal/client/entities"
	orm "github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
	"github.com/AndrXxX/goph-keeper/internal/client/services/convertors"
)

type Factory struct {
}

func (f Factory) ORMPasswordsAdapter(s Storage[orm.PasswordItem]) Storage[e.PasswordItem] {
	return &ORMAdapter[e.PasswordItem, orm.PasswordItem]{
		ORMConvertor:    convertors.PasswordOrmEntityConvertor{},
		EntityConvertor: convertors.PasswordEntityOrmConvertor{},
		Storage:         s,
	}
}

func (f Factory) ORMNotesAdapter(s Storage[orm.NoteItem]) Storage[e.NoteItem] {
	return &ORMAdapter[e.NoteItem, orm.NoteItem]{
		ORMConvertor:    convertors.NoteOrmEntityConvertor{},
		EntityConvertor: convertors.NoteEntityOrmConvertor{},
		Storage:         s,
	}
}

func (f Factory) ORMBankCardAdapter(s Storage[orm.BankCardItem]) Storage[e.BankCardItem] {
	return &ORMAdapter[e.BankCardItem, orm.BankCardItem]{
		ORMConvertor:    convertors.BankCardOrmEntityConvertor{},
		EntityConvertor: convertors.BankCardEntityOrmConvertor{},
		Storage:         s,
	}
}

func (f Factory) ORMUserAdapter(s Storage[orm.User]) Storage[e.User] {
	return &ORMAdapter[e.User, orm.User]{
		ORMConvertor:    &convertors.UserOrmEntityConvertor{},
		EntityConvertor: &convertors.UserEntityOrmConvertor{},
		Storage:         s,
	}
}
