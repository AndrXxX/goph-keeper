package synchronize

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/services/itemsloader"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/convertors"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/synchronizers"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

func Synchronizers() map[string]Synchronizer {
	return map[string]Synchronizer{
		datatypes.Passwords: &synchronizers.PasswordSynchronizer{
			LC: convertors.ListConvertor[entities.PasswordItem]{},
			L:  &itemsloader.ItemsLoader[entities.PasswordItem]{},
		},
		datatypes.Notes: &synchronizers.NotesSynchronizer{
			LC: convertors.ListConvertor[entities.NoteItem]{},
			L:  &itemsloader.ItemsLoader[entities.NoteItem]{},
		},
		datatypes.BankCards: &synchronizers.BankCardSynchronizer{
			LC: convertors.ListConvertor[entities.BankCardItem]{},
			L:  &itemsloader.ItemsLoader[entities.BankCardItem]{},
		},
	}
}
