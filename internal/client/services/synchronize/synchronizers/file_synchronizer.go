package synchronizers

import (
	"fmt"
	"net/http"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/convertors"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/syncerr"
)

type FileSynchronizer struct {
	LC convertors.ListConvertor[entities.FileItem]
	L  fileLoader
	S  Storage[entities.FileItem]
}

func (s *FileSynchronizer) Sync(updates []any) error {
	list, cErr := s.LC.Convert(updates)
	if cErr != nil {
		return fmt.Errorf("convert file updates: %w", cErr)
	}
	code, uErr := s.L.Upload(list)
	if code == http.StatusUnauthorized {
		return syncerr.UnauthorizedError
	}
	if uErr != nil {
		return fmt.Errorf("upload file updates: %w", cErr)
	}
	if code != http.StatusOK {
		return fmt.Errorf("upload file updates - unexpected code: %v", code)
	}
	dErr := s.Download()
	if dErr != nil {
		return fmt.Errorf("download file updates: %w", dErr)
	}
	return nil
}

func (s *FileSynchronizer) Download() error {
	code, list := s.L.Download()
	if code == http.StatusUnauthorized {
		return syncerr.UnauthorizedError
	}
	if code != http.StatusOK {
		return fmt.Errorf("could not download files - unexpected code: %v", code)
	}
	for i := range list {
		exist := s.S.Find(&entities.FileItem{StoredItem: entities.StoredItem{ID: list[i].ID}})
		if exist == nil {
			_, err := s.S.Create(&list[i])
			if err != nil {
				return fmt.Errorf("create file item: %w", err)
			}
			continue
		}
		exist.StoredItem = list[i].StoredItem
		err := s.S.Update(exist)
		if err != nil {
			return fmt.Errorf("update file item: %w", err)
		}
	}
	return nil
}
