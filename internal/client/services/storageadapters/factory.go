package storageadapters

import (
	e "github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
	orm "github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
	"github.com/AndrXxX/goph-keeper/internal/client/services/convertors"
)

type Factory struct {
}

func (f Factory) ORMPasswordsAdapter(s interfaces.Storage[orm.PasswordItem]) ORMAdapter[e.PasswordItem, orm.PasswordItem] {
	return ORMAdapter[e.PasswordItem, orm.PasswordItem]{
		ORMConvertor:    convertors.PasswordOrmEntityConvertor{},
		EntityConvertor: convertors.PasswordEntityOrmConvertor{},
		Storage:         s,
	}
}

func (f Factory) ORMNotesAdapter(s interfaces.Storage[orm.NoteItem]) ORMAdapter[e.NoteItem, orm.NoteItem] {
	return ORMAdapter[e.NoteItem, orm.NoteItem]{
		ORMConvertor:    convertors.NoteOrmEntityConvertor{},
		EntityConvertor: convertors.NoteEntityOrmConvertor{},
		Storage:         s,
	}
}

func (f Factory) ORMBankCardAdapter(s interfaces.Storage[orm.BankCardItem]) ORMAdapter[e.BankCardItem, orm.BankCardItem] {
	return ORMAdapter[e.BankCardItem, orm.BankCardItem]{
		ORMConvertor:    convertors.BankCardOrmEntityConvertor{},
		EntityConvertor: convertors.BankCardEntityOrmConvertor{},
		Storage:         s,
	}
}

func (f Factory) ORMUserAdapter(s interfaces.Storage[orm.User]) ORMAdapter[e.User, orm.User] {
	return ORMAdapter[e.User, orm.User]{
		ORMConvertor:    convertors.UserOrmEntityConvertor{},
		EntityConvertor: convertors.UserEntityOrmConvertor{},
		Storage:         s,
	}
}
