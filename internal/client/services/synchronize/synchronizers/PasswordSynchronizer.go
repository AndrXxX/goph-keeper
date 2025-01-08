package synchronizers

import (
	"fmt"
	"net/http"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/convertors"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

type PasswordSynchronizer struct {
	LC convertors.ListConvertor[entities.PasswordItem]
	L  loader[entities.PasswordItem]
	S  interfaces.Storage[entities.PasswordItem]
}

func (s *PasswordSynchronizer) Sync(updates []any) error {
	list, cErr := s.LC.Convert(updates)
	if cErr != nil {
		return fmt.Errorf("convert password updates: %w", cErr)
	}
	uErr := s.L.Upload(datatypes.Passwords, list)
	if uErr != nil {
		return fmt.Errorf("upload password updates: %w", cErr)
	}
	dErr := s.Download()
	if dErr != nil {
		return fmt.Errorf("download password updates: %w", dErr)
	}
	return nil
}

func (s *PasswordSynchronizer) Download() error {
	code, list := s.L.Download(datatypes.Passwords)
	if code != http.StatusOK {
		return fmt.Errorf("could not download passwords - unexpected code: %v", code)
	}
	for i := range list {
		exist := s.S.Find(&list[i])
		if exist == nil {
			_, err := s.S.Create(&list[i])
			if err != nil {
				return fmt.Errorf("create password item: %w", err)
			}
			continue
		}
		exist.Login = list[i].Login
		exist.Password = list[i].Password
		exist.StoredItem = list[i].StoredItem
		err := s.S.Update(exist)
		if err != nil {
			return fmt.Errorf("update password item: %w", err)
		}
	}
	return nil
}
