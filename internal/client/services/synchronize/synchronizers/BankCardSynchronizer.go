package synchronizers

import (
	"fmt"
	"net/http"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/convertors"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

type BankCardSynchronizer struct {
	LC convertors.ListConvertor[entities.BankCardItem]
	L  loader[entities.BankCardItem]
	S  interfaces.Storage[entities.BankCardItem]
}

func (s *BankCardSynchronizer) Sync(updates []any) error {
	list, cErr := s.LC.Convert(updates)
	if cErr != nil {
		return fmt.Errorf("convert bank card updates: %w", cErr)
	}
	uErr := s.L.Upload(datatypes.BankCards, list)
	if uErr != nil {
		return fmt.Errorf("upload bank card updates: %w", cErr)
	}
	dErr := s.Download()
	if dErr != nil {
		return fmt.Errorf("download bank card updates: %w", dErr)
	}
	return nil
}

func (s *BankCardSynchronizer) Download() error {
	code, list := s.L.Download(datatypes.BankCards)
	if code != http.StatusOK {
		return fmt.Errorf("download bank cards - unexpected code: %v", code)
	}
	for i := range list {
		exist := s.S.Find(&list[i])
		if exist == nil {
			_, err := s.S.Create(&list[i])
			if err != nil {
				return fmt.Errorf("create bank card item: %w", err)
			}
			continue
		}
		exist.Number = list[i].Number
		exist.CVCCode = list[i].CVCCode
		exist.Validity = list[i].Validity
		exist.Cardholder = list[i].Cardholder
		exist.StoredItem = list[i].StoredItem
		err := s.S.Update(exist)
		if err != nil {
			return fmt.Errorf("update bank card item: %w", err)
		}
	}
	return nil
}
