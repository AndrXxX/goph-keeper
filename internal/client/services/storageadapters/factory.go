package storageadapters

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
	"github.com/AndrXxX/goph-keeper/internal/client/services/convertors"
)

type Factory struct {
}

func (f Factory) ORMPasswordsAdapter() ORMAdapter[entities.PasswordItem, ormmodels.PasswordItem] {
	return ORMAdapter[entities.PasswordItem, ormmodels.PasswordItem]{
		ORMConvertor:    convertors.PasswordOrmEntityConvertor{},
		EntityConvertor: convertors.PasswordEntityOrmConvertor{},
	}
}

func (f Factory) ORMNotesAdapter() ORMAdapter[entities.NoteItem, ormmodels.NoteItem] {
	return ORMAdapter[entities.NoteItem, ormmodels.NoteItem]{
		ORMConvertor:    convertors.NoteOrmEntityConvertor{},
		EntityConvertor: convertors.NoteEntityOrmConvertor{},
	}
}

func (f Factory) ORMBankCardAdapter() ORMAdapter[entities.BankCardItem, ormmodels.BankCardItem] {
	return ORMAdapter[entities.BankCardItem, ormmodels.BankCardItem]{
		ORMConvertor:    convertors.BankCardOrmEntityConvertor{},
		EntityConvertor: convertors.BankCardEntityOrmConvertor{},
	}
}
