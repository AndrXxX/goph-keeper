package synchronize

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/services/itemsloader"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/convertors"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/synchronizers"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	"github.com/AndrXxX/goph-keeper/pkg/requestjsonentity"
)

type Factory struct {
	RS       requestSender
	UB       urlBuilder
	Storages *Storages
}

func (f *Factory) PasswordSynchronizer() *synchronizers.PasswordSynchronizer {
	return &synchronizers.PasswordSynchronizer{
		LC: convertors.ListConvertor[entities.PasswordItem]{},
		L: &itemsloader.ItemsLoader[entities.PasswordItem]{
			Sender: f.RS, URLBuilder: f.UB, Fetcher: &requestjsonentity.Fetcher[entities.PasswordItem]{},
		},
		S: f.Storages.Password,
	}
}

func (f *Factory) NotesSynchronizer() *synchronizers.NotesSynchronizer {
	return &synchronizers.NotesSynchronizer{
		LC: convertors.ListConvertor[entities.NoteItem]{},
		L: &itemsloader.ItemsLoader[entities.NoteItem]{
			Sender: f.RS, URLBuilder: f.UB, Fetcher: &requestjsonentity.Fetcher[entities.NoteItem]{},
		},
		S: f.Storages.Note,
	}
}

func (f *Factory) BankCardSynchronizer() *synchronizers.BankCardSynchronizer {
	return &synchronizers.BankCardSynchronizer{
		LC: convertors.ListConvertor[entities.BankCardItem]{},
		L: &itemsloader.ItemsLoader[entities.BankCardItem]{
			Sender: f.RS, URLBuilder: f.UB, Fetcher: &requestjsonentity.Fetcher[entities.BankCardItem]{},
		},
		S: f.Storages.BankCard,
	}
}

func (f *Factory) Map() map[string]Synchronizer {
	return map[string]Synchronizer{
		datatypes.Passwords: f.PasswordSynchronizer(),
		datatypes.Notes:     f.NotesSynchronizer(),
		datatypes.BankCards: f.BankCardSynchronizer(),
	}
}
